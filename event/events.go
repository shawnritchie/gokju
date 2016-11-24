package event

import "reflect"

type (
	Identifier string

	Eventer interface {
		EventID() Identifier
	}
)

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
