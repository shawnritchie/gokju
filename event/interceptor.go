package event

import "reflect"

type (
	Interceptor struct {
		Event reflect.Type
		Version int
		in <- chan Eventer
		out chan <- Eventer
		Intercept func(Eventer) Eventer
	}

	Interceptors []Interceptor
)


func (s Interceptors) Len() int {
	return len(s)
}

func (s Interceptors) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Interceptors) Less(i, j int) bool {
	return s[i].Version < s[j].Version
}
