package handler

import (
"reflect"
"strings"
"github.com/shawnritchie/gokju/eventbus"
	"github.com/shawnritchie/gokju/event"
)

type Simple struct {
	eventbus.Consumer
	Listener interface{}
	handlers map[reflect.Type]func(c event.Container)
}

func NewSimpleHandler(context event.ContainerContext, consumer eventbus.Consumer, listener interface{}) Simple {
	b := Simple{
		Listener: listener,
		Consumer: consumer,
	}
	b.handlers = extractHandlers(context, b.Listener)
	go eventRouter(b)
	return b
}

func (s *Simple)Handle(t reflect.Type) func(c event.Container) {
	return s.handlers[t]
}

func extractHandlers(context event.ContainerContext, listener interface{}) map[reflect.Type]func(c event.Container) {
	handlers := make(map[reflect.Type]func(c event.Container))
	t, v := reflect.TypeOf(listener), reflect.ValueOf(listener)

	for i := 0; i < t.NumMethod(); i++ {
		methodType := t.Method(i)
		def := matchesEventHandlerDefinitions(context, methodType.Type)
		if def != nil && followsNamingConvention(methodType) {
			methodVal := v.MethodByName(methodType.Name)
			handlers[methodVal.Type().In(0)] = func(c event.Container) {
				v := methodVal
				d := def
				context.MapFunctionType(d)(v, c)
			}
		}
	}

	return handlers
}

func eventRouter(r Simple) {
	for {
		select {
		case c := <-r.ChannelOut():
			fx, ok := r.handlers[reflect.TypeOf(c.Event)]
			if (ok) {
				fx(c)
			}
		}

	}
}

func matchesEventHandlerDefinitions(context event.ContainerContext, t reflect.Type) reflect.Type {
	for _, k := range context.AllFunctionType() {
		if k.NumIn() == t.NumIn() {
			s := true;
			for i:= 0; i<k.NumIn(); i++ {
				if k.In(i).Kind() == reflect.Interface {
					s = s && t.In(i).Implements(k.In(i))
				} else {
					s = s && t.In(i)==k.In(i)
				}
			}
			if s {
				return k
			}
		}
	}
	return nil
}

func followsNamingConvention(m reflect.Method) bool {
	return strings.HasSuffix(m.Name, "Handler")
}

