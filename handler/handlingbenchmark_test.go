package handler

import (
	"time"
	"testing"
	"fmt"
	"github.com/shawnritchie/gokju/eventbus"
	"reflect"
	"github.com/shawnritchie/gokju/event"
)

type BenchmarkEvent struct {
	amount int
}

func (e BenchmarkEvent)EventID() event.Identifier {
	return event.Identifier("this.is.the.unique.event.identifier.DummyEvent1")
}


type BenchmarkEventListener struct {
	*BlockingHandler
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
		case event := <- l.ChannelOut():
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
	queue := eventbus.NewBlockingQueue()
	listener.BlockingHandler = &BlockingHandler{
		Emitter: &queue,
		Handler: &Simple{
			Consumer: &queue,
			Listener: listener,
			handlers: map[reflect.Type]func(c event.Container){},
		},
	}
	go listener.HardCodedRouter()

	m := event.MetaData{
		seqKey: uint64(1),
		timestampKey: time.Now(),
	}

	for n := 0; n < b.N; n++ {
		listener.Emitter.ChannelIn() <- event.Container{
			Event:BenchmarkEvent{amount:n},
			MetaData: m,
		}
	}
}

func BenchmarkBlockingRouter(b *testing.B) {
	routingContext := event.NewContainerContext(containeKeyDef)
	listener := BenchmarkEventListener{total: 0}
	listener.BlockingHandler = NewBlockingHandler(&routingContext, &listener)

	m := event.MetaData{
		seqKey: uint64(1),
		timestampKey: time.Now(),
	}

	for n := 0; n < b.N; n++ {
		listener.SendAndWait(event.Container{
			Event:BenchmarkEvent{amount:n},
			MetaData: m,
		})
	}
	b.StopTimer()
}