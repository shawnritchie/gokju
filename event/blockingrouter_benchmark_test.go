package event

import (
	"time"
	"fmt"
	"testing"
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
	seq     uint64
	timestamp time.Time
}

func (c BenchmarkEventContainer)Event() Eventer{
	return c.event
}

func (c BenchmarkEventContainer)Seq() uint64 {
	return c.seq
}

func (c BenchmarkEventContainer)Timestamp() time.Time{
	return c.timestamp
}


type BenchmarkEventListener struct {
	*BlockingRouter
	total uint64
	lastProcessed time.Time
}

func (l *BenchmarkEventListener)BenchmarkEventHandler(event BenchmarkEvent, seqNo uint64, timestamp time.Time) {
	l.total = l.total + uint64(event.amount)
	l.lastProcessed = timestamp
}


func (l *BenchmarkEventListener)HardCodedRouter() {
	for {
		select {
		case event := <- l.Consumer.channelOut():
			switch t := event.Event().(type) {
			case BenchmarkEvent:
				c := event.(BenchmarkEventContainer)
				l.BenchmarkEventHandler(event.Event().(BenchmarkEvent), c.seq, c.timestamp)
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

	for n := 0; n < b.N; n++ {
		listener.Emitter.channelIn() <- BenchmarkEventContainer{
			event:BenchmarkEvent{amount:n},
			seq:uint64(n),
			timestamp:time.Now(),
		}
	}
}

func BenchmarkBlockingRouter(b *testing.B) {
	listener := BenchmarkEventListener{total: 0}
	listener.BlockingRouter = NewBlockingRouter(&listener)

	for n := 0; n < b.N; n++ {
		listener.SendAndWait(BenchmarkEventContainer{
			event:BenchmarkEvent{amount:n},
			seq:uint64(n),
			timestamp:time.Now(),
		})
	}
}