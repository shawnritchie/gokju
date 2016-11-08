package main

import (
	"github.com/shawnritchie/gokju/routing"
	"fmt"
	"time"
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

type DummyEventContainer struct {
	event     routing.Eventer
	metaData  structs.MetaData
}

func (c DummyEventContainer)Event() routing.Eventer{
	return c.event
}

func (c DummyEventContainer)MetaData() structs.MetaData {
	return c.metaData
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

	m := structs.MetaData{}
	m.Add(uint64(1))
	m.Add(time.Now())

	aggregate.SendAndWait(DummyEventContainer{
		event: DummyEvent1{ v1:"test", v2:15 },
		metaData: m,
	})

	aggregate.SendAndWait(DummyEventContainer{
		event: DummyEvent2{ v1:"test", v2:15, close:close },
		metaData: m,
	})

	<- close
}