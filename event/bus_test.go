package event

import (
	"testing"
	"time"
	"fmt"
)

type (
	eventBusListener struct {
		total int
		results chan int
	}

	valueAddedEvent struct {
		add int
	}

	valueAddedEventV1 struct {
		add uint64
	}
)

func (a *eventBusListener)Address() Address {
	return Address("unique.address.bro")
}

func (a *eventBusListener)ValueAddedEventHandler(e valueAddedEvent) {
	a.total = a.total + e.add
	a.results <- a.total
}

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

func TestSimpleEventBus_Register_MutateEvent(t *testing.T) {
	eventBus := NewSimpleEventBus()
	listener := eventBusListener{total:0, results: make(chan int)}
	eventBus.Subscribe(NewBlockingProcessor(&routingContext, &listener))
	eventBus.Register(Interceptor{
		Identifier: Identifier("event.valueAddedEvent"),
		Intercept: func(c Container) (Container, error) {
			e := c.Event.(valueAddedEvent)
			e.add = e.add + 1
			c.Event = e
			c.MetaData[seqKey] = 101
			return c, nil
		},
	})

	eventBus.Publish(Container{
		Event: valueAddedEvent{add: 5},
		MetaData : MetaData{
			seqKey: uint64(1),
		},
	})

	result := <- listener.results
	if (result != 6) {
		t.Errorf("The incorrect total has been passed down the event handler expect 5 received %v", result)
	}
}