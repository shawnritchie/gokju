package event

import (
	"testing"
	"time"
	"reflect"
)

type DummyEvent1 struct {
	v1 string
	v2 int
	close chan struct{}
}

func (e DummyEvent1)EventID() Identifier {
	return Identifier("this.is.the.unique.identifier.DummyEvent1")
}

func (e DummyEvent1)Version() int {
	return 0
}

type DummyEvent2 struct {
	v1 string
	v2 int
	close chan struct{}
}

func (e DummyEvent2)EventID() Identifier {
	return Identifier("this.is.the.unique.identifier.DummyEvent1")
}

func (e DummyEvent2)Version() int {
	return 0
}

type DummyEvent3 struct {
	v1 string
	v2 int
	close chan struct{}
}

func (e DummyEvent3)EventID() Identifier {
	return Identifier("this.is.the.unique.identifier.DummyEvent1")
}

func (e DummyEvent3)Version() int {
	return 0
}

type DummyEvent4 struct {
	v1 string
	v2 int
	close chan struct{}
}

func (e DummyEvent4)EventID() Identifier {
	return Identifier("this.is.the.unique.identifier.DummyEvent4")
}

func (e DummyEvent4)Version() int {
	return 0
}

type DummyEvent5 struct {
	v1 string
	v2 int
	close chan struct{}
}

func (e DummyEvent5)EventID() Identifier {
	return Identifier("this.is.the.unique.identifier.DummyEvent5")
}

func (e DummyEvent5)Version() int {
	return 0
}


type DummyListener struct {
	v1 string
	v2 int
	close chan int
}

func (a *DummyListener)Address() Address{
	return Address("unique.address.bro")
}

func (a *DummyListener)DummyEvent1Handler(event DummyEvent1, timestamp time.Time, seqNo uint64) {
	a.v1 = event.v1
	a.v2 = event.v2
	a.close <- 50
}

func (a *DummyListener)DummyEvent2Handler(event DummyEvent2, timestamp time.Time, seqNo uint64) {
	a.v1 = event.v1
	a.v2 = event.v2
}

func (a *DummyListener)DummyEvent4Handler(event DummyEvent4, timestamp time.Time) {
	a.v1 = event.v1
	a.v2 = event.v2
}

func (a *DummyListener)DummyEvent5Handler(event DummyEvent5) {
	a.v1 = event.v1
	a.v2 = event.v2
}

func (a *DummyListener)ThisIsNotAHandler(event struct{}) {
}

func (a *DummyListener)RandomMethod(event struct{}) {
}

func (a *DummyListener)privateEventHandler(event DummyEvent1, seqNo uint64, timestamp time.Time) {
}

var routingContext DefaultContainerContext = NewContainerContext(containeKeyDef)


func TestRouting_FollowsNamingConvention_Match(t *testing.T) {
	t.Log("Test that the method actually end with Handler")
	aggregate := DummyListener{v1: "", v2: 0}
	m, _ := reflect.TypeOf(&aggregate).MethodByName("ThisIsNotAHandler")

	if (!followsNamingConvention(m)) {
		t.Errorf("Expected followsNamingConvention to match ThisIsNotAHandler")
	}
}

func TestRouting_FollowsNamingConvention_NoMatch(t *testing.T) {
	t.Log("Test methods that do not end with Handler to not Match")
	aggregate := DummyListener{v1: "", v2: 0}
	m, _ := reflect.TypeOf(&aggregate).MethodByName("RandomMethod")

	if (followsNamingConvention(m)) {
		t.Errorf("Expected followsNamingConvention to not match RandomMethod")
	}
}

func TestRouting_ExtractHandlers(t *testing.T) {
	t.Log("Creating an Aggregate with 4 handlers")
	aggregate := DummyListener{v1: "", v2: 0}
	handlers := extractHandlers(&routingContext, &aggregate)
	if (handlers == nil) || (len(handlers) != 4) {
		t.Errorf("Expected 4 handlers, but instead %d were created.", len(handlers))
	}
}

func TestRouting_ShowThatPointerToStructIsImplementingAggregate(t *testing.T) {
	t.Log("Creating an Aggregate but pass the value instead of the pointer, should yield 0 handlers")
	aggregate := DummyListener{v1: "", v2: 0}
	handlers := extractHandlers(&routingContext, aggregate)
	if (handlers == nil) || (len(handlers) != 0) {
		t.Errorf("Expected 0 handlers, but instead %d were created.", len(handlers))
	}
}

func TestRouting_EventRouter_RoutingCorrectly(t *testing.T) {
	t.Log("Check if routing is working correctly that is passing an event and making sure the event handler is invoked")
	listener := DummyListener{v1: "", v2: 0, close: make(chan int)}
	eventProcessor := NewBlockingProcessor(&routingContext, &listener)

	eventProcessor.Send(Container{
		Event: DummyEvent1{v1:"test", v2:15 },
		MetaData : MetaData{
			seqKey: uint64(1),
			timestampKey: time.Now(),
		},
	})

	result := <- listener.close

	if (result != 50) {
		t.Errorf("The incorrect secret code has been passed down the event handler")
	}
}

func TestRouting_EventRouterForDifferentFunctionDefintions(t *testing.T) {
	t.Log("Check if routing is working correctly that is passing an event and making sure the event handler is invoked")
	listener := DummyListener{v1: "", v2: 0, close: make(chan int)}
	eventProcessor := NewBlockingProcessor(&routingContext, &listener)

	eventProcessor.Send(Container{
		Event: DummyEvent2{ v1:"test", v2:15 },
		MetaData : MetaData{
			seqKey: uint64(1),
			timestampKey: time.Now(),
		},
	})

	eventProcessor.Send(Container{
		Event: DummyEvent4{ v1:"test", v2:15 },
		MetaData : MetaData{
			seqKey: uint64(2),
			timestampKey: time.Now(),
		},
	})

	eventProcessor.Send(Container{
		Event: DummyEvent5{ v1:"test", v2:15 },
		MetaData : MetaData{
			seqKey: uint64(3),
			timestampKey: time.Now(),
		},
	})
}

func TestRouting_EventRouter_SendingUnknownEventDoesNotBreakStuff(t *testing.T) {
	t.Log("Check if routing is working correctly that is passing an event which doesn't have a handler doesn't cuase a panic")
	t.Log("Check if routing is working correctly that is passing an event and making sure the event handler is invoked")
	listener := DummyListener{v1: "", v2: 0, close: make(chan int)}
	eventProcessor := NewBlockingProcessor(&routingContext, &listener)

	eventProcessor.Send(Container{
		Event: DummyEvent3{ v1:"test", v2:15 },
		MetaData : MetaData{
			seqKey: uint64(1),
			timestampKey: time.Now(),
		},
	})
}


func TestReflection_IsFuncTypeTheSameBasedOnDefinition(t *testing.T) {
	aggregate := DummyListener{v1: "", v2: 0, close: make(chan int)}
	a := reflect.TypeOf(aggregate.DummyEvent1Handler)
	b := reflect.TypeOf(func(event DummyEvent1, timestamp time.Time, seqNo uint64){})
	if (a != b) {
		t.Errorf("Function definition are not the same :(")
	}
}

