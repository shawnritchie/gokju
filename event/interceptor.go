package event

type (
	Interceptor struct {
		Identifier
		Intercept func(E Container) (Container, error)
	}

	Interceptors []Interceptor
)