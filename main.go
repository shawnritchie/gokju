package main

import (
	"github.com/shawnritchie/gokju/event"
	"fmt"
	"time"
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


type DummyEventContainer struct {
	event     event.Eventer
	seq       uint64
	timestamp time.Time
}

func (c DummyEventContainer)Event() event.Eventer{
	return c.event
}

func (c DummyEventContainer)Seq() uint64 {
	return c.seq
}

func (c DummyEventContainer)Timestamp() time.Time{
	return c.timestamp
}


type EventListener struct {
	*event.BlockingRouter
	v1 string
	v2 int
}

func (a *EventListener)DummyEvent1Handler(event DummyEvent1, timestamp time.Time, seqNo uint64) {
	fmt.Printf("Event:%v, seqNo:%v, timestamp:%v\n", event, seqNo, timestamp)
	a.v1 = event.v1
	a.v2 = event.v2
}

func (a *EventListener)DummyEvent2Handler(event DummyEvent2, timestamp time.Time, seqNo uint64) {
	fmt.Printf("Event:%v, seqNo:%v, timestamp:%v\n", event, seqNo, timestamp)
	a.v1 = event.v1
	a.v2 = event.v2
	event.close <- struct{}{}
}


func main() {
	close := make(chan struct{})

	eventListener := EventListener{v1: "", v2: 0}
	eventListener.BlockingRouter = event.NewBlockingRouter(&eventListener)

	eventListener.SendAndWait(DummyEventContainer{
		event: DummyEvent1{ v1:"test", v2:15 },
		seq: uint64(1),
		timestamp: time.Now(),
	})

	eventListener.SendAndWait(DummyEventContainer{
		event: DummyEvent2{ v1:"test", v2:15, close:close },
		seq: uint64(2),
		timestamp: time.Now(),
	})

	<- close
}