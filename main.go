package main

import (
	"github.com/shawnritchie/gokju/routing"
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
	event     routing.Eventer
	SeqNo     uint64
	Timestamp time.Time
}

func (c DummyEventContainer)Event() routing.Eventer{
	return c.event
}

type Aggregate struct {
	*routing.BlockingRouter
	v1 string
	v2 int
}

func (a *Aggregate)Router() routing.Router {
	return a.BlockingRouter
}

func (a *Aggregate)DummyEvent1Handler(event DummyEvent1, seqNo uint64, timestamp time.Time) {
	fmt.Printf("Event:%v, seqNo:%v, timestamp:%v\n", event, seqNo, timestamp)
	a.v1 = event.v1
	a.v2 = event.v2
}

func (a *Aggregate)DummyEvent2Handler(event DummyEvent2, seqNo uint64, timestamp time.Time) {
	fmt.Printf("Event:%v, seqNo:%v, timestamp:%v\n", event, seqNo, timestamp)
	a.v1 = event.v1
	a.v2 = event.v2
	event.close <- struct{}{}
}


func main() {
	close := make(chan struct{})

	aggregate := Aggregate{v1: "", v2: 0}
	aggregate.BlockingRouter = routing.NewBlockingRouter(&aggregate)

	aggregate.SendAndWait(DummyEventContainer{
		event: DummyEvent1{ v1:"test", v2:15 },
		SeqNo: uint64(1),
		Timestamp: time.Now(),
	})

	aggregate.SendAndWait(DummyEventContainer{
		event: DummyEvent2{ v1:"test", v2:15, close:close },
		SeqNo: uint64(2),
		Timestamp: time.Now(),
	})

	<- close
}