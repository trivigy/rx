package rx

type EventHandler interface {
	Handle(interface{})
}
