package rx

// Observer represents a group of EventHandlers.
type Observer struct {
	NextHandler     NextFunc
	ErrorHandler    ErrorFunc
	CompleteHandler CompleteFunc
}

// DefaultObserver guarantees any handler won't be nil.
var DefaultObserver = Observer{
	NextHandler:     func(interface{}) {},
	ErrorHandler:    func(err error) {},
	CompleteHandler: func() {},
}

// Handle registers Observer to EventHandler.
func (ob Observer) Handle(item interface{}) {
	switch item := item.(type) {
	case error:
		ob.ErrorHandler(item)
		return
	default:
		ob.NextHandler(item)
	}
}

// NewObserver constructs a new Observer instance with default Observer and
// accept any number of EventHandler.
func NewObserver(handlers ...EventHandler) Observer {
	ob := DefaultObserver
	if len(handlers) > 0 {
		for _, handler := range handlers {
			switch handler := handler.(type) {
			case NextFunc:
				ob.NextHandler = handler
			case ErrorFunc:
				ob.ErrorHandler = handler
			case CompleteFunc:
				ob.CompleteHandler = handler
			case Observer:
				ob = handler
			}
		}
	}
	return ob
}

// Next applies Observer's NextHandler to an Item
func (ob Observer) Next(item interface{}) {
	switch item := item.(type) {
	case error:
		return
	default:
		if ob.NextHandler != nil {
			ob.NextHandler.Handle(item)
		}
	}
}

// Error applies Observer's ErrorHandler to an error
func (ob Observer) Error(err error) {
	if ob.ErrorHandler != nil {
		ob.ErrorHandler.Handle(err)
	}
}

// Complete terminates the Observer's internal Observable
func (ob Observer) Complete() {
	if ob.CompleteHandler != nil {
		ob.CompleteHandler.Handle(nil)
	}
}
