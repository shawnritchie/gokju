package event

import (
	"time"
)

type blockingQueue struct {
	queue chan Container
}

func (s *blockingQueue)Emit() chan <- Container {
	return s.queue
}

func (s *blockingQueue)Consume() <- chan Container {
	return s.queue
}

func (s *blockingQueue)Send(e Container) {
	go func() {
		s.queue <- e
	}()
}

func (s *blockingQueue)SendAck(e Container, d time.Duration, ack func(Container), fail func(Container)) {
	go func() {
		select {
		case s.queue <- e:
			ack(e)
		case <-time.After(d):
			fail(e)
		}
	}()
}

func (s *blockingQueue)SendAndWait(e Container) {
	s.queue <- e
}

func (s *blockingQueue)SendAndWaitWithTimeout(e Container, d time.Duration) (succ bool) {
	select {
	case s.queue <- e:
		succ = true
	case <-time.After(d):
		succ = false
	}
	return succ
}

func (s *blockingQueue)Destroy() {
	close(s.queue)
}

func NewBlockingQueue() blockingQueue {
	return blockingQueue{queue:make(chan Container)}
}