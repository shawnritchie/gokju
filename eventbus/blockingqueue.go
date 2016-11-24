package eventbus

import (
	"github.com/shawnritchie/gokju/event"
	"time"
)

type blockingQueue struct {
	queue chan event.Container
}

func (s *blockingQueue)ChannelIn() chan <- event.Container {
	return s.queue
}

func (s *blockingQueue)ChannelOut() <- chan event.Container {
	return s.queue
}

func (s *blockingQueue)Send(e event.Container) {
	go func() {
		s.queue <- e
	}()
}

func (s *blockingQueue)SendAck(e event.Container, d time.Duration, ack func(event.Container), fail func(event.Container)) {
	go func() {
		select {
		case s.queue <- e:
			ack(e)
		case <-time.After(d):
			fail(e)
		}
	}()
}

func (s *blockingQueue)SendAndWait(e event.Container) {
	s.queue <- e
}

func (s *blockingQueue)SendAndWaitWithTimeout(e event.Container, d time.Duration) (succ bool) {
	select {
	case s.queue <- e:
		succ = true
	case <-time.After(d):
		succ = false
	}
	return succ
}

func NewBlockingQueue() blockingQueue {
	return blockingQueue{queue:make(chan event.Container)}
}