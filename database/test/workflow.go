package test

import (
	"errors"
	"github.com/AsimovNetwork/asimov/database"
)



type TestCaseContainer struct {
	testCases []TestCase
}

func ( tcc *TestCaseContainer ) execute( ctx DbContext ) error {
	if tcc.testCases == nil || len(tcc.testCases) == 0 {
		return nil
	}

	var err error
	for _, tc := range tcc.testCases {
		err = tc.execute( ctx )
		if err != nil {
			//break
		}
	}

	return err
}

func ( tcc *TestCaseContainer ) addCase( testCase TestCase ) error {
	if tcc == nil || testCase == nil {
		return errors.New("TestCaseContainer addCase param error " )
	}

	if tcc.testCases == nil {
		tcc.testCases = make( []TestCase, 0, 0 )
	}

	tcc.testCases = append( tcc.testCases, testCase )
	return nil
}



type Action interface {
	execute( ctx DbContext ) error
	addCase( testCase TestCase ) error
}


// TxViewAction
type DbTxViewAction struct {
	TestCaseContainer
}

func NewDbTxViewAction() *DbTxViewAction {
	return &DbTxViewAction{}
}

func ( action *DbTxViewAction ) addCase( testCase TestCase ) error {
	if testCase == nil || testCase.isWritable() {
		return errors.New( "test case is invalid" )
	}
	return action.TestCaseContainer.addCase( testCase )
}

// TxUpdateAction
type DbTxUpdateAction struct {
	TestCaseContainer
}

func NewDbTxUpdateAction() *DbTxUpdateAction {
	return &DbTxUpdateAction{}
}

// DbAction
type DbBaseAction struct {
	TestCaseContainer
}

func NewDbBaseAction() *DbBaseAction {
	return &DbBaseAction{}
}

// Db batch action
type DbBatchAction struct {
	TestCaseContainer
	batch database.Batch

	db database.Database

}

func NewDbBatchAction( db database.Database ) *DbBatchAction {
	return &DbBatchAction{
		db:db,
	}
}

func ( action *DbBatchAction ) execute( ctx DbContext ) error {
	ctx.batch = action.db.NewBatch()
	action.batch = ctx.batch
	return action.TestCaseContainer.execute( ctx )
}


// Db Bucket action
type DbBucketAction struct {
	TestCaseContainer

	db database.Database
	isWriable bool
}

func NewDbBucketAction( db database.Database, isWritable bool ) *DbBucketAction {
	return &DbBucketAction{
		db:db,
		isWriable:isWritable,
	}
}


func ( action *DbBucketAction ) execute( ctx DbContext ) error {
	var err error
	if action.isWriable {
		err = action.db.Update(func(tx database.Tx) error {
			ctx.bucket = tx.Metadata()
			return  action.TestCaseContainer.execute( ctx )
		} )
	} else {
		err = action.db.View(func(tx database.Tx) error {
			ctx.bucket = tx.Metadata()
			return  action.TestCaseContainer.execute( ctx )
		} )
	}

	return err
}


// Db Cursor action
type DbCursorAction struct {
	TestCaseContainer

	db database.Database
	isWriable bool
}

func NewDbCursorAction(
	db database.Database,
	isWriable bool ) *DbCursorAction {
	return &DbCursorAction{
		db:db,
		isWriable:isWriable,
	}
}

func (action *DbCursorAction) execute( ctx DbContext ) error {
	var err error
	if action.isWriable {
		err = action.db.Update(func(tx database.Tx) error {
			ctx.bucket = tx.Metadata()
			ctx.cursor = ctx.bucket.Cursor()
			return  action.TestCaseContainer.execute( ctx )
		} )
	} else {
		err = action.db.View(func(tx database.Tx) error {
			ctx.bucket = tx.Metadata()
			ctx.cursor = ctx.bucket.Cursor()
			return  action.TestCaseContainer.execute( ctx )
		} )
	}

	return err
}

// db context
type DbContext struct {
	db database.Database
	batch database.Batch
	bucket database.Bucket
	cursor database.Cursor
	tx database.Tx
}

// WorkFlow
type WorkFlow struct {
	ctx DbContext
	actions []Action
}

func NewWorkFlow(
	db database.Database ) *WorkFlow {
	ctx := DbContext{
		db:db,
	}
	return &WorkFlow{
		ctx:ctx,
	}
}

func (wf *WorkFlow) addAction( action Action ) error {
	if wf == nil || action == nil {
		return errors.New( "WorkFlow addAction param error" )
	}

	if wf.actions == nil {
		wf.actions = make( []Action, 0, 0 )
	}

	wf.actions = append( wf.actions, action )
	return nil
}

func (wf *WorkFlow) execute() error {
	if wf == nil || wf.actions == nil {
		return errors.New( "WorkFlow execute error" )
	}

	for _, action := range wf.actions {
		if action != nil {
			err := action.execute( wf.ctx )
			if err != nil {
				//return err
			}
		}
	}

	return nil
}






