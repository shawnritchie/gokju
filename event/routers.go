package event

type BlockingRouter struct {
	Emitter
	Router
}

func NewBlockingRouter(context RouterContext, listener interface{}) *BlockingRouter {
	q := newBlockingQueue()
	r := NewRouter(context, q, listener)
	return &BlockingRouter{
		Emitter: q,
		Router: r,
	}
}