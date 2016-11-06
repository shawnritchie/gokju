package work

import (
	"github.com/shawnritchie/gokju/message"
	"github.com/shawnritchie/gokju/structs"
)

type CorrelationDataProvider func(message.Messenger) structs.MetaData

type Phase struct {
	phase string
	started bool
	reverseCallbackOrder bool
}

const (
	NOT_STARTED Phase = Phase{"NOT_STARTED", false, false}
	STARTED Phase = Phase{"STARTED", true, false}
	PREPARE_COMMIT Phase = Phase{"PREPARE_COMMIT", true, false}
	COMMIT Phase = Phase{"COMMIT", true, false}
	ROLLBACK Phase = Phase{"ROLLBACK", true, true}
	AFTER_COMMIT Phase = Phase{"AFTER_COMMIT", true, true}
	CLEANUP Phase = Phase{"CLEANUP", false, true}
	CLOSED Phase = Phase{"CLOSED", false, true}
)

type UnitOfWorker interface {
	parent() UnitOfWorker
	phase() Phase

	start()
	commit()
	rollback()
	rollbackWitError(err error)

	onPrepareCommit()
	onCommit()
	afterCommit()
	onRollback()
	onCleanup()

	getMessage() message.Messenger

	getCorrelationData() structs.MetaData
	registerCorrelationDataProvider() CorrelationDataProvider

	rollbackOnError() bool
	execute(func())
	executeWithResults(func() interface{}) interface{}
}
