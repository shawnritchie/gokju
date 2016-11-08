package routing

import (
	"reflect"
	"strings"
)

type Router interface {
	Aggregate() Aggregate
	Handlers() map[reflect.Type]reflect.Value
	EventConsumer() <- chan EventContainer
	EventEmitter() chan <- EventContainer
	Send(e EventContainer)
	SendAndWait(e EventContainer)
}

var eventerType reflect.Type = reflect.TypeOf((*Eventer)(nil)).Elem()

func extractHandlers(agg interface{}) map[reflect.Type]reflect.Value {
	handlers := make(map[reflect.Type]reflect.Value)
	t, v := reflect.TypeOf(agg), reflect.ValueOf(agg)

	for i := 0; i < t.NumMethod(); i++ {
		methodType := t.Method(i)
		methodVal := v.MethodByName(methodType.Name)
		if followsNamingConvention(methodType) && doesConsumeEventer(methodVal.Type()) {
			handlers[methodVal.Type().In(0)] = methodVal
		}
	}

	return handlers
}

func eventRouter(r Router) {
	for {
		select {
		case eventContainer := <-r.EventConsumer():
			m, ok := r.Handlers()[reflect.TypeOf(eventContainer.Event())]
			if (ok) {
				inputs := make([]reflect.Value, m.Type().NumIn())
				inputs[0] = reflect.ValueOf(eventContainer.Event())
				for i := 1; i < m.Type().NumIn(); i++ {
					inputs[i] = reflect.ValueOf(eventContainer.MetaData().Get(m.Type().In(i)))
				}
				m.Call(inputs)
			}
		}

	}
}

func followsNamingConvention(m reflect.Method) bool {
	return strings.HasSuffix(m.Name, "Handler")
}

func doesConsumeEventer(m reflect.Type) bool {
	return (m.NumIn() > 0) && (m.In(0).Implements(eventerType))
}

