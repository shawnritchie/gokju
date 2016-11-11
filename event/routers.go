package event

type BlockingRouter struct {
	Emitter
	Router
}

func NewBlockingRouter(listener interface{}) *BlockingRouter {
	q := newBlockingQueue()
	r := NewRouter(q, listener)
	return &BlockingRouter{
		Emitter: q,
		Router: r,
	}
}