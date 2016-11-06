package routing

import (
	"reflect"
)

type Router interface {
	Aggregate() Aggregate
	Handlers() map[reflect.Type]reflect.Value
	EventConsumer() <- chan EventContainer
	EventEmitter() chan <- EventContainer
	Send(e EventContainer)
	SendAndWait(e EventContainer)
}