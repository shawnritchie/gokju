package event

import (
	"time"
	"fmt"
	"testing"
	"reflect"
	"github.com/shawnritchie/gokju/structs"
)

type BenchmarkEvent struct {
	amount int
}

func (e BenchmarkEvent)EventID() string {
	return "this.is.the.unique.event.identifier.DummyEvent1"
}


type BenchmarkEventListener struct {
	*BlockingRouter
	total uint64
	lastProcessed time.Time
}

func (l *BenchmarkEventListener)BenchmarkEventHandler(event BenchmarkEvent, timestamp time.Time, seqNo uint64) {
	l.total = l.total + uint64(event.amount)
	l.lastProcessed = timestamp
}


func (l *BenchmarkEventListener)HardCodedRouter() {
	for {
		select {
		case event := <- l.Consumer.channelOut():
			switch t := event.Event.(type) {
			case BenchmarkEvent:
				c := event
				l.BenchmarkEventHandler(event.Event.(BenchmarkEvent),
					c.MetaData[timestampKey].(time.Time),
					c.MetaData[seqKey].(uint64))
			default:
				fmt.Printf("unexpected type %T\n", t)
			}
		}
	}
}

func BenchmarkHardCodedRouter(b *testing.B) {
	listener := BenchmarkEventListener{total: 0}
	queue := newBlockingQueue()
	listener.BlockingRouter = &BlockingRouter{
		Emitter: queue,
		Router: Router{
			Consumer: queue,
			listener: listener,
			handlers: map[reflect.Type]func(c EventContainer){},
		},
	}
	go listener.HardCodedRouter()

	m := structs.MetaData{
		seqKey: uint64(1),
		timestampKey: time.Now(),
	}

	for n := 0; n < b.N; n++ {
		listener.Emitter.channelIn() <- EventContainer{
			Event:BenchmarkEvent{amount:n},
			MetaData: m,
		}
	}
}

func BenchmarkBlockingRouter(b *testing.B) {
	routingContext := NewRouterContext(containeKeyDef)
	listener := BenchmarkEventListener{total: 0}
	listener.BlockingRouter = NewBlockingRouter(routingContext, &listener)

	m := structs.MetaData{
		seqKey: uint64(1),
		timestampKey: time.Now(),
	}

	for n := 0; n < b.N; n++ {
		listener.SendAndWait(EventContainer{
			Event:BenchmarkEvent{amount:n},
			MetaData: m,
		})
	}
	b.StopTimer()
}