package rx

type subscribed struct{}

type unsubscribed struct{}

type closed struct{}

type message struct {
	ch chan<- interface{}
	data interface{}
}
