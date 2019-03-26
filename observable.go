package rx

type Observable interface {
	source() Observable
	input() chan interface{}

	Subscribe(handlers ...EventHandler) *Subscription
	Pipe(operations ...OperatorFunction) Observable
	ToPromise() func() chan interface{}
}
