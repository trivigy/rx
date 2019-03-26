package rx

// BehaviorSubject is a variant of the simple subject struct that requires an
// initial value and emits its current value whenever it is subscribed to.
type BehaviorSubject struct {
	observers    []BehaviorObservable
	currentValue interface{}
	isStopped    bool
}

// NewBehaviorSubject creates a new behavior subject structure.
func NewBehaviorSubject(initialValue interface{}) BehaviorSubject {
	return BehaviorSubject{
		currentValue: initialValue,
	}
}

// Next notifies all subscribed observers with the given value.
func (s *BehaviorSubject) Next(value interface{}) {
	if !s.isStopped {
		s.currentValue = value
		for _, observer := range s.observers {
			observer.SetCurrentValue(s.currentValue)
			observer.input() <- s.currentValue
		}
	}
}

// Error notifies all subscribed observers with the given error.
func (s *BehaviorSubject) Error(err error) {
	s.isStopped = true
	for _, observer := range s.observers {
		observer.input() <- err
	}
}

// Complete terminates this subject and notifies all subscribed observers.
func (s *BehaviorSubject) Complete() {
	s.isStopped = true
	for _, observer := range s.observers {
		close(observer.input())
	}
}

// Observable returns a new observable with this subject as source.
func (s *BehaviorSubject) Observable() Observable {
	observable := NewBehaviorObservable(nil, s.currentValue)
	s.observers = append(s.observers, observable.(BehaviorObservable))
	return observable
}
