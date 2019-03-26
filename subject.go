package rx

// Subject represents a special type of observable that allows multi-casting.
type Subject struct {
	observers []SimpleObservable
	isStopped bool
}

// NewSubject creates a new subject structure.
func NewSubject() Subject {
	return Subject{}
}

// Next notifies all subscribed observers with the given value.
func (s Subject) Next(value interface{}) {
	if !s.isStopped {
		for _, observer := range s.observers {
			observer.input() <- value
		}
	}
}

// Error notifies all subscribed observers with the given error.
func (s *Subject) Error(err error) {
	s.isStopped = true
	for _, observer := range s.observers {
		observer.input() <- err
	}
}

// Complete terminates this subject and notifies all subscribed observers.
func (s *Subject) Complete() {
	s.isStopped = true
	for _, observer := range s.observers {
		close(observer.input())
	}
}

// Observable returns a new observable with this subject as source.
func (s *Subject) Observable() Observable {
	observable := NewSimpleObservable(nil)
	s.observers = append(s.observers, observable.(SimpleObservable))
	return observable
}
