package routing

import (
	"time"
	"fmt"
	"testing"
)

type BenchmarkEvent struct {
	amount int
}

func (e BenchmarkEvent)EventID() string {
	return "this.is.the.unique.event.identifier.DummyEvent1"
}

type BenchmarkEventContainer struct {
	event     Eventer
	SeqNo     uint64
	Timestamp time.Time
}

func (c BenchmarkEventContainer)Event() Eventer{
	return c.event
}

type BenchmarkAggregate struct {
	*BlockingRouter
	total uint64
	lastProcessed time.Time
}

func (a *BenchmarkAggregate)Router() Router {
	return a.BlockingRouter
}

func (a *BenchmarkAggregate)BenchmarkEventHandler(event BenchmarkEvent, seqNo uint64, timestamp time.Time) {
	a.total = a.total + uint64(event.amount)
	a.lastProcessed = timestamp
}

func (a *BenchmarkAggregate)HardCodedRouter() {
	for {
		select {
		case event := <- a.events:
			switch t := event.Event().(type) {
			case BenchmarkEvent:
				c := event.(BenchmarkEventContainer)
				a.BenchmarkEventHandler(event.Event().(BenchmarkEvent), c.SeqNo, c.Timestamp)
			default:
				fmt.Printf("unexpected type %T\n", t)
			}
		}
	}
}

func BenchmarkHardCodedRouter(b *testing.B) {
	aggregate := BenchmarkAggregate{total: 0}
	aggregate.BlockingRouter = &BlockingRouter{
		aggregate: &aggregate,
		events: make(chan EventContainer),
		lifecycle: make(chan interface{}),
	}
	go aggregate.HardCodedRouter()

	for n := 0; n < b.N; n++ {
		aggregate.events <- BenchmarkEventContainer{
			event:BenchmarkEvent{amount:n},
			SeqNo:uint64(n),
			Timestamp:time.Now(),
		}
	}
}

func BenchmarkBlockingRouter(b *testing.B) {
	aggregate := BenchmarkAggregate{total: 0}
	aggregate.BlockingRouter = NewBlockingRouter(&aggregate)

	for n := 0; n < b.N; n++ {
		aggregate.SendAndWait(BenchmarkEventContainer{
			event:BenchmarkEvent{amount:n},
			SeqNo:uint64(n),
			Timestamp:time.Now(),
		})
	}
}