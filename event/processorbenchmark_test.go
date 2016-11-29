package event

import (
	"time"
	"testing"
	"reflect"
	"fmt"
)

type BenchmarkEvent struct {
	amount int
}

func (e BenchmarkEvent)EventID() Identifier {
	return Identifier("this.is.the.unique.identifier.DummyEvent1")
}

func (e BenchmarkEvent)Version() int {
	return 0
}

type BenchmarkEventListener struct {
	*BlockingEventProcessor
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
		case event := <- l.Consume():
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
	queue := NewBlockingQueue()
	listener.BlockingEventProcessor = &BlockingEventProcessor{
		Emitter: &queue,
		Consumer: &queue,
		Addressable: listener,
		handlers: map[reflect.Type]func(c Container){},
	}
	go listener.HardCodedRouter()

	m := MetaData{
		seqKey: uint64(1),
		timestampKey: time.Now(),
	}

	for n := 0; n < b.N; n++ {
		listener.Emitter.Emit() <- Container{
			Event:BenchmarkEvent{amount:n},
			MetaData: m,
		}
	}
}

func BenchmarkBlockingRouter(b *testing.B) {
	routingContext := NewContainerContext(containeKeyDef)
	listener := BenchmarkEventListener{total: 0}
	listener.BlockingEventProcessor = NewBlockingProcessor(&routingContext, &listener)

	m := MetaData{
		seqKey: uint64(1),
		timestampKey: time.Now(),
	}

	for n := 0; n < b.N; n++ {
		listener.SendAndWait(Container{
			Event:BenchmarkEvent{amount:n},
			MetaData: m,
		})
	}
	b.StopTimer()
}