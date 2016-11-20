package event

import (
	"reflect"
	"strings"
)

type Router struct {
	RouterContext
	Consumer
	listener interface{}
	handlers map[reflect.Type]func(c EventContainer)
}

func NewRouter(context RouterContext, consumer Consumer, listener interface{}) Router {
	b := Router{
		RouterContext: context,
		listener: listener,
		Consumer: consumer,
	}
	b.handlers = extractHandlers(context, b.listener)
	go eventRouter(b)
	return b
}

func extractHandlers(context RouterContext, listener interface{}) map[reflect.Type]func(c EventContainer) {
	handlers := make(map[reflect.Type]func(c EventContainer))
	t, v := reflect.TypeOf(listener), reflect.ValueOf(listener)

	for i := 0; i < t.NumMethod(); i++ {
		methodType := t.Method(i)
		def := matchesEventHandlerDefinitions(context, methodType.Type)
		if def != nil && followsNamingConvention(methodType) {
			methodVal := v.MethodByName(methodType.Name)
			handlers[methodVal.Type().In(0)] = func(c EventContainer) {
				v := methodVal
				d := def
				context.methodDefinitions[d](v, c)
			}
		}
	}

	return handlers
}

func eventRouter(r Router) {
	for {
		select {
		case c := <-r.channelOut():
			fx, ok := r.handlers[reflect.TypeOf(c.Event)]
			if (ok) {
				fx(c)
			}
		}

	}
}

func matchesEventHandlerDefinitions(context RouterContext, t reflect.Type) reflect.Type {
	for k, _ := range context.methodDefinitions {
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

