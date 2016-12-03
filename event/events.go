package event

import (
	"reflect"
	"time"
	"strconv"
)

type (
	Identifier string
	Version int

	Event interface{}
)

func extractUnderlyingType(value interface{}) reflect.Type {
	t := reflect.TypeOf(value)
	if (t.Kind() == reflect.Ptr) {
		return t.Elem()
	}
	return t
}

func EventIdentifier(e interface{}) Identifier {
	t := extractUnderlyingType(e)
	f, ok := t.FieldByName("Event")
	if (ok) {
		return Identifier(f.Tag.Get("id"))
	}
	return Identifier(t.String())
}

func EventVersion(e interface{}) Version {
	t := extractUnderlyingType(e)
	f, ok := t.FieldByName("Event")
	if (ok) {
		i, err := strconv.Atoi(f.Tag.Get("v"))
		if err == nil {
			return Version(i)
		}
	}
	return Version(0)
}

type (
	MetaDataIdentifier interface {
		ToInt() int
	}

	MetaData map[MetaDataIdentifier]interface{}

	MetaDataMap map[MetaDataIdentifier]reflect.Type
)

type MetaDataDefinition struct {
	Keys MetaDataMap
	Generator func(i int) MetaDataIdentifier
}

func (t MetaDataDefinition)Type(i MetaDataIdentifier) reflect.Type {
	return t.Keys[i]
}

func (t MetaDataDefinition)Get(i int) MetaDataIdentifier {
	return t.Generator(i)
}

func (t MetaDataDefinition)Len() int {
	return len(t.Keys)
}

type (
	Emitter interface {
		Emit() chan <- Container
		Send(e Container)
		SendAndWait(e Container)
		SendAndWaitWithTimeout(e Container, d time.Duration) (succ bool)
		Destroy()
	}

	Consumer interface {
		Consume() <- chan Container
	}
)