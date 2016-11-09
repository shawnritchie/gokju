package routing

import "time"

type Eventer interface {
	EventID() string
}

type EventContainer interface {
	Event() Eventer
	Seq() uint64
	Timestamp() time.Time
}




