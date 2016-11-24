package gokju

import (
	"github.com/shawnritchie/gokju/event"
	"reflect"
	"github.com/shawnritchie/gokju/structs"
	"errors"
	"fmt"
)

type application struct {
	routerContext event.DefaultContainerContext
	aggregateFactory map[reflect.Type]event.AggregateFactory
	aggregateInstances map[reflect.Type]map[event.AggregateID]event.Router
}

var app application = application{
	routerContext:(event.DefaultContainerContext{}),
	aggregateFactory:make(map[reflect.Type]event.AggregateFactory),
	aggregateInstances:make(map[reflect.Type]map[event.AggregateID]event.Router),
}

func RegisterMetaDataDefinition(defintion structs.MetaDataDefinition) {
	app.routerContext = event.NewContainerContext(defintion)
}

func RegisterAggregateFactory(aggregateType reflect.Type, aggregateFactory event.AggregateFactory) {
	if _, ok := app.aggregateInstances[aggregateType]; !ok {
		app.aggregateInstances[aggregateType] = make(map[event.AggregateID]event.Router)
	}

	app.aggregateFactory[aggregateType] = func(id event.AggregateID) event.Aggregate {
		agg := aggregateFactory(id)
		//TODO:: ADD REPLAYING STATE SOME TIME
		app.aggregateInstances[aggregateType][agg.AggregateID()] = event.NewBlockingHandler(app.routerContext, agg)
		return agg
	}
}

func fetchAggregate(aggregateType reflect.Type, id event.AggregateID) (event.Aggregate, error) {
	if factory, ok := app.aggregateFactory[aggregateType]; ok {
		if instances, ok := app.aggregateInstances[aggregateType]; ok {
			if instance, ok := instances[id]; ok {
				return instance.Listener.(event.Aggregate), nil
			}
		}
		return factory(id), nil
	}
	return nil, errors.New(fmt.Sprintf("Unsupported Aggregate Type: %v", aggregateType))
}




