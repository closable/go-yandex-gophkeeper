package errors

import "errors"

var (
	ErrorJWTToken    = errors.New("error JWT token")
	ErrorCrypoSeq    = errors.New("invalid crypto keystring")
	ErrorAuthInfo    = errors.New("login or password is empty")
	ErrorLoginInfo   = errors.New("login info already present")
	ErrorLogin       = errors.New("login error")
	ErrorAddItem     = errors.New("error adding item")
	ErrorUpdateItem  = errors.New("error updating item")
	ErrorListData    = errors.New("error list data")
	ErrorDeleteData  = errors.New("delete data error")
	ErrorFileCreate  = errors.New("file create error")
	ErrorStreamData  = errors.New("stream data error")
	ErrorFileWriting = errors.New("file data writing error")
	ErrorFileReading = errors.New("file data reading error")

	ErrorTxDB      = errors.New("starting transaction error")
	ErrorPrepareDB = errors.New("prepearing sql error")
	ErrorExecDB    = errors.New("executing sql error")
	ErrorCommitDB  = errors.New("commiting sql error")
)
