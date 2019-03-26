package rx

// Observable represents an interface template for the public behavior of an
// observable.
type Observable interface {
	source() Observable
	input() chan interface{}

	Subscribe(handlers ...EventHandler) *Subscription
	Pipe(operations ...OperatorFunction) Observable
	ToPromise() func() chan interface{}
}
