package rx

type BehaviorSubject struct {
	observers    []BehaviorObservable
	currentValue interface{}
	isStopped    bool
}

func NewBehaviorSubject(initialValue interface{}) BehaviorSubject {
	return BehaviorSubject{
		currentValue: initialValue,
	}
}

func (s *BehaviorSubject) Next(value interface{}) {
	if !s.isStopped {
		s.currentValue = value
		for _, observer := range s.observers {
			observer.SetCurrentValue(s.currentValue)
			observer.input() <- s.currentValue
		}
	}
}

func (s *BehaviorSubject) Error(err error) {
	s.isStopped = true
	for _, observer := range s.observers {
		observer.input() <- err
	}
}

func (s *BehaviorSubject) Complete() {
	s.isStopped = true
	for _, observer := range s.observers {
		close(observer.input())
	}
}

func (s *BehaviorSubject) Observable() Observable {
	observable := NewBehaviorObservable(nil, s.currentValue)
	s.observers = append(s.observers, observable.(BehaviorObservable))
	return observable
}
