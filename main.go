package main

import (
	"github.com/shawnritchie/gokju/event"
	"fmt"
	"time"
	"context"
	"github.com/shawnritchie/gokju/structs"
	"reflect"
)

type metaKey int

const (
	seqKey metaKey = iota
	timestampKey
	stringKey
)

func (k metaKey)ToInt() int {
	return int(k)
}

var metaKeyDef structs.MetaDataDefinition = structs.MetaDataDefinition {
	Keys: structs.MetaDataMap{seqKey: reflect.TypeOf((*uint64)(nil)).Elem(),
		timestampKey: reflect.TypeOf((*time.Time)(nil)).Elem(),
		stringKey: reflect.TypeOf((*string)(nil)).Elem(),
	},
	Generator:func(i int) structs.MetaDataIdentifier{
		return metaKey(i)
	},
}


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

	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(10 * time.Second))
	ctx = context.WithValue(ctx, "test", "test")
	ctx = context.WithValue(ctx, "test2", "test2")

	eventListener := EventListener{v1: "", v2: 0}
	routerContext := event.NewRouterContext(metaKeyDef)
	eventListener.BlockingRouter = event.NewBlockingRouter(routerContext, &eventListener)

	eventListener.SendAndWait(event.EventContainer{
		Event: DummyEvent1{ v1:"test", v2:15 },
		MetaData: structs.MetaData{
			seqKey: uint64(1),
			timestampKey: time.Now(),
		},
	})

	eventListener.SendAndWait(event.EventContainer{
		Event: DummyEvent2{ v1:"test", v2:15, close:close },
		MetaData: structs.MetaData{
			seqKey: uint64(2),
			timestampKey: time.Now(),
		},
	})

	<- close
}