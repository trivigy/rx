package rx

import (
	"sync"
)

// SimpleObservable represents an object that can emit a set of values over a
// period of time.
type SimpleObservable struct {
	baseObservable
}

// NewSimpleObservable creates new simple observable.
func NewSimpleObservable(source Observable) Observable {
	observable := SimpleObservable{
		baseObservable: baseObservable{
			_source:    source,
			_input:     make(chan interface{}),
			register:   make(chan chan<- interface{}),
			unregister: make(chan chan<- interface{}),
			outputs:    make(map[chan<- interface{}]*outputMetadata),
			mutex:      &sync.RWMutex{},
		},
	}
	go observable.run()
	return observable
}

// Subscribe registers Observer handlers for notifications it will emit.
func (o SimpleObservable) Subscribe(handlers ...EventHandler) *Subscription {
	return o.subscribe(true, handlers...)
}

// Pipe is used to stitch together functional operators into a chain.
func (o SimpleObservable) Pipe(operations ...OperatorFunction) Observable {
	var observable Observable = o
	for _, operation := range operations {
		observable = operation(observable)
	}
	return observable
}

// ToPromise returns a awaitable channel which closes after receiving one
// notification.
func (o SimpleObservable) ToPromise() func() chan interface{} {
	channel := make(chan interface{})
	sub := o.Subscribe(
		NextFunc(func(item interface{}) {
			channel <- item
		}),
		ErrorFunc(func(err error) {
			channel <- err
		}),
		CompleteFunc(func() {
		}),
	)

	var closed bool
	return func() chan interface{} {
		output := make(chan interface{})
		func() {
			value, ok := <-channel
			if ok {
				var wg sync.WaitGroup
				wg.Add(1)
				go func() {
					defer func() {
						recover()
					}()
					wg.Done()
					output <- value
				}()
				wg.Wait()
				close(channel)
				sub.Unsubscribe()
			}
		}()

		if !closed {
			closed = true
			return output
		}
		close(output)
		return output
	}
}
