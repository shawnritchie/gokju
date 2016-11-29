package event


import (
	"sync"
	"reflect"
	"errors"
	"sort"
	"fmt"
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
	interceptors map[reflect.Type]Interceptors
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
		interceptors: map[reflect.Type]Interceptors{},
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

	t := reflect.TypeOf(interceptor)
	if _, prs := s.interceptors[t]; !prs {
		s.interceptors[t] = []Interceptor{}
	}

	if err := isValidInterceptor(s.interceptors[t], interceptor); err != nil {
		return err
	}

	s.interceptors[t] = append(s.interceptors[t], interceptor)
	sort.Sort(s.interceptors[t])

	return nil
}

func isValidInterceptor(chain []Interceptor, add Interceptor) error {
	if len(chain) == 0 && add.Version != 1 {
		return errors.New("Expected the first interceptor to start with version 1")
	} else if chain[len(chain) - 1].Version+1 != add.Version {
		return errors.New(fmt.Sprintf("Expected Version %v, Received Version %v", chain[len(chain) - 1].Version+1, add.Version))
	}

	return nil
}


func (s *SimpleEventBus)Publish(events ...Container) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, e := range events {
		for _, m := range s.eventProcessors {
			go func () {
				e.Event = s.intercept(e.Event)
				m.consume <- e
			}()
		}
	}
}

func (s *SimpleEventBus)intercept(e Eventer) Eventer {
	t := reflect.TypeOf(e)
	if _, prs := s.interceptors[t]; prs {
		for _, interceptor := range s.interceptors[t] {
			if (interceptor.Version - 1) ==  e.Version() {
				e = interceptor.Intercept(e)
			}
		}
	}
	return e
}

func(s *SimpleEventBus)findProcessor(p Processor) int {
	for i, m := range s.eventProcessors {
		if (p.Address() == m.processor.Address()) {
			return i
		}
	}
	return -1
}
