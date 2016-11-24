package ddd

import (
	"github.com/shawnritchie/gokju/aggregate"
	"github.com/shawnritchie/gokju/eventbus"
)

type Aggregate interface {
	Address() eventbus.Address
	AggregateID() aggregate.Identifier
}

type AggregateFactory func(id aggregate.Identifier) Aggregate
