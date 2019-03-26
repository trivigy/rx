package rx

type (
	// NextFunc handles a next item in a stream.
	NextFunc func(interface{})

	// ErrorFunc handles an error in a stream.
	ErrorFunc func(error)

	// CompleteFunc handles the end of a stream.
	CompleteFunc func()
)

// Handle registers NextFunc to EventHandler.
func (handle NextFunc) Handle(item interface{}) {
	switch item := item.(type) {
	case error:
		return
	default:
		handle(item)
	}
}

// Handle registers ErrorFunc to EventHandler.
func (handle ErrorFunc) Handle(item interface{}) {
	switch item := item.(type) {
	case error:
		handle(item)
	default:
		return
	}
}

// Handle registers CompleteFunc to EventHandler.
func (handle CompleteFunc) Handle(item interface{}) {
	handle()
}
