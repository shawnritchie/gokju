package routing

type Eventer interface {
	EventID() string
}

type EventContainer interface {
	Event() Eventer
}




