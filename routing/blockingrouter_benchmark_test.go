package routing

import (
	"time"
	"fmt"
	"testing"
	"github.com/shawnritchie/gokju/structs"
	"reflect"
)

type BenchmarkEvent struct {
	amount int
}

func (e BenchmarkEvent)EventID() string {
	return "this.is.the.unique.event.identifier.DummyEvent1"
}

type BenchmarkEventContainer struct {
	event     Eventer
	metadata structs.MetaData
}

func (c BenchmarkEventContainer)Event() Eventer {
	return c.event
}

func (c BenchmarkEventContainer)MetaData() structs.MetaData {
	return c.metadata
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
				a.BenchmarkEventHandler(event.Event().(BenchmarkEvent),
					event.MetaData().Get(reflect.TypeOf(uint64(0))).(uint64),
					event.MetaData().Get(reflect.TypeOf(time.Time{})).(time.Time))
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

	m := structs.MetaData{}
	m.Add(uint64(1))
	m.Add(time.Now())

	for n := 0; n < b.N; n++ {
		aggregate.events <- BenchmarkEventContainer{
			event:BenchmarkEvent{amount:n},
			metadata: m,
		}
	}
}

func BenchmarkBlockingRouter(b *testing.B) {
	aggregate := BenchmarkAggregate{total: 0}
	aggregate.BlockingRouter = NewBlockingRouter(&aggregate)

	m := structs.MetaData{}
	m.Add(uint64(1))
	m.Add(time.Now())

	for n := 0; n < b.N; n++ {
		aggregate.SendAndWait(BenchmarkEventContainer{
			event:BenchmarkEvent{amount:n},
			metadata: m,
		})
	}
}