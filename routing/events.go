package routing

import (
	"github.com/shawnritchie/gokju/structs"
)

type Eventer interface {
	EventID() string
}

type EventContainer interface {
	Event() Eventer
	MetaData() structs.MetaData
}




