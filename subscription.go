package rx

import (
	"sync"
	"time"
)

// Subscription is usually returned from any subscription
type Subscription struct {
	observable baseObservable
	channel    chan interface{}
	observer   Observer
	leaf       bool

	subLock   sync.WaitGroup
	unsubLock sync.WaitGroup

	SubscribeAt   time.Time
	UnsubscribeAt time.Time
	Error         error
}

// NewSubscription creates a new subscription and attaches handler functions to
// the specified observable.
func NewSubscription(observable baseObservable, leaf bool, handlers ...EventHandler) *Subscription {
	sub := &Subscription{
		leaf:       leaf,
		observable: observable,
		channel:    make(chan interface{}),
		observer:   NewObserver(handlers...),
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		for item := range sub.channel {
			switch item := item.(type) {
			case subscribed:
				sub.subLock.Done()
			case unsubscribed:
				sub.unsubLock.Done()
			case closed:
				close(sub.channel)
			case error:
				sub.observer.Error(item)
				sub.Error = item
				go sub.Unsubscribe()
			case message:
				if !sub.leaf {
					sub.observer.Next(item)
					continue
				}

				if sub.channel == item.ch {
					sub.observer.Next(item.data)
				}
			default:
				sub.observer.Next(item)
			}
		}

		if sub.Error == nil {
			sub.observer.Complete()
		}
	}()

	wg.Wait()
	sub.subscribe()
	return sub
}

// subscribe records the time of subscription.
func (s *Subscription) subscribe() {
	s.subLock.Add(1)
	s.observable.register <- s.channel
	s.SubscribeAt = time.Now()
	s.subLock.Wait()
}

// Unsubscribe records the time of unsubscription.
func (s *Subscription) Unsubscribe() {
	s.unsubLock.Add(1)
	s.observable.unregister <- s.channel
	s.UnsubscribeAt = time.Now()
	s.unsubLock.Wait()
	close(s.channel)
}
