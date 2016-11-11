package event

import (
	"reflect"
	"strings"
	"time"
)

type Router struct {
	Consumer
	listener interface{}
	handlers map[reflect.Type]func(c EventContainer)
}

func NewRouter(c Consumer, listener interface{}) Router {
	b := Router{
		listener: listener,
		Consumer: c,
	}
	b.handlers = extractHandlers(b.listener)
	go eventRouter(b)
	return b
}

var eventHandlerTemplates map[reflect.Type]func(v reflect.Value, c EventContainer) = map[reflect.Type]func(v reflect.Value, c EventContainer){
	reflect.TypeOf(func(a interface{}, event Eventer){}): func(v reflect.Value, c EventContainer) {
		v.Call([]reflect.Value {
			reflect.ValueOf(c.Event()),
		})
	},
	reflect.TypeOf(func(a interface{}, event Eventer, timestamp time.Time){}): func(v reflect.Value, c EventContainer) {
		v.Call([]reflect.Value {
			reflect.ValueOf(c.Event()),
			reflect.ValueOf(c.Timestamp()),
		})
	},
	reflect.TypeOf(func(a interface{}, event Eventer, timestamp time.Time, seq uint64){}): func(v reflect.Value, c EventContainer) {
		v.Call([]reflect.Value {
			reflect.ValueOf(c.Event()),
			reflect.ValueOf(c.Timestamp()),
			reflect.ValueOf(c.Seq()),
		})
	},
}

func extractHandlers(listener interface{}) map[reflect.Type]func(c EventContainer) {
	handlers := make(map[reflect.Type]func(c EventContainer))
	t, v := reflect.TypeOf(listener), reflect.ValueOf(listener)

	for i := 0; i < t.NumMethod(); i++ {
		methodType := t.Method(i)
		def := matchesEventHandlerDefinitions(methodType.Type)
		if def != nil && followsNamingConvention(methodType) {
			methodVal := v.MethodByName(methodType.Name)
			handlers[methodVal.Type().In(0)] = func(c EventContainer) {
				v := methodVal
				d := def
				eventHandlerTemplates[d](v, c)
			}
		}
	}

	return handlers
}

func eventRouter(r Router) {
	for {
		select {
		case c := <-r.channelOut():
			fx, ok := r.handlers[reflect.TypeOf(c.Event())]
			if (ok) {
				fx(c)
			}
		}

	}
}

func matchesEventHandlerDefinitions(t reflect.Type) reflect.Type {
	for k, _ := range eventHandlerTemplates {
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
