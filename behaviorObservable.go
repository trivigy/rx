package rx

import (
	"sync"
)

// BehaviorObservable represents an object that can emit a set of values over a
// period of time. This object is a special case as opposed to simple observable
// differs in that it emits the last value emitted to new subscribers upon
// subscription.
type BehaviorObservable struct {
	baseObservable
	currentValue interface{}
}

// NewBehaviorObservable creates new behavior observable.
func NewBehaviorObservable(source Observable, initialValue interface{}) Observable {
	observable := BehaviorObservable{
		baseObservable: baseObservable{
			_source:    source,
			_input:     make(chan interface{}),
			register:   make(chan chan<- interface{}),
			unregister: make(chan chan<- interface{}),
			outputs:    make(map[chan<- interface{}]*outputMetadata),
			mutex:      &sync.RWMutex{},
		},
		currentValue: initialValue,
	}
	go observable.run()
	return observable
}

// Subscribe registers Observer handlers for notifications it will emit.
func (o BehaviorObservable) Subscribe(handlers ...EventHandler) *Subscription {
	sub := o.subscribe(true, handlers...)

	var observable Observable = o
	for observable.source() != nil {
		observable = observable.source()
	}

	observable.input() <- message{ch: sub.channel, data: o.currentValue}
	return sub
}

// Pipe is used to stitch together functional operators into a chain.
func (o BehaviorObservable) Pipe(operations ...OperatorFunction) Observable {
	var observable Observable = o
	for _, operation := range operations {
		observable = operation(observable)
	}
	return observable
}

// ToPromise returns a awaitable channel which closes after receiving one
// notification.
func (o BehaviorObservable) ToPromise() func() chan interface{} {
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

// SetCurrentValue sets the current value of the observable which will be
// emitted to the next subscriber. This value mostly adjusts automatically
// based on the last emitted value.
func (o *BehaviorObservable) SetCurrentValue(value interface{}) {
	o.currentValue = value
}
