package eventbus

import (
	"time"
	"github.com/shawnritchie/gokju/event"
)

type Address string


type (
	Emitter interface {
		ChannelIn() chan <- event.Container
		Send(e event.Container)
		SendAndWait(e event.Container)
		SendAndWaitWithTimeout(e event.Container, d time.Duration) (succ bool)
	}

	Consumer interface {
		ChannelOut() <- chan event.Container
	}
)