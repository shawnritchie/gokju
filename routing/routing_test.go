package routing

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

func (e DummyEvent1)EventID() string {
	return "this.is.the.unique.event.identifier.DummyEvent1"
}

type DummyEvent2 struct {
	v1 string
	v2 int
	close chan struct{}
}

func (e DummyEvent2)EventID() string {
	return "this.is.the.unique.event.identifier.DummyEvent1"
}

type DummyEvent3 struct {
	v1 string
	v2 int
	close chan struct{}
}

func (e DummyEvent3)EventID() string {
	return "this.is.the.unique.event.identifier.DummyEvent1"
}

type DummyEvent4 struct {
	v1 string
	v2 int
	close chan struct{}
}

func (e DummyEvent4)EventID() string {
	return "this.is.the.unique.event.identifier.DummyEvent4"
}

type DummyEvent5 struct {
	v1 string
	v2 int
	close chan struct{}
}

func (e DummyEvent5)EventID() string {
	return "this.is.the.unique.event.identifier.DummyEvent5"
}

type DummyAggregate struct {
	*BlockingRouter
	v1 string
	v2 int
	close chan int
}

type DummyEventContainer struct {
	event     Eventer
	seq     uint64
	timestamp time.Time
}

func (c DummyEventContainer)Event() Eventer{
	return c.event
}

func (c DummyEventContainer)Seq() uint64 {
	return c.seq
}

func (c DummyEventContainer)Timestamp() time.Time{
	return c.timestamp
}

func (a *DummyAggregate)Router() Router {
	return a.BlockingRouter
}

func (a *DummyAggregate)DummyEvent1Handler(event DummyEvent1, timestamp time.Time, seqNo uint64) {
	a.v1 = event.v1
	a.v2 = event.v2
	a.close <- 50
}

func (a *DummyAggregate)DummyEvent2Handler(event DummyEvent2, timestamp time.Time, seqNo uint64) {
	a.v1 = event.v1
	a.v2 = event.v2
}

func (a *DummyAggregate)DummyEvent4Handler(event DummyEvent4, timestamp time.Time) {
	a.v1 = event.v1
	a.v2 = event.v2
}

func (a *DummyAggregate)DummyEvent5Handler(event DummyEvent5) {
	a.v1 = event.v1
	a.v2 = event.v2
}

func (a *DummyAggregate)ThisIsNotAHandler(event struct{}) {
}

func (a *DummyAggregate)RandomMethod(event struct{}) {
}

func (a *DummyAggregate)privateEventHandler(event DummyEvent1, seqNo uint64, timestamp time.Time) {
}

func TestRouting_FollowsNamingConvention_Match(t *testing.T) {
	t.Log("Test that the method actually end with Handler")
	aggregate := DummyAggregate{v1: "", v2: 0}
	m, _ := reflect.TypeOf(&aggregate).MethodByName("ThisIsNotAHandler")

	if (!followsNamingConvention(m)) {
		t.Errorf("Expected followsNamingConvention to match ThisIsNotAHandler")
	}
}

func TestRouting_FollowsNamingConvention_NoMatch(t *testing.T) {
	t.Log("Test methods that do not end with Handler to not Match")
	aggregate := DummyAggregate{v1: "", v2: 0}
	m, _ := reflect.TypeOf(&aggregate).MethodByName("RandomMethod")

	if (followsNamingConvention(m)) {
		t.Errorf("Expected followsNamingConvention to not match RandomMethod")
	}
}

func TestRouting_ExtractHandlers(t *testing.T) {
	t.Log("Creating an Aggregate with 2 handlers")
	aggregate := DummyAggregate{v1: "", v2: 0}
	handlers := extractHandlers(&aggregate)
	if (handlers == nil) || (len(handlers) != 4) {
		t.Errorf("Expected 4 handlers, but instead %d were created.", len(handlers))
	}
}

func TestRouting_ShowThatPointerToStructIsImplementingAggregate(t *testing.T) {
	t.Log("Creating an Aggregate but pass the value instead of the pointer, should yield 0 handlers")
	aggregate := DummyAggregate{v1: "", v2: 0}
	handlers := extractHandlers(aggregate)
	if (handlers == nil) || (len(handlers) != 0) {
		t.Errorf("Expected 0 handlers, but instead %d were created.", len(handlers))
	}
}

func TestRouting_EventRouter_RoutingCorrectly(t *testing.T) {
	t.Log("Check if routing is working correctly that is passing an event and making sure the event handler is invoked")
	aggregate := DummyAggregate{v1: "", v2: 0, close: make(chan int)}
	b := &BlockingRouter{
		aggregate: &aggregate,
		events: make(chan EventContainer),
		lifecycle: make(chan interface{}),
	}
	b.handlers = extractHandlers(b.aggregate)
	go eventRouter(b)

	b.events <- DummyEventContainer{
		event: DummyEvent1{ v1:"test", v2:15 },
		seq: uint64(1),
		timestamp: time.Now(),
	}

	result := <- aggregate.close

	if (result != 50) {
		t.Errorf("The incorrect secret code has been passed down the event handler")
	}
}

func TestRouting_EventRouterForDifferentFunctionDefintions(t *testing.T) {
	t.Log("Check if routing is working correctly that is passing an event and making sure the event handler is invoked")
	aggregate := DummyAggregate{v1: "", v2: 0, close: make(chan int)}
	b := &BlockingRouter{
		aggregate: &aggregate,
		events: make(chan EventContainer),
		lifecycle: make(chan interface{}),
	}
	b.handlers = extractHandlers(b.aggregate)
	go eventRouter(b)

	b.events <- DummyEventContainer{
		event: DummyEvent2{ v1:"test", v2:15 },
		seq: uint64(1),
		timestamp: time.Now(),
	}

	b.events <- DummyEventContainer{
		event: DummyEvent4{ v1:"test", v2:15 },
		seq: uint64(2),
		timestamp: time.Now(),
	}

	b.events <- DummyEventContainer{
		event: DummyEvent5{ v1:"test", v2:15 },
		seq: uint64(3),
		timestamp: time.Now(),
	}
}

func TestRouting_EventRouter_SendingUnknownEventDoesNotBreakStuff(t *testing.T) {
	t.Log("Check if routing is working correctly that is passing an event which doesn't have a handler doesn't cuase a panic")
	aggregate := DummyAggregate{v1: "", v2: 0, close: make(chan int)}
	b := &BlockingRouter{
		aggregate: &aggregate,
		events: make(chan EventContainer),
		lifecycle: make(chan interface{}),
	}
	b.handlers = extractHandlers(b.aggregate)
	go eventRouter(b)

	b.events <- DummyEventContainer{
		event: DummyEvent3{ v1:"test", v2:15 },
		seq: uint64(1),
		timestamp: time.Now(),
	}
}


func TestReflection_IsFuncTypeTheSameBasedOnDefinition(t *testing.T) {
	aggregate := DummyAggregate{v1: "", v2: 0, close: make(chan int)}
	a := reflect.TypeOf(aggregate.DummyEvent1Handler)
	b := reflect.TypeOf(func(event DummyEvent1, timestamp time.Time, seqNo uint64){})
	if (a != b) {
		t.Errorf("Function definition are not the same :(")
	}
}