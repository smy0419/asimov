package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

///// Bitcoin address generation

const VERSION = byte(0x00)
const CHECKSUM_LENGTH = 4

var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func Base58Encode(input []byte) []byte {
	var result []byte

	x := big.NewInt(0).SetBytes(input)

	base := big.NewInt(int64(len(b58Alphabet)))
	zero := big.NewInt(0)
	mod := &big.Int{}

	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod)
		result = append(result, b58Alphabet[mod.Int64()])
	}

	ReverseBytes(result)

	for _, b := range input {
		if b == 0x00 {
			result = append([]byte{b58Alphabet[0]}, result...)
		} else {
			break
		}
	}
	return result

}

func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 0
	for _, b := range input {
		if b != b58Alphabet[0] {
			break
		}
		zeroBytes++
	}
	payload := input[zeroBytes:]
	for _, b := range payload {
		charIndex := bytes.IndexByte(b58Alphabet, b)
		result.Mul(result, big.NewInt(int64(len(b58Alphabet))))
		result.Add(result, big.NewInt(int64(charIndex)))
	}

	decoded := result.Bytes()
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zeroBytes), decoded...)

	return decoded
}

func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

type BitcoinKeys struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
}

func GetBitcoinKeys() *BitcoinKeys {
	b := &BitcoinKeys{nil, nil}
	b.newKeyPair()
	return b
}

func (b *BitcoinKeys) newKeyPair() {
	curve := elliptic.P256()
	var err error
	b.PrivateKey, err = ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	b.PublicKey = append(b.PrivateKey.PublicKey.X.Bytes(), b.PrivateKey.PublicKey.Y.Bytes()...)
}

//获取地址
func (b *BitcoinKeys) GetAddress() []byte {
	//1.ripemd160(sha256(publickey))
	ripPubKey := GeneratePublicKeyHash(b.PublicKey)
	//2.最前面添加一个字节的版本信息获得 versionPublickeyHash
	versionPublickeyHash := append([]byte{VERSION}, ripPubKey[:]...)
	//3.sha256(sha256(versionPublickeyHash))  取最后四个字节的值
	tailHash := CheckSumHash(versionPublickeyHash)
	//4.拼接最终hash versionPublickeyHash + checksumHash
	finalHash := append(versionPublickeyHash, tailHash...)
	//进行base58加密
	address := Base58Encode(finalHash)
	return address
}

func GeneratePublicKeyHash(publicKey []byte) []byte {
	sha256PubKey := sha256.Sum256(publicKey)
	r := ripemd160.New()
	r.Write(sha256PubKey[:])
	ripPubKey := r.Sum(nil)
	return ripPubKey
}

func CheckSumHash(versionPublickeyHash []byte) []byte {
	versionPublickeyHashSha1 := sha256.Sum256(versionPublickeyHash)
	versionPublickeyHashSha2 := sha256.Sum256(versionPublickeyHashSha1[:])
	tailHash := versionPublickeyHashSha2[:CHECKSUM_LENGTH]
	return tailHash
}

///// BACKUP

//通过地址获得公钥
// func GetPublicKeyHashFromAddress(address string) []byte {
// 	addressBytes := []byte(address)
// 	fullHash := Base58Decode(addressBytes)
// 	publicKeyHash := fullHash[1 : len(fullHash)-CHECKSUM_LENGTH]
// 	return publicKeyHash
// }

//检测比特币地址是否有效
// func IsVaildBitcoinAddress(address string) bool {
// 	adddressByte := []byte(address)
// 	fullHash := Base58Decode(adddressByte)
// 	if len(fullHash) != 25 {
// 		return false
// 	}
// 	prefixHash := fullHash[:len(fullHash)-CHECKSUM_LENGTH]
// 	tailHash := fullHash[len(fullHash)-CHECKSUM_LENGTH:]
// 	tailHash2 := CheckSumHash(prefixHash)
// 	if bytes.Compare(tailHash, tailHash2[:]) == 0 {
// 		return true
// 	} else {
// 		return false
// 	}
// }

func createNewAddress() (string, error) {
	keys := GetBitcoinKeys()
	bitcoinAddress := keys.GetAddress()
	fmt.Println("生成新的比特币地址:", string(bitcoinAddress))
	// fmt.Printf("比特币地址是否有效:%v\n：", IsVaildBitcoinAddress(string(bitcoinAddress)))
	return string(bitcoinAddress), nil
}
