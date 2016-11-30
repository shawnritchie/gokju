package event

type (
	Interceptor struct {
		Identifier Identifier
		Version int
		in <- chan Event
		out chan <- Event
		Intercept func(Event) Event
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
