package event

import (
	"reflect"
	"strings"
)

type Processor interface {
	Addressable
	Consumer
	Emitter
	Handle(reflect.Type) func(c Container)
}

type BlockingEventProcessor struct {
	Addressable
	Consumer
	Emitter
	handlers map[reflect.Type]func(c Container)
}

func NewBlockingProcessor(context ContainerContext, listener Addressable) *BlockingEventProcessor {
	q := NewBlockingQueue()
	b := &BlockingEventProcessor{
		Addressable: listener,
		Consumer: &q,
		Emitter: &q,
	}
	b.handlers = extractHandlers(context, listener)
	go process(b)
	return b
}

func (s *BlockingEventProcessor)Handle(t reflect.Type) func(c Container) {
	if fx, prs := s.handlers[t]; prs {
		return fx
	} else {
		return func(c Container) {}
	}
}


func extractHandlers(context ContainerContext, listener interface{}) map[reflect.Type]func(c Container) {
	handlers := make(map[reflect.Type]func(c Container))
	t, v := reflect.TypeOf(listener), reflect.ValueOf(listener)

	for i := 0; i < t.NumMethod(); i++ {
		methodType := t.Method(i)
		def := matchesEventHandlerDefinitions(context, methodType.Type)
		if def != nil && followsNamingConvention(methodType) {
			methodVal := v.MethodByName(methodType.Name)
			handlers[methodVal.Type().In(0)] = func(c Container) {
				v := methodVal
				d := def
				context.MapFunctionType(d)(v, c)
			}
		}
	}

	return handlers
}

func process(r Processor) {
	for c := range r.Consume() {
		r.Handle(reflect.TypeOf(c.Event))(c)
	}
}

func matchesEventHandlerDefinitions(context ContainerContext, t reflect.Type) reflect.Type {
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

