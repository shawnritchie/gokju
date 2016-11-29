package event

import (
	"testing"
	"time"
)

type (
	eventBusListener struct {
		total int
		results chan int
	}

	valueAddedEvent struct {
		add int
	}
)

func (a *eventBusListener)Address() Address {
	return Address("unique.address.bro")
}

func (a *eventBusListener)ValueAddedEventHandler(e valueAddedEvent) {
	a.total = a.total + e.add
	a.results <- a.total
}

func (e valueAddedEvent)EventID() Identifier { return Identifier("EventID") }
func (e valueAddedEvent)Version() int { return 0 }

func TestSimpleEventBus_Subscribe(t *testing.T) {
	eventBus := NewSimpleEventBus()
	listener := eventBusListener{total:0, results: make(chan int)}
	eventBus.Subscribe(NewBlockingProcessor(&routingContext, &listener))
	eventBus.Publish(Container{
		Event: valueAddedEvent{add: 5},
		MetaData : MetaData{
			seqKey: uint64(1),
			timestampKey: time.Now(),
		},
	})

	result := <- listener.results
	if (result != 5) {
		t.Errorf("The incorrect total has been passed down the event handler expect 5 received %v", result)
	}
}

func TestSimpleEventBus_UnSubscribe(t *testing.T) {
	eventBus := NewSimpleEventBus()
	listener := eventBusListener{total:0, results: make(chan int)}
	unsubscribe, _ := eventBus.Subscribe(NewBlockingProcessor(&routingContext, &listener))
	unsubscribe()
	eventBus.Publish(Container{
		Event: valueAddedEvent{add: 5},
		MetaData : MetaData{
			seqKey: uint64(1),
			timestampKey: time.Now(),
		},
	})

	if len(eventBus.eventProcessors) != 0 {
		t.Errorf("Expected number of even processors 0 while we have %v registered", len(eventBus.eventProcessors))
	}
}