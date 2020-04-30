// Copyright (c) 2018-2020 The asimov developers
// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package mempool

import (
	"container/list"
	"fmt"
	"github.com/AsimovNetwork/asimov/ainterface"
	"github.com/AsimovNetwork/asimov/common"
	"github.com/AsimovNetwork/asimov/rpcs/rpcjson"
	"sync"
	"time"

	"github.com/AsimovNetwork/asimov/asiutil"
	"github.com/AsimovNetwork/asimov/blockchain"
	"github.com/AsimovNetwork/asimov/blockchain/indexers"
	"github.com/AsimovNetwork/asimov/mining"
	"github.com/AsimovNetwork/asimov/protos"
	"github.com/AsimovNetwork/asimov/txscript"
)

const (
	// orphanTTL is the maximum amount of time an orphan is allowed to
	// stay in the orphan pool before it expires and is evicted during the
	// next scan.
	orphanTTL = time.Minute * 15

	// orphanExpireScanInterval is the minimum amount of time in between
	// scans of the orphan pool to evict expired transactions.
	orphanExpireScanInterval = time.Minute * 5

	// forbiddenTxLimit is the maximum count of forbidden transactions in cache.
	forbiddenTxLimit = 1000

	// outpointTxLimit is the maximum count of transactions in each outpoints value.
	outpointTxLimit = 5
)

// Tag represents an identifier to use for tagging orphan transactions.  The
// caller may choose any scheme it desires, however it is common to use peer IDs
// so that orphans can be identified by which peer first relayed them.
type Tag uint64

// Config is a descriptor containing the memory pool configuration.
type Config struct {
	// Policy defines the various mempool configuration options related
	// to policy.
	Policy Policy

	// FetchUtxoView defines the function to use to fetch unspent
	// transaction output information.
	FetchUtxoView func(*asiutil.Tx, bool) (ainterface.IUtxoViewpoint, error)

	// BestHeight defines the function to use to access the block height of
	// the current best chain.
	BestHeight func() int32

	// MedianTimePast defines the function to use in order to access the
	// median time past calculated from the point-of-view of the current
	// chain tip within the best chain.
	MedianTimePast func() int64

	// CalcSequenceLock defines the function to use in order to generate
	// the current sequence lock for the given transaction using the passed
	// utxo view.
	CalcSequenceLock func(*asiutil.Tx, ainterface.IUtxoViewpoint) (*blockchain.SequenceLock, error)

	// AddrIndex defines the optional address index instance to use for
	// indexing the unconfirmed transactions in the memory pool.
	// This can be nil if the address index is not enabled.
	AddrIndex *indexers.AddrIndex

	Chain *blockchain.BlockChain

	FeesChan chan interface{}

	CheckTransactionInputs func(tx *asiutil.Tx, txHeight int32, utxoView ainterface.IUtxoViewpoint,
		b *blockchain.BlockChain) (int64, *map[protos.Assets]int64, error)
}

// Policy houses the policy (configuration parameters) which is used to
// control the mempool.
type Policy struct {
	// MaxTxVersion is the transaction version that the mempool should
	// accept.  All transactions above this version are rejected as
	// non-standard.
	MaxTxVersion uint32

	// MaxOrphanTxs is the maximum number of orphan transactions
	// that can be queued.
	MaxOrphanTxs int

	// MaxOrphanTxSize is the maximum size allowed for orphan transactions.
	// This helps prevent memory exhaustion attacks from sending a lot of
	// of big orphans.
	MaxOrphanTxSize int

	// MaxSigOpCostPerTx is the cumulative maximum cost of all the signature
	// operations in a single transaction we will relay or mine.  It is a
	// fraction of the max signature operations for a block.
	MaxSigOpCostPerTx int

	// MinTxPrice defines the minimum transaction price to be
	// considered a non-zero fee.
	MinRelayTxPrice float64
}

// orphanTx is normal transaction that references an ancestor transaction
// that is not yet available.  It also contains additional information related
// to it such as an expiration time to help prevent caching the orphan forever.
type orphanTx struct {
	tx         *asiutil.Tx
	tag        Tag
	expiration time.Time
}

// TxPool is used as a source of transactions that need to be mined into blocks
// and relayed to other peers.  It is safe for concurrent access from multiple
// peers.
type TxPool struct {
	routeMtx sync.Mutex
	existCh  chan interface{}

	mtx           sync.RWMutex
	cfg           Config
	pool          map[common.Hash]*mining.TxDesc
	orphans       map[common.Hash]*orphanTx
	orphansByPrev map[protos.OutPoint]map[common.Hash]*asiutil.Tx
	outpoints     map[protos.OutPoint]map[common.Hash]*asiutil.Tx
	forbiddenTxs  map[common.Hash]int64
	forbiddenList []common.Hash
	pennyTotal    float64 // exponentially decaying total for penny spends.
	lastPennyUnix int64   // unix time of last ``penny spend''

	// nextExpireScan is the time after which the orphan pool will be
	// scanned in order to evict orphans.  This is NOT a hard deadline as
	// the scan will only run when an orphan is added to the pool as opposed
	// to on an unconditional timer.
	nextExpireScan time.Time

	fees map[protos.Assets]int32
}

// Ensure the TxPool type implements the mining.TxSource interface.
var _ mining.TxSource = (*TxPool)(nil)

// removeOrphan is the internal function which implements the public
// RemoveOrphan.  See the comment for RemoveOrphan for more details.
//
// This function MUST be called with the mempool lock held (for writes).
func (mp *TxPool) removeOrphan(tx *asiutil.Tx, removeRedeemers bool) {
	// Nothing to do if passed tx is not an orphan.
	txHash := tx.Hash()
	otx, exists := mp.orphans[*txHash]
	if !exists {
		return
	}

	// Remove the reference from the previous orphan index.
	for _, txIn := range otx.tx.MsgTx().TxIn {
		orphans, exists := mp.orphansByPrev[txIn.PreviousOutPoint]
		if exists {
			delete(orphans, *txHash)

			// Remove the map entry altogether if there are no
			// longer any orphans which depend on it.
			if len(orphans) == 0 {
				delete(mp.orphansByPrev, txIn.PreviousOutPoint)
			}
		}
	}

	// Remove any orphans that redeem outputs from this one if requested.
	if removeRedeemers {
		prevOut := protos.OutPoint{Hash: *txHash}
		for txOutIdx := range tx.MsgTx().TxOut {
			prevOut.Index = uint32(txOutIdx)
			for _, orphan := range mp.orphansByPrev[prevOut] {
				mp.removeOrphan(orphan, true)
			}
		}
	}

	// Remove the transaction from the orphan pool.
	delete(mp.orphans, *txHash)
}

