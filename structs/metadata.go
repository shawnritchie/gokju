package structs

import "reflect"

type MetaData map[reflect.Type]interface{}

func (m MetaData) Add(i interface{}) {
	m[reflect.TypeOf(i)] = i
}

func (m MetaData) Get(t reflect.Type) interface{} {
	i, ok := m[t]
	if (!ok) {
		return reflect.Zero(t)
	}
	return i
}