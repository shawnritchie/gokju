package handler

import (
	"github.com/shawnritchie/gokju/eventbus"
	"reflect"
	"github.com/shawnritchie/gokju/event"
)

type Handler interface {
	eventbus.Consumer
	Handle(reflect.Type) func(c event.Container)
}

type BlockingHandler struct {
	eventbus.Emitter
	Handler
}

func NewBlockingHandler(context event.ContainerContext, listener interface{}) *BlockingHandler {
	q := eventbus.NewBlockingQueue()
	h := NewSimpleHandler(context, &q, listener)
	return &BlockingHandler{
		Emitter: &q,
		Handler: &h,
	}
}
