package routing

import (
	"reflect"
)

type BlockingRouter struct {
	aggregate Aggregate
	handlers map[reflect.Type]reflect.Value
	events chan EventContainer
	lifecycle <-chan interface{}
}

func NewBlockingRouter(a Aggregate) *BlockingRouter {
	b := &BlockingRouter{
		aggregate: a,
		events: make(chan EventContainer),
		lifecycle: make(chan interface{}),
	}
	b.handlers = extractHandlers(b.aggregate)
	go eventRouter(b)
	return b
}

func (b *BlockingRouter)Aggregate() Aggregate {
	return b.aggregate
}

func (b *BlockingRouter)Handlers() map[reflect.Type]reflect.Value{
	return b.handlers
}

func (b *BlockingRouter)EventConsumer() <- chan EventContainer {
	return b.events
}

func (b *BlockingRouter)EventEmitter() chan <- EventContainer {
	return b.events
}

func (b *BlockingRouter)Send(e EventContainer) {
	go func() {
		b.events <- e
	}()
}

func (b *BlockingRouter)SendAndWait(e EventContainer) {
	b.events <- e
}
