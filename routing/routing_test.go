package routing

import (
	"testing"
	"time"
	"reflect"
	"github.com/shawnritchie/gokju/structs"
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

type DummyAggregate struct {
	*BlockingRouter
	v1 string
	v2 int
	close chan int
}

type DummyEventContainer struct {
	event     Eventer
	metadata structs.MetaData
}

func NewDummyEventContainer(e Eventer) DummyEventContainer {
	m := structs.MetaData{}
	m.Add(uint64(1))
	m.Add(time.Now())

	return DummyEventContainer{
		event:e,
		metadata: m,
	}
}

func (c DummyEventContainer)Event() Eventer {
	return c.event
}

func (c DummyEventContainer)MetaData() structs.MetaData {
	return c.metadata
}

func (a *DummyAggregate)Router() Router {
	return a.BlockingRouter
}

func (a *DummyAggregate)DummyEvent1Handler(event DummyEvent1, seqNo uint64, timestamp time.Time) {
	a.v1 = event.v1
	a.v2 = event.v2
	a.close <- 50
}

func (a *DummyAggregate)DummyEvent2Handler(event DummyEvent2, seqNo uint64, timestamp time.Time) {
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

func TestRouting_DoesConsumeEventer_Match(t *testing.T) {
	t.Log("Test that the first parameter to a method is actually implementing Eventer interface")
	aggregate := DummyAggregate{v1: "", v2: 0}
	m := reflect.ValueOf(&aggregate).MethodByName("DummyEvent1Handler").Type()

	if (!doesConsumeEventer(m)) {
		t.Errorf("privateEventHandler first paramter is implementing Eventer")
	}
}

func TestRouting_DoesConsumeEventer_NoMatch(t *testing.T) {
	t.Log("Test that the first parameter to a method is not actually implementing Eventer interface")
	aggregate := DummyAggregate{v1: "", v2: 0}
	m := reflect.ValueOf(&aggregate).MethodByName("ThisIsNotAHandler").Type()

	if (doesConsumeEventer(m)) {
		t.Errorf("ThisIsNotAHandler first paramter is not implementing Eventer")
	}
}

func TestRouting_ExtractHandlers(t *testing.T) {
	t.Log("Creating an Aggregate with 2 handlers")
	aggregate := DummyAggregate{v1: "", v2: 0}
	handlers := extractHandlers(&aggregate)
	if (handlers == nil) || (len(handlers) != 2) {
		t.Errorf("Expected 2 handlers, but instead %d were created.", len(handlers))
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

	var check int = 50
	b.events <- NewDummyEventContainer(DummyEvent1{ v1:"test", v2:check })

	result := <- aggregate.close
	if (result != check) {
		t.Errorf("The incorrect secret code has been passed down the event handler")
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

	b.events <- NewDummyEventContainer(DummyEvent3{ v1:"test", v2:15 })
}