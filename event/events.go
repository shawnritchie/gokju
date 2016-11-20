package event

import (
	"time"
	"github.com/shawnritchie/gokju/structs"
)

type Eventer interface {
	EventID() string
}


type EventContainer struct {
	Event    Eventer
	MetaData structs.MetaData
}

type Emitter interface {
	channelIn() chan <- EventContainer
	Send(e EventContainer)
	SendAndWait(e EventContainer)
	SendAndWaitWithTimeout(e EventContainer, d time.Duration) (succ bool)
}

type Consumer interface {
	channelOut() <- chan EventContainer
}


type blockingQueue struct {
	queue chan EventContainer
}

func (s *blockingQueue)channelIn() chan <- EventContainer {
	return s.queue
}

func (s *blockingQueue)channelOut() <- chan EventContainer {
	return s.queue
}

func (s *blockingQueue)Send(e EventContainer) {
	go func() {
		s.queue <- e
	}()
}

func (s *blockingQueue)SendAck(e EventContainer, d time.Duration, ack func(EventContainer), fail func(EventContainer)) {
	go func() {
		select {
		case s.queue <- e:
			ack(e)
		case <-time.After(d):
			fail(e)
		}
	}()
}

func (s *blockingQueue)SendAndWait(e EventContainer) {
	s.queue <- e
}

func (s *blockingQueue)SendAndWaitWithTimeout(e EventContainer, d time.Duration) (succ bool) {
	select {
	case s.queue <- e:
		succ = true
	case <-time.After(d):
		succ = false
	}
	return succ
}

func newBlockingQueue() *blockingQueue {
	return &blockingQueue{queue:make(chan EventContainer)}
}