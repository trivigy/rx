package rx

import (
	"sync"
)

type BehaviorObservable struct {
	baseObservable
	currentValue interface{}
}

func NewBehaviorObservable(source Observable, initialValue interface{}) Observable {
	observable := BehaviorObservable{
		baseObservable: baseObservable{
			_source:    source,
			_input:     make(chan interface{}),
			register:   make(chan chan<- interface{}),
			unregister: make(chan chan<- interface{}),
			outputs:    make(map[chan<- interface{}]*outputMetadata),
		},
		currentValue: initialValue,
	}
	go observable.run()
	return observable
}

func (o BehaviorObservable) Subscribe(handlers ...EventHandler) *Subscription {
	sub := o.subscribe(true, handlers...)

	var observable Observable = o
	for observable.source() != nil {
		observable = observable.source()
	}

	observable.input() <- message{ch: sub.channel, data: o.currentValue}
	return sub
}

func (o BehaviorObservable) Pipe(operations ...OperatorFunction) Observable {
	var observable Observable = o
	for _, operation := range operations {
		observable = operation(observable)
	}
	return observable
}

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

func (o *BehaviorObservable) SetCurrentValue(value interface{}) {
	o.currentValue = value
}
