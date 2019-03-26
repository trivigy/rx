package rx

type Subject struct {
	observers []SimpleObservable
	isStopped bool
}

func NewSubject() Subject {
	return Subject{}
}

func (s Subject) Next(value interface{}) {
	if !s.isStopped {
		for _, observer := range s.observers {
			observer.input() <- value
		}
	}
}

func (s *Subject) Error(err error) {
	s.isStopped = true
	for _, observer := range s.observers {
		observer.input() <- err
	}
}

func (s *Subject) Complete() {
	s.isStopped = true
	for _, observer := range s.observers {
		close(observer.input())
	}
}

func (s *Subject) Observable() Observable {
	observable := NewSimpleObservable(nil)
	s.observers = append(s.observers, observable.(SimpleObservable))
	return observable
}