// RemoveOrphan removes the passed orphan transaction from the orphan pool and
// previous orphan index.
//
// This function is safe for concurrent access.
func (mp *TxPool) RemoveOrphan(tx *asiutil.Tx) {
	mp.mtx.Lock()
	mp.removeOrphan(tx, false)
	mp.mtx.Unlock()
}

// RemoveOrphansByTag removes all orphan transactions tagged with the provided
// identifier.
//
// This function is safe for concurrent access.
func (mp *TxPool) RemoveOrphansByTag(tag Tag) uint64 {
	var numEvicted uint64
	mp.mtx.Lock()
	for _, otx := range mp.orphans {
		if otx.tag == tag {
			mp.removeOrphan(otx.tx, true)
			numEvicted++
		}
	}
	mp.mtx.Unlock()
	return numEvicted
}

// limitNumOrphans limits the number of orphan transactions by evicting a random
// orphan if adding a new one would cause it to overflow the max allowed.
//
// This function MUST be called with the mempool lock held (for writes).
func (mp *TxPool) limitNumOrphans() error {
	// Scan through the orphan pool and remove any expired orphans when it's
	// time.  This is done for efficiency so the scan only happens
	// periodically instead of on every orphan added to the pool.
	if now := time.Now(); now.After(mp.nextExpireScan) {
		origNumOrphans := len(mp.orphans)
		for _, otx := range mp.orphans {
			if now.After(otx.expiration) {
				// Remove redeemers too because the missing
				// parents are very unlikely to ever materialize
				// since the orphan has already been around more
				// than long enough for them to be delivered.
				mp.removeOrphan(otx.tx, true)
			}
		}

		// Set next expiration scan to occur after the scan interval.
		mp.nextExpireScan = now.Add(orphanExpireScanInterval)

		numOrphans := len(mp.orphans)
		if numExpired := origNumOrphans - numOrphans; numExpired > 0 {
			log.Debugf("Expired %d %s (remaining: %d)", numExpired,
				pickNoun(numExpired, "orphan", "orphans"),
				numOrphans)
		}
	}

	// Nothing to do if adding another orphan will not cause the pool to
	// exceed the limit.
	if len(mp.orphans)+1 <= mp.cfg.Policy.MaxOrphanTxs {
		return nil
	}

	// Remove a random entry from the map.  For most compilers, Go's
	// range statement iterates starting at a random item although
	// that is not 100% guaranteed by the spec.  The iteration order
	// is not important here because an adversary would have to be
	// able to pull off preimage attacks on the hashing function in
	// order to target eviction of specific entries anyways.
	for _, otx := range mp.orphans {
		// Don't remove redeemers in the case of a random eviction since
		// it is quite possible it might be needed again shortly.
		mp.removeOrphan(otx.tx, false)
		break
	}

	return nil
}

// addOrphan adds an orphan transaction to the orphan pool.
//
// This function MUST be called with the mempool lock held (for writes).
func (mp *TxPool) addOrphan(tx *asiutil.Tx, tag Tag) {
	// Nothing to do if no orphans are allowed.
	if mp.cfg.Policy.MaxOrphanTxs <= 0 {
		return
	}

	// Limit the number orphan transactions to prevent memory exhaustion.
	// This will periodically remove any expired orphans and evict a random
	// orphan if space is still needed.
	mp.limitNumOrphans()

	mp.orphans[*tx.Hash()] = &orphanTx{
		tx:         tx,
		tag:        tag,
		expiration: time.Now().Add(orphanTTL),
	}
	for _, txIn := range tx.MsgTx().TxIn {
		if _, exists := mp.orphansByPrev[txIn.PreviousOutPoint]; !exists {
			mp.orphansByPrev[txIn.PreviousOutPoint] =
				make(map[common.Hash]*asiutil.Tx)
		}
		mp.orphansByPrev[txIn.PreviousOutPoint][*tx.Hash()] = tx
	}

	log.Debugf("Stored orphan transaction %v (total: %d)", tx.Hash(),
		len(mp.orphans))
}

// maybeAddOrphan potentially adds an orphan to the orphan pool.
//
// This function MUST be called with the mempool lock held (for writes).
func (mp *TxPool) maybeAddOrphan(tx *asiutil.Tx, tag Tag) error {
	// Ignore orphan transactions that are too large.  This helps avoid
	// a memory exhaustion attack based on sending a lot of really large
	// orphans.  In the case there is a valid transaction larger than this,
	// it will ultimtely be rebroadcast after the parent transactions
	// have been mined or otherwise received.
	//
	// Note that the number of orphan transactions in the orphan pool is
	// also limited, so this equates to a maximum memory used of
	// mp.cfg.Policy.MaxOrphanTxSize * mp.cfg.Policy.MaxOrphanTxs (which is ~5MB
	// using the default values at the time this comment was written).
	serializedLen := tx.MsgTx().SerializeSize()
	if serializedLen > mp.cfg.Policy.MaxOrphanTxSize {
		str := fmt.Sprintf("orphan transaction size of %d bytes is "+
			"larger than max allowed size of %d bytes",
			serializedLen, mp.cfg.Policy.MaxOrphanTxSize)
		return txRuleError(protos.RejectNonstandard, str)
	}

	// Add the orphan if the none of the above disqualified it.
	mp.addOrphan(tx, tag)

	return nil
}

// removeOrphanDoubleSpends removes all orphans which spend outputs spent by the
// passed transaction from the orphan pool.  Removing those orphans then leads
// to removing all orphans which rely on them, recursively.  This is necessary
// when a transaction is added to the main pool because it may spend outputs
// that orphans also spend.
//
// This function MUST be called with the mempool lock held (for writes).
func (mp *TxPool) removeOrphanDoubleSpends(tx *asiutil.Tx) {
	msgTx := tx.MsgTx()
	for _, txIn := range msgTx.TxIn {
		for _, orphan := range mp.orphansByPrev[txIn.PreviousOutPoint] {
			mp.removeOrphan(orphan, true)
		}
	}
}

// isTransactionInPool returns whether or not the passed transaction already
// exists in the main pool.
//
// This function MUST be called with the mempool lock held (for reads).
func (mp *TxPool) isTransactionInPool(hash *common.Hash) bool {
	if _, exists := mp.pool[*hash]; exists {
		return true
	}

	return false
}

// IsTransactionInPool returns whether or not the passed transaction already
// exists in the main pool.
//
// This function is safe for concurrent access.
func (mp *TxPool) IsTransactionInPool(hash *common.Hash) bool {
	// Protect concurrent access.
	mp.mtx.RLock()
	inPool := mp.isTransactionInPool(hash)
	mp.mtx.RUnlock()

	return inPool
}

// isTransactionInForbiddenPool returns whether or not the passed transaction
// already exists in the main pool.
//
// This function MUST be called with the mempool lock held (for reads).
func (mp *TxPool) isTransactionInForbiddenPool(hash *common.Hash) bool {
	if _, exists := mp.forbiddenTxs[*hash]; exists {
		return true
	}

	return false
}

// isOrphanInPool returns whether or not the passed transaction already exists
// in the orphan pool.
//
// This function MUST be called with the mempool lock held (for reads).
func (mp *TxPool) isOrphanInPool(hash *common.Hash) bool {
	if _, exists := mp.orphans[*hash]; exists {
		return true
	}

	return false
}

// IsOrphanInPool returns whether or not the passed transaction already exists
// in the orphan pool.
//
// This function is safe for concurrent access.
func (mp *TxPool) IsOrphanInPool(hash *common.Hash) bool {
	// Protect concurrent access.
	mp.mtx.RLock()
	inPool := mp.isOrphanInPool(hash)
	mp.mtx.RUnlock()

	return inPool
}

// haveTransaction returns whether or not the passed transaction already exists
// in the main pool, in the orphan pool, or the forbidden pool.
//
// This function MUST be called with the mempool lock held (for reads).
func (mp *TxPool) haveTransaction(hash *common.Hash) bool {
	return mp.isTransactionInPool(hash) || mp.isOrphanInPool(hash) || mp.isTransactionInForbiddenPool(hash)
}

// HaveTransaction returns whether or not the passed transaction already exists
// in the main pool or in the orphan pool.
//
// This function is safe for concurrent access.
func (mp *TxPool) HaveTransaction(hash *common.Hash) bool {
	// Protect concurrent access.
	mp.mtx.RLock()
	haveTx := mp.haveTransaction(hash)
	mp.mtx.RUnlock()

	return haveTx
}

// removeTransaction is the internal function which implements the public
// RemoveTransaction.  See the comment for RemoveTransaction for more details.
//
// This function MUST be called with the mempool lock held (for writes).
func (mp *TxPool) removeTransaction(tx *asiutil.Tx, removeRedeemers bool) {
	txHash := tx.Hash()
	if removeRedeemers {
		// Remove any transactions which rely on this one.
		for i := uint32(0); i < uint32(len(tx.MsgTx().TxOut)); i++ {
			prevOut := protos.OutPoint{Hash: *txHash, Index: i}
			if txRedeemers, exists := mp.outpoints[prevOut]; exists {
				for _, txRedeemer := range txRedeemers {
					mp.removeTransaction(txRedeemer, true)
				}
			}
		}
	}

	// Remove the transaction if needed.
	if txDesc, exists := mp.pool[*txHash]; exists {
		// Remove unconfirmed address index entries associated with the
		// transaction if enabled.
		if mp.cfg.AddrIndex != nil {
			mp.cfg.AddrIndex.RemoveUnconfirmedTx(txHash)
		}

		// Mark the referenced outpoints as unspent by the pool.
		for _, txIn := range txDesc.Tx.MsgTx().TxIn {
			if txRedeemers, exists := mp.outpoints[txIn.PreviousOutPoint]; exists {
				delete(txRedeemers, *txDesc.Tx.Hash())
				if len(txRedeemers) == 0 {
					delete(mp.outpoints, txIn.PreviousOutPoint)
				}
			}
		}
		delete(mp.pool, *txHash)
	}
}

// RemoveTransaction removes the passed transaction from the mempool. When the
// removeRedeemers flag is set, any transactions that redeem outputs from the
// removed transaction will also be removed recursively from the mempool, as
// they would otherwise become orphans.
//
// This function is safe for concurrent access.
func (mp *TxPool) RemoveTransaction(tx *asiutil.Tx, removeRedeemers bool) {
	// Protect concurrent access.
	mp.mtx.Lock()
	mp.removeTransaction(tx, removeRedeemers)
	mp.mtx.Unlock()
}

// RemoveDoubleSpends removes all transactions which spend outputs spent by the
// passed transaction from the memory pool.  Removing those transactions then
// leads to removing all transactions which rely on them, recursively.  This is
// necessary when a block is connected to the main chain because the block may
// contain transactions which were previously unknown to the memory pool.
//
// This function is safe for concurrent access.
func (mp *TxPool) RemoveDoubleSpends(tx *asiutil.Tx) {
	// Protect concurrent access.
	mp.mtx.Lock()
	var needRemoved []*asiutil.Tx
	for _, txIn := range tx.MsgTx().TxIn {
		if txRedeemers, ok := mp.outpoints[txIn.PreviousOutPoint]; ok {
			for txHash, txRedeemer := range txRedeemers {
				if !txHash.IsEqual(tx.Hash()) {
					needRemoved = append(needRemoved, txRedeemer)
				}
			}
		}
	}
	for _, txRedeemer := range needRemoved {
		mp.removeTransaction(txRedeemer, true)
	}
	mp.mtx.Unlock()
}

func (mp *TxPool) HasSpentInTxPool(PreviousOutPoints *[]protos.OutPoint) []protos.OutPoint {
	// Protect concurrent access.
	mp.mtx.Lock()
	defer mp.mtx.Unlock()
	spentPreOuts := make([]protos.OutPoint, 0)
	for _, preOut := range *PreviousOutPoints {
		if _, ok := mp.outpoints[preOut]; ok {
			spentPreOuts = append(spentPreOuts, preOut)
		}
	}
	return spentPreOuts
}

// addTransaction adds the passed transaction to the memory pool.  It should
// not be called directly as it doesn't perform any validation.  This is a
// helper for maybeAcceptTransaction.
//
// This function MUST be called with the mempool lock held (for writes).
func (mp *TxPool) addTransaction(utxoView ainterface.IUtxoViewpoint, tx *asiutil.Tx, height int32, fee int64, feeList *map[protos.Assets]int64) *mining.TxDesc {
	// Add the transaction to the pool and mark the referenced outpoints
	// as spent by the pool.
	txD := &mining.TxDesc{
		Tx:       tx,
		Added:    time.Now(),
		Height:   height,
		Fee:      fee,
		FeeList:  feeList,
		GasPrice: float64(fee) / float64(tx.MsgTx().TxContract.GasLimit),
	}

	mp.pool[*tx.Hash()] = txD
	for _, txIn := range tx.MsgTx().TxIn {
		txs, exist := mp.outpoints[txIn.PreviousOutPoint]
		if !exist {
			txs = make(map[common.Hash]*asiutil.Tx)
			mp.outpoints[txIn.PreviousOutPoint] = txs
		} else if len(txs) >= outpointTxLimit {
			for _, txR := range txs {
				mp.removeTransaction(txR, true)
				break
			}
		}
		txs[*tx.Hash()] = tx
	}

	// Add unconfirmed address index entries associated with the transaction
	// if enabled.
	if mp.cfg.AddrIndex != nil {
		mp.cfg.AddrIndex.AddUnconfirmedTx(tx, utxoView)
	}

	return txD
}

// checkPoolDoubleSpend checks whether or not the passed transaction is
// attempting to spend coins already spent by other transactions in the pool.
// Note it does not check for double spends against transactions already in the
// main chain.
//
// This function MUST be called with the mempool lock held (for reads).
func (mp *TxPool) checkPoolDoubleSpend(tx *asiutil.Tx, gasPrice float64) error {
	for _, txIn := range tx.MsgTx().TxIn {
		if txs, exists := mp.outpoints[txIn.PreviousOutPoint]; exists {
			for _, txR := range txs {
				if _, isForbidden := mp.forbiddenTxs[*txR.Hash()]; isForbidden {
					continue
				}
				if txDesc, poolExists := mp.pool[*txR.Hash()]; poolExists{
					if txDesc.GasPrice > gasPrice {
						str := fmt.Sprintf("output %v already spent by "+
							"transaction %v in the memory pool",
							txIn.PreviousOutPoint, txR.Hash())
						return txRuleError(protos.RejectDuplicate, str)
					}
				}
			}
		}
	}
	return nil
}

// fetchInputUtxos loads utxo details about the input transactions referenced by
// the passed transaction.  First, it loads the details form the viewpoint of
// the main chain, then it adjusts them based upon the contents of the
// transaction pool.
//
// This function MUST be called with the mempool lock held (for reads).
func (mp *TxPool) fetchInputUtxos(tx *asiutil.Tx) (ainterface.IUtxoViewpoint, error) {
	utxoView, err := mp.cfg.FetchUtxoView(tx, true)
	if err != nil {
		return nil, err
	}

	// Attempt to populate any missing inputs from the transaction pool.
	for _, txIn := range tx.MsgTx().TxIn {
		prevOut := &txIn.PreviousOutPoint
		entry := utxoView.LookupEntry(*prevOut)
		if entry != nil && !entry.IsSpent() {
			continue
		}

		if poolTxDesc, exists := mp.pool[prevOut.Hash]; exists {
			// AddTxOut ignores out of range index values, so it is
			// safe to call without bounds checking here.
			utxoView.AddTxOut(poolTxDesc.Tx, prevOut.Index,
				mining.UnminedHeight)
		}
	}

	return utxoView, nil
}

// FetchTransaction returns the requested transaction from the transaction pool.
// This only fetches from the main transaction pool and does not include
// orphans.
//
// This function is safe for concurrent access.
func (mp *TxPool) FetchTransaction(txHash *common.Hash) (*asiutil.Tx, bool, error) {
	// Protect concurrent access.
	mp.mtx.RLock()
	txDesc, exists := mp.pool[*txHash]
	_, f := mp.forbiddenTxs[*txHash]
	mp.mtx.RUnlock()

	if exists {
		return txDesc.Tx, f, nil
	}

	return nil, f, fmt.Errorf("transaction is not in the pool")
}

// FetchAnyTransaction returns the requested transaction from the transaction pool.
// This fetches from the main transaction pool and orphans.
//
// This function is safe for concurrent access.
func (mp *TxPool) FetchAnyTransaction(txHash *common.Hash) (*asiutil.Tx, error) {
	// Protect concurrent access.
	mp.mtx.RLock()
	defer mp.mtx.RUnlock()
	if txDesc, exists := mp.pool[*txHash]; exists {
		return txDesc.Tx, nil
	}
	if orphanTx, exists := mp.orphans[*txHash]; exists {
		return orphanTx.tx, nil
	}

	return nil, fmt.Errorf("transaction is not in the pool")
}

// maybeAcceptTransaction is the internal function which implements the public
// MaybeAcceptTransaction.  See the comment for MaybeAcceptTransaction for
// more details.
//
// This function MUST be called with the mempool lock held (for writes).
func (mp *TxPool) maybeAcceptTransaction(tx *asiutil.Tx, isNew, rejectDupOrphans bool) ([]*common.Hash, *mining.TxDesc, error) {
	txHash := tx.Hash()

	// Don't accept the transaction if it already exists in the pool.  This
	// applies to orphan transactions as well when the reject duplicate
	// orphans flag is set.  This check is intended to be a quick check to
	// weed out duplicates.
	if mp.isTransactionInPool(txHash) || mp.isTransactionInForbiddenPool(txHash) ||
		(rejectDupOrphans && mp.isOrphanInPool(txHash)) {
		str := fmt.Sprintf("already have transaction %v", txHash)
		return nil, nil, txRuleError(protos.RejectDuplicate, str)
	}

	if tx.Type() != asiutil.TxTypeNormal {
		str := fmt.Sprintf("bad transaction type %v", tx.Type())
		return nil, nil, txRuleError(protos.RejectMalformed, str)
	}

	// Perform preliminary sanity checks on the transaction.  This makes
	// use of blockchain which contains the invariant rules for what
	// transactions are allowed into blocks.
	err := blockchain.CheckTransactionSanity(tx)
	if err != nil {
		if cerr, ok := err.(blockchain.RuleError); ok {
			return nil, nil, chainRuleError(cerr)
		}
		return nil, nil, err
	}

	// A standalone transaction must not be a coinbase transaction.
	if blockchain.IsCoinBase(tx) {
		str := fmt.Sprintf("transaction %v is an individual coinbase",
			txHash)
		return nil, nil, txRuleError(protos.RejectInvalid, str)
	}

	// Get the current height of the main chain.  A standalone transaction
	// will be mined into the next block at best, so its height is at least
	// one more than the current height.
	bestHeight := mp.cfg.BestHeight()
	nextBlockHeight := bestHeight + 1

	medianTimePast := mp.cfg.MedianTimePast()

	// Don't allow non-standard transactions
	err = checkTransactionStandard(tx, nextBlockHeight,
		medianTimePast, mp.cfg.Policy.MaxTxVersion)
	if err != nil {
		// Attempt to extract a reject code from the error so
		// it can be retained.  When not possible, fall back to
		// a non standard error.
		rejectCode, found := extractRejectCode(err)
		if !found {
			rejectCode = protos.RejectNonstandard
		}
		str := fmt.Sprintf("transaction %v is not standard: %v",
			txHash, err)
		return nil, nil, txRuleError(rejectCode, str)
	}

	// Fetch all of the unspent transaction outputs referenced by the inputs
	// to this transaction.  This function also attempts to fetch the
	// transaction itself to be used for detecting a duplicate transaction
	// without needing to do a separate lookup.
	utxoView, err := mp.fetchInputUtxos(tx)
	if err != nil {
		if cerr, ok := err.(blockchain.RuleError); ok {
			return nil, nil, chainRuleError(cerr)
		}
		return nil, nil, err
	}

	// Don't allow the transaction if it exists in the main chain and is not
	// not already fully spent.
	prevOut := protos.OutPoint{Hash: *txHash}
	for txOutIdx := range tx.MsgTx().TxOut {
		prevOut.Index = uint32(txOutIdx)
		entry := utxoView.LookupEntry(prevOut)
		if entry != nil && !entry.IsSpent() {
			return nil, nil, txRuleError(protos.RejectDuplicate,
				"transaction already exists")
		}
		utxoView.RemoveEntry(prevOut)
	}

	// Transaction is an orphan if any of the referenced transaction outputs
	// don't exist or are already spent.  Adding orphans to the orphan pool
	// is not handled by this function, and the caller should use
	// maybeAddOrphan if this behavior is desired.
	var missingParents []*common.Hash
	for outpoint, entry := range utxoView.Entries() {
		if entry == nil || entry.IsSpent() {
			// Must make a copy of the hash here since the iterator
			// is replaced and taking its address directly would
			// result in all of the entries pointing to the same
			// memory location and thus all be the final hash.
			hashCopy := outpoint.Hash
			missingParents = append(missingParents, &hashCopy)
		}
	}
	if len(missingParents) > 0 {
		return missingParents, nil, nil
	}

	// Don't allow the transaction into the mempool unless its sequence
	// lock is active, meaning that it'll be allowed into the next block
	// with respect to its defined relative lock times.
	sequenceLock, err := mp.cfg.CalcSequenceLock(tx, utxoView)
	if err != nil {
		if cerr, ok := err.(blockchain.RuleError); ok {
			return nil, nil, chainRuleError(cerr)
		}
		return nil, nil, err
	}
	if !blockchain.SequenceLockActive(sequenceLock, nextBlockHeight,
		medianTimePast) {
		return nil, nil, txRuleError(protos.RejectNonstandard,
			"transaction's sequence locks on inputs not met")
	}

	// Perform several checks on the transaction inputs using the invariant
	// rules in blockchain for what transactions are allowed into blocks.
	// Also returns the fees associated with the transaction which will be
	// used later.
	txFee, feeList, err := mp.cfg.CheckTransactionInputs(tx, nextBlockHeight,
		utxoView, mp.cfg.Chain)
	if err != nil {
		if cerr, ok := err.(blockchain.RuleError); ok {
			return nil, nil, chainRuleError(cerr)
		}
		return nil, nil, err
	}

	feepool := mp.getFees()
	if feepool != nil && feeList != nil {
		for k := range *feeList {
			if _, ok := feepool[k]; !ok {
				errstr := fmt.Sprintf("Asset is not support as fee %v", k)
				return nil, nil, txRuleError(protos.RejectFeeUnsupport, errstr)
			}
		}
	}

	// Don't allow transactions with price too low to get into a mined block.
	gasPrice := float64(txFee) / float64(tx.MsgTx().TxContract.GasLimit)
	if gasPrice < mp.cfg.Policy.MinRelayTxPrice {
		str := fmt.Sprintf("transaction %v gas price too low: %f > %f",
			txHash, gasPrice, mp.cfg.Policy.MinRelayTxPrice)
		return nil, nil, txRuleError(protos.RejectLowGasPrice, str)
	}

	// Don't allow transactions with non-standard inputs
	err = checkInputsStandard(tx, utxoView)
	if err != nil {
		// Attempt to extract a reject code from the error so
		// it can be retained.  When not possible, fall back to
		// a non standard error.
		rejectCode, found := extractRejectCode(err)
		if !found {
			rejectCode = protos.RejectNonstandard
		}
		str := fmt.Sprintf("transaction %v has a non-standard "+
			"input: %v", txHash, err)
		return nil, nil, txRuleError(rejectCode, str)
	}

	// NOTE: if you modify this code to accept non-standard transactions,
	// you should add code here to check that the transaction does a
	// reasonable number of ECDSA signature verifications.

	// Don't allow transactions with an excessive number of signature
	// operations which would result in making it impossible to mine.  Since
	// the coinbase address itself can contain signature operations, the
	// maximum allowed signature operations per transaction is less than
	// the maximum allowed signature operations per block.
	sigOpCost, err := blockchain.GetSigOpCost(tx, false, utxoView)
	if err != nil {
		if cerr, ok := err.(blockchain.RuleError); ok {
			return nil, nil, chainRuleError(cerr)
		}
		return nil, nil, err
	}
	if sigOpCost > mp.cfg.Policy.MaxSigOpCostPerTx {
		str := fmt.Sprintf("transaction %v sigop cost is too high: %d > %d",
			txHash, sigOpCost, mp.cfg.Policy.MaxSigOpCostPerTx)
		return nil, nil, txRuleError(protos.RejectNonstandard, str)
	}

	// Verify crypto signatures for each input and reject the transaction if
	// any don't verify.
	err = blockchain.ValidateTransactionScripts(tx, utxoView,
		txscript.StandardVerifyFlags)
	if err != nil {
		if cerr, ok := err.(blockchain.RuleError); ok {
			return nil, nil, chainRuleError(cerr)
		}
		return nil, nil, err
	}

	// The transaction may not use any of the same outputs as other
	// transactions already in the pool as that would ultimately result in a
	// double spend.  This check is intended to be quick and therefore only
	// detects double spends within the transaction pool itself.  The
	// transaction could still be double spending coins from the main chain
	// at this point.  There is a more in-depth check that happens later
	// after fetching the referenced transaction inputs from the main chain
	// which examines the actual spend data and prevents double spends.
	err = mp.checkPoolDoubleSpend(tx, gasPrice)
	if err != nil {
		return nil, nil, err
	}

	// Add to transaction pool.
	txD := mp.addTransaction(utxoView, tx, bestHeight, txFee, feeList)

	log.Debugf("Accepted transaction %v (pool size: %v)", txHash,
		len(mp.pool))

	return nil, txD, nil
}

// MaybeAcceptTransaction is the main workhorse for handling insertion of new
// free-standing transactions into a memory pool.  It includes functionality
// such as rejecting duplicate transactions, ensuring transactions follow all
// rules, detecting orphan transactions, and insertion into the memory pool.
//
// If the transaction is an orphan (missing parent transactions), the
// transaction is NOT added to the orphan pool, but each unknown referenced
// parent is returned.  Use ProcessTransaction instead if new orphans should
// be added to the orphan pool.
//
// This function is safe for concurrent access.
func (mp *TxPool) MaybeAcceptTransaction(tx *asiutil.Tx, isNew bool) ([]*common.Hash, *mining.TxDesc, error) {
	// Protect concurrent access.
	mp.mtx.Lock()
	hashes, txD, err := mp.maybeAcceptTransaction(tx, isNew, true)
	mp.mtx.Unlock()

	return hashes, txD, err
}

// processOrphans is the internal function which implements the public
// ProcessOrphans.  See the comment for ProcessOrphans for more details.
//
// This function MUST be called with the mempool lock held (for writes).
func (mp *TxPool) processOrphans(acceptedTx *asiutil.Tx) []*mining.TxDesc {
	var acceptedTxns []*mining.TxDesc

	// Start with processing at least the passed transaction.
	processList := list.New()
	processList.PushBack(acceptedTx)
	for processList.Len() > 0 {
		// Pop the transaction to process from the front of the list.
		firstElement := processList.Remove(processList.Front())
		processItem := firstElement.(*asiutil.Tx)

		prevOut := protos.OutPoint{Hash: *processItem.Hash()}
		for txOutIdx := range processItem.MsgTx().TxOut {
			// Look up all orphans that redeem the output that is
			// now available.  This will typically only be one, but
			// it could be multiple if the orphan pool contains
			// double spends.  While it may seem odd that the orphan
			// pool would allow this since there can only possibly
			// ultimately be a single redeemer, it's important to
			// track it this way to prevent malicious actors from
			// being able to purposely constructing orphans that
			// would otherwise make outputs unspendable.
			//
			// Skip to the next available output if there are none.
			prevOut.Index = uint32(txOutIdx)
			orphans, exists := mp.orphansByPrev[prevOut]
			if !exists {
				continue
			}

			// Potentially accept an orphan into the tx pool.
			for _, tx := range orphans {
				missing, txD, err := mp.maybeAcceptTransaction(
					tx, true, false)
				if err != nil {
					// The orphan is now invalid, so there
					// is no way any other orphans which
					// redeem any of its outputs can be
					// accepted.  Remove them.
					mp.removeOrphan(tx, true)
					break
				}

				// Transaction is still an orphan.  Try the next
				// orphan which redeems this output.
				if len(missing) > 0 {
					continue
				}

				// Transaction was accepted into the main pool.
				//
				// Add it to the list of accepted transactions
				// that are no longer orphans, remove it from
				// the orphan pool, and add it to the list of
				// transactions to process so any orphans that
				// depend on it are handled too.
				acceptedTxns = append(acceptedTxns, txD)
				mp.removeOrphan(tx, false)
				processList.PushBack(tx)

				// Only one transaction for this outpoint can be
				// accepted, so the rest are now double spends
				// and are removed later.
				break
			}
		}
	}

	// Recursively remove any orphans that also redeem any outputs redeemed
	// by the accepted transactions since those are now definitive double
	// spends.
	mp.removeOrphanDoubleSpends(acceptedTx)
	for _, txD := range acceptedTxns {
		mp.removeOrphanDoubleSpends(txD.Tx)
	}

	return acceptedTxns
}

// ProcessOrphans determines if there are any orphans which depend on the passed
// transaction hash (it is possible that they are no longer orphans) and
// potentially accepts them to the memory pool.  It repeats the process for the
// newly accepted transactions (to detect further orphans which may no longer be
// orphans) until there are no more.
//
// It returns a slice of transactions added to the mempool.  A nil slice means
// no transactions were moved from the orphan pool to the mempool.
//
// This function is safe for concurrent access.
func (mp *TxPool) ProcessOrphans(acceptedTx *asiutil.Tx) []*mining.TxDesc {
	mp.mtx.Lock()
	acceptedTxns := mp.processOrphans(acceptedTx)
	mp.mtx.Unlock()

	return acceptedTxns
}

// ProcessTransaction is the main workhorse for handling insertion of new
// free-standing transactions into the memory pool.  It includes functionality
// such as rejecting duplicate transactions, ensuring transactions follow all
// rules, orphan transaction handling, and insertion into the memory pool.
//
// It returns a slice of transactions added to the mempool.  When the
// error is nil, the list will include the passed transaction itself along
// with any additional orphan transaactions that were added as a result of
// the passed one being accepted.
//
// This function is safe for concurrent access.
func (mp *TxPool) ProcessTransaction(tx *asiutil.Tx, allowOrphan, rateLimit bool, tag Tag) ([]*mining.TxDesc, error) {
	log.Tracef("Processing transaction %v", tx.Hash())

	// Protect concurrent access.
	mp.mtx.Lock()
	defer mp.mtx.Unlock()

	// Potentially accept the transaction to the memory pool.
	missingParents, txD, err := mp.maybeAcceptTransaction(tx, true, true)
	if err != nil {
		return nil, err
	}

	if len(missingParents) == 0 {
		// Accept any orphan transactions that depend on this
		// transaction (they may no longer be orphans if all inputs
		// are now available) and repeat for those accepted
		// transactions until there are no more.
		newTxs := mp.processOrphans(tx)
		acceptedTxs := make([]*mining.TxDesc, len(newTxs)+1)

		// Add the parent transaction first so remote nodes
		// do not add orphans.
		acceptedTxs[0] = txD
		copy(acceptedTxs[1:], newTxs)

		return acceptedTxs, nil
	}

	// The transaction is an orphan (has inputs missing).  Reject
	// it if the flag to allow orphans is not set.
	if !allowOrphan {
		// Only use the first missing parent transaction in
		// the error message.
		//
		// NOTE: RejectDuplicate is really not an accurate
		// reject code here, but it matches the reference
		// implementation and there isn't a better choice due
		// to the limited number of reject codes.  Missing
		// inputs is assumed to mean they are already spent
		// which is not really always the case.
		str := fmt.Sprintf("orphan transaction %v references "+
			"outputs of unknown or fully-spent "+
			"transaction %v", tx.Hash(), missingParents[0])
		return nil, txRuleError(protos.RejectDuplicate, str)
	}

	// Potentially add the orphan transaction to the orphan pool.
	err = mp.maybeAddOrphan(tx, tag)
	return nil, err
}

// Count returns the number of transactions in the main pool.  It does not
// include the orphan pool.
//
// This function is safe for concurrent access.
func (mp *TxPool) Count() int {
	mp.mtx.RLock()
	count := len(mp.pool)
	mp.mtx.RUnlock()

	return count
}

// TxHashes returns a slice of hashes for all of the transactions in the memory
// pool.
//
// This function is safe for concurrent access.
func (mp *TxPool) TxHashes() []*common.Hash {
	mp.mtx.RLock()
	hashes := make([]*common.Hash, len(mp.pool))
	i := 0
	for hash := range mp.pool {
		hashCopy := hash
		hashes[i] = &hashCopy
		i++
	}
	mp.mtx.RUnlock()

	return hashes
}

// TxDescs returns a slice of descriptors for all the transactions in the pool.
// The descriptors are to be treated as read only.
//
// This function is safe for concurrent access.
func (mp *TxPool) TxDescs() mining.TxDescList {
	mp.mtx.RLock()
	descs := make(mining.TxDescList, 0, len(mp.pool))
	for _, desc := range mp.pool {
		if _, exist := mp.forbiddenTxs[*desc.Tx.Hash()]; !exist {
			descs = append(descs, desc)
		}
	}
	mp.mtx.RUnlock()

	return descs
}

// UpdateForbiddenTxs put given txhashes into forbiddenTxs.
// If size of forbiddenTxs exceed limit, clear some olders.
//
// This function is safe for concurrent access.
func (mp *TxPool) UpdateForbiddenTxs(txHashes []*common.Hash, height int64) {
	mp.mtx.RLock()
	for _, txHash := range txHashes {
		mp.forbiddenTxs[*txHash] = height
		mp.forbiddenList = append(mp.forbiddenList, *txHash)
	}
	size := len(mp.forbiddenTxs)
	if size > forbiddenTxLimit {
		for i := 0; i < size - forbiddenTxLimit; i++ {
			delete(mp.forbiddenTxs, mp.forbiddenList[i])
		}
		mp.forbiddenList = mp.forbiddenList[size - forbiddenTxLimit:]
	}
	mp.mtx.RUnlock()
}

// RawMempoolVerbose returns all of the entries in the mempool as a fully
// populated rpcjson result.
//
// This function is safe for concurrent access.
func (mp *TxPool) RawMempoolVerbose() map[string]*rpcjson.GetRawMempoolVerboseResult {
	mp.mtx.RLock()
	defer mp.mtx.RUnlock()

	result := make(map[string]*rpcjson.GetRawMempoolVerboseResult,
		len(mp.pool))

	for _, desc := range mp.pool {
		// Calculate the current priority based on the inputs to
		// the transaction.  Use zero if one or more of the
		// input transactions can't be found for some reason.
		tx := desc.Tx

		mpd := &rpcjson.GetRawMempoolVerboseResult{
			Size:    int32(tx.MsgTx().SerializeSize()),
			Fee:     desc.Fee,
			Time:    desc.Added.Unix(),
			Height:  int64(desc.Height),
			Depends: make([]string, 0),
		}
		for _, txIn := range tx.MsgTx().TxIn {
			hash := &txIn.PreviousOutPoint.Hash
			if mp.haveTransaction(hash) {
				mpd.Depends = append(mpd.Depends,
					hash.UnprefixString())
			}
		}

		result[tx.Hash().UnprefixString()] = mpd
	}

	return result
}

func (mp *TxPool) updateFees(fees map[protos.Assets]int32) {
	mp.fees = fees
	for _, txdec := range mp.pool {
		for assets := range *txdec.FeeList {
			if _, ok := fees[assets]; !ok {
				mp.removeTransaction(txdec.Tx, true)
				break
			}
		}
	}
}

func (mp *TxPool) getFees() map[protos.Assets]int32 {
	return mp.fees
}

func (mp *TxPool) handleUpdateFees() {
	existCh := mp.existCh
mainloop:
	for {
		select {
		case arg := <-mp.cfg.FeesChan:
			if fees := arg.(map[protos.Assets]int32); fees != nil {
				mp.mtx.Lock()
				mp.updateFees(fees)
				mp.mtx.Unlock()
			}
		case <-existCh:
			break mainloop
		}
	}
}

func (mp *TxPool) Start() {
	mp.routeMtx.Lock()
	defer mp.routeMtx.Unlock()
	if mp.existCh != nil {
		panic("txpoll is already started")
	}
	log.Info("TxPool start")
	mp.existCh = make(chan interface{})
	go mp.handleUpdateFees()
}

func (mp *TxPool) Halt() {
	mp.routeMtx.Lock()
	defer mp.routeMtx.Unlock()
	log.Info("TxPool Stop")

	if mp.existCh != nil {
		close(mp.existCh)
		mp.existCh = nil
	}
}

// New returns a new memory pool for validating and storing standalone
// transactions until they are mined into a block.
func New(cfg *Config) *TxPool {
	return &TxPool{
		cfg:            *cfg,
		pool:           make(map[common.Hash]*mining.TxDesc),
		orphans:        make(map[common.Hash]*orphanTx),
		orphansByPrev:  make(map[protos.OutPoint]map[common.Hash]*asiutil.Tx),
		nextExpireScan: time.Now().Add(orphanExpireScanInterval),
		outpoints:      make(map[protos.OutPoint]map[common.Hash]*asiutil.Tx),
		forbiddenTxs:   make(map[common.Hash]int64),
	}
}

//////////////////////////////////////////
////////////// block signature pool
//////////////////////////////////////////
type SignatureDesc struct {
	sig *asiutil.BlockSign
	packageHeight int32
}

type SigPool struct {
	mtx  sync.RWMutex
	pool map[common.Hash]SignatureDesc
}

// Ensure the TxPool type implements the mining.TxSource interface.
var _ mining.SigSource = (*SigPool)(nil)

// MiningDescs return a block signature list used for mining.
//
// This function is safe for concurrent access.
func (mp *SigPool) MiningDescs(height int32) []*asiutil.BlockSign {
	mp.mtx.RLock()
	defer mp.mtx.RUnlock()

	descs := make([]*asiutil.BlockSign, 0, len(mp.pool))
	for _, v := range mp.pool {
		if v.packageHeight == 0 && v.sig.MsgSign.BlockHeight >= height - common.BlockSignDepth {
			descs = append(descs, v.sig)
		}
	}

	return descs
}

func NewSigPool() *SigPool {
	return &SigPool{
		pool: make(map[common.Hash]SignatureDesc),
	}
}

// HaveSignature return whether the blocksign in the mempool.
//
// This function is safe for concurrent access.
func (mp *SigPool) HaveSignature(hash *common.Hash) bool {
	mp.mtx.RLock()
	defer mp.mtx.RUnlock()

	return mp.haveSignature(hash)
}

func (mp *SigPool) haveSignature(hash *common.Hash) bool {
	signature, ok := mp.pool[*hash]
	if ok && signature.sig != nil {
		return true
	}
	return false
}

// ConnectSigns mark the passed block signature with packageHeight in the mempool.
//
// This function is safe for concurrent access.
func (mp *SigPool) ConnectSigns(sigs []*asiutil.BlockSign, height int32) {
	mp.mtx.Lock()
	defer mp.mtx.Unlock()

	deprecatedlist := make([]*asiutil.BlockSign, 0, len(mp.pool))
	for _, v := range mp.pool {
		if v.sig.MsgSign.BlockHeight < height - common.BlockSignDepth * 2 {
			deprecatedlist = append(deprecatedlist, v.sig)
		}
	}
	for _, sig := range deprecatedlist {
		delete(mp.pool, *sig.Hash())
	}
	for _, v := range sigs {
		mp.pool[*v.Hash()] = SignatureDesc{
			sig :         v,
			packageHeight:height,
		}
	}
}

// DisConnectSigns clear signatures' packageHeight which equals passed height
// from the mempool.
//
// This function is safe for concurrent access.
func (mp *SigPool) DisConnectSigns(height int32) {
	mp.mtx.Lock()
	defer mp.mtx.Unlock()

	for _, v := range mp.pool {
		if v.packageHeight == height {
			v.packageHeight = 0
		}
	}
}

// FetchSignature return the block signature via signature hash
// from the mempool.
//
// This function is safe for concurrent access.
func (mp *SigPool) FetchSignature(sigHash *common.Hash) (*asiutil.BlockSign, error) {
	mp.mtx.RLock()
	v, exists := mp.pool[*sigHash]
	mp.mtx.RUnlock()

	if exists {
		return v.sig, nil
	}

	return nil, fmt.Errorf("signature is not in the pool")
}

// ProcessSig is the main workhorse for handling insertion of new
// block signature into the memory pool.
//
// This function is safe for concurrent access.
func (mp *SigPool) ProcessSig(sig *asiutil.BlockSign) error {
	if sig == nil {
		return nil
	}
	mp.mtx.Lock()
	defer mp.mtx.Unlock()
	_, ok := mp.pool[*sig.Hash()]
	if ok {
		return nil
	}

	mp.pool[*sig.Hash()] = SignatureDesc{sig: sig}

	log.Debugf("[SigPool] Accepted signature height: %v, hash:%v (pool size: %v)",
		sig.MsgSign.BlockHeight, sig.Hash(), len(mp.pool))

	return nil
}
