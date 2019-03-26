package rx

type OperatorFunction func(Observable) Observable

func Filter(fn func(interface{}) bool) OperatorFunction {
	return func(observable Observable) Observable {
		switch observable := observable.(type) {
		case SimpleObservable:
			nextObservable := NewSimpleObservable(observable)
			observable.Subscribe(
				NextFunc(func(item interface{}) {
					switch item := item.(type) {
					case message:
						if fn(item.data) {
							nextObservable.input() <- item
						}
					default:
						if fn(item) {
							nextObservable.input() <- item
						}
					}
				}),
				ErrorFunc(func(err error) {
					nextObservable.input() <- err
				}),
				CompleteFunc(func() {
					close(nextObservable.input())
				}),
			)
			return nextObservable
		case BehaviorObservable:
			nextObservable := NewBehaviorObservable(
				observable,
				observable.currentValue,
			).(BehaviorObservable)
			observable.subscribe(false,
				NextFunc(func(item interface{}) {
					switch item := item.(type) {
					case message:
						if fn(item.data) {
							nextObservable.SetCurrentValue(item.data)
							nextObservable.input() <- item
						}
					default:
						if fn(item) {
							nextObservable.SetCurrentValue(item)
							nextObservable.input() <- item
						}
					}
				}),
				ErrorFunc(func(err error) {
					nextObservable.input() <- err
				}),
				CompleteFunc(func() {
					close(nextObservable.input())
				}),
			)
			return nextObservable
		default:
			panic("should not happen")
		}
	}
}

func Take(num interface{}) OperatorFunction {
	return func(observable Observable) Observable {
		counter := num.(int)
		switch observable := observable.(type) {
		case SimpleObservable:
			nextObservable := NewSimpleObservable(observable)
			observable.Subscribe(
				NextFunc(func(item interface{}) {
					if counter > 0 {
						nextObservable.input() <- item
						counter--
					}
				}),
				ErrorFunc(func(err error) {
					nextObservable.input() <- err
				}),
				CompleteFunc(func() {
					close(nextObservable.input())
				}),
			)
			return nextObservable
		case BehaviorObservable:
			nextObservable := NewBehaviorObservable(
				observable,
				observable.currentValue,
			).(BehaviorObservable)
			observable.subscribe(false,
				NextFunc(func(item interface{}) {
					if counter > 0 {
						nextObservable.SetCurrentValue(observable.currentValue)
						nextObservable.input() <- item
						counter--
					}
				}),
				ErrorFunc(func(err error) {
					nextObservable.input() <- err
				}),
				CompleteFunc(func() {
					close(nextObservable.input())
				}),
			)
			return nextObservable
		default:
			panic("should not happen")
		}
	}
}
