package event

import "time"

type Eventer interface {
	EventID() string
}

type EventContainer interface {
	Event() Eventer
	Seq() uint64
	Timestamp() time.Time
}

type Emitter interface {
	channelIn() chan <- EventContainer
	Send(e EventContainer)
	SendAndWait(e EventContainer)
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

func (s *blockingQueue)SendAndWait(e EventContainer) {
	s.queue <- e
}

func newBlockingQueue() *blockingQueue {
	return &blockingQueue{queue:make(chan EventContainer)}
}