package event

import (
	"sync"
	"reflect"
	"errors"
)

type Address string

type Addressable interface {
	Address() Address
}

type Bus interface {
	Subscribe(p Processor) (unsubscribe func(), err error)
	Register(interceptor Interceptor) error
	Publish(events ...Container)
}

type SimpleEventBus struct {
	mutex sync.RWMutex
	eventProcessors []managedEventProcessor
	interceptors map[Identifier]Interceptors
}

type managedEventProcessor struct {
	destroy chan struct{}
	consume chan Container
	processor Processor
}

func newManagedEventProcessor(processor Processor) managedEventProcessor {
	p := managedEventProcessor{
		destroy: make(chan struct{}),
		consume: make(chan Container),
		processor: processor,
	}
	go p.process()
	return p
}

func (p *managedEventProcessor)process() {
	select {
	case <- p.destroy:
		p.processor.Destroy()
		close(p.destroy)
		close(p.consume)
		return
	case c := <- p.consume:
		p.processor.Handle(reflect.TypeOf(c.Event))(c)
	}
}

func NewSimpleEventBus() SimpleEventBus {
	return SimpleEventBus{
		mutex: sync.RWMutex{},
		eventProcessors: []managedEventProcessor{},
		interceptors: map[Identifier]Interceptors{},
	}
}

func (s *SimpleEventBus)Subscribe(p Processor) (unsubscribe func(), err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if i := s.findProcessor(p); i == -1 {
		m := newManagedEventProcessor(p)
		s.eventProcessors = append(s.eventProcessors, m)

		return func() {
			s.mutex.Lock()
			defer s.mutex.Unlock()

			m.destroy <- struct{}{}
			if i := s.findProcessor(p); i != -1 {
				s.eventProcessors = append(s.eventProcessors[:i], s.eventProcessors[i+1:]...)
			}
		}, nil
	}
	return nil, errors.New("Duplicate Entry")
}

func (s *SimpleEventBus)Register(interceptor Interceptor) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	id := interceptor.Identifier
	if _, prs := s.interceptors[id]; !prs {
		s.interceptors[id] = Interceptors{}
	}

	s.interceptors[id] = append(s.interceptors[id], interceptor)

	return nil
}


func (s *SimpleEventBus)Publish(containers ...Container) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, c := range containers {
		for _, m := range s.eventProcessors {
			go func () {
				container, err := s.intercept(c)
				if (err != nil) {
					//TODO:: LOGGING
				} else {
					m.consume <- container
				}
			}()
		}
	}
}

func (s *SimpleEventBus)intercept(c Container) (Container, error) {
	id := EventIdentifier(c.Event)
	var e error = nil
	container := c
	if _, prs := s.interceptors[id]; prs {
		for _, interceptor := range s.interceptors[id] {
			container, e = interceptor.Intercept(c)
			if (e != nil) {
				return Container{}, e
			}
		}
	}
	return container, nil
}

func(s *SimpleEventBus)findProcessor(p Processor) int {
	for i, m := range s.eventProcessors {
		if (p.Address() == m.processor.Address()) {
			return i
		}
	}
	return -1
}
