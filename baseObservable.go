package rx

import (
	"fmt"
	"sync"
	"time"
)

type outputMetadata struct {
	ch   chan<- interface{}
	seq  uint64
	done uint64
}

type baseObservable struct {
	_source    Observable
	_input     chan interface{}
	register   chan chan<- interface{}
	unregister chan chan<- interface{}
	outputs    map[chan<- interface{}]*outputMetadata
	mutex      *sync.RWMutex
}

func (o *baseObservable) push(event interface{}, ch chan<- interface{}, seq uint64) {
	o.mutex.RLock()
	meta, ok := o.outputs[ch]
	o.mutex.RUnlock()
	if !ok {
		return
	}

	if meta.done+1 < seq {
		time.Sleep(time.Duration(seq-meta.done-1) * time.Millisecond)
		go o.push(event, ch, seq)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%#v\n", err)
		}
		meta.done++
	}()
	ch <- event
}

func (o *baseObservable) run() {
	for {
		select {
		case event, ok := <-o._input:
			if ok {
				o.mutex.RLock()
				for ch, meta := range o.outputs {
					go o.push(event, ch, meta.seq)
					meta.seq++
				}
				o.mutex.RUnlock()
			} else {
				for ch := range o.outputs {
					o.mutex.Lock()
					delete(o.outputs, ch)
					o.mutex.Unlock()
					ch <- closed{}
				}
			}
		case ch, ok := <-o.register:
			if ok {
				o.mutex.Lock()
				o.outputs[ch] = &outputMetadata{ch: ch, seq: 1}
				o.mutex.Unlock()
				ch <- subscribed{}
			}
		case ch := <-o.unregister:
			o.mutex.Lock()
			delete(o.outputs, ch)
			o.mutex.Unlock()
			ch <- unsubscribed{}
		}
	}
}

func (o baseObservable) source() Observable {
	return o._source
}

func (o baseObservable) input() chan interface{} {
	return o._input
}

func (o baseObservable) subscribe(leaf bool, handlers ...EventHandler) *Subscription {
	return NewSubscription(o, leaf, handlers...)
}
