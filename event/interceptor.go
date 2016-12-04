package event

type (
	Interceptor struct {
		Identifier
		Intercept func(c Container) (Container, error)
	}

	Interceptors []Interceptor
)