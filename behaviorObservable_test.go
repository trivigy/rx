package rx

import (
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BehaviorObservableSuite struct {
	suite.Suite
}

func (s *BehaviorObservableSuite) TestNewBehaviorObservable_WithObserver() {
	wg := sync.WaitGroup{}
	watcher1 := Observer{
		NextHandler: func(item interface{}) {
			assert.Equal(s.T(), item, "hello world")
			wg.Done()
		},
		ErrorHandler: func(err error) {
			assert.Fail(s.T(), err.Error())
		},
		CompleteHandler: func() {
			wg.Done()
		},
	}

	watcher2 := Observer{
		NextHandler: func(item interface{}) {
			assert.Equal(s.T(), item, "hello world")
			wg.Done()
		},
		ErrorHandler: func(err error) {
			assert.Fail(s.T(), err.Error())
		},
		CompleteHandler: func() {
			wg.Done()
		},
	}

	wg.Add(4)
	observable := NewBehaviorObservable(nil, "hello world")
	sub1 := observable.Subscribe(watcher1)
	sub2 := observable.Subscribe(watcher2)

	observable.input() <- "hello world"
	wg.Wait()

	wg.Add(2)
	sub1.Unsubscribe()
	sub2.Unsubscribe()
	wg.Wait()
}

func (s *BehaviorObservableSuite) TestNewBehaviorObservable_WithHandlers() {
	wg := sync.WaitGroup{}

	wg.Add(4)
	observable := NewBehaviorObservable(nil, "hello world")
	sub1 := observable.Subscribe(
		NextFunc(func(item interface{}) {
			assert.Equal(s.T(), item, "hello world")
			wg.Done()
		}),
		ErrorFunc(func(err error) {
			assert.Fail(s.T(), err.Error())
		}),
		CompleteFunc(func() {
			wg.Done()
		}),
	)
	sub2 := observable.Subscribe(
		NextFunc(func(item interface{}) {
			assert.Equal(s.T(), item, "hello world")
			wg.Done()
		}),
		ErrorFunc(func(err error) {
			assert.Fail(s.T(), err.Error())
		}),
		CompleteFunc(func() {
			wg.Done()
		}),
	)

	observable.input() <- "hello world"
	wg.Wait()

	wg.Add(2)
	sub1.Unsubscribe()
	sub2.Unsubscribe()
	wg.Wait()
}

func (s *BehaviorObservableSuite) TestNewBehaviorObservable_Error() {
	wg := sync.WaitGroup{}

	wg.Add(4)
	observable := NewBehaviorObservable(nil, "hello world")
	observable.Subscribe(
		NextFunc(func(item interface{}) {
			assert.Equal(s.T(), item, "hello world")
			wg.Done()
		}),
		ErrorFunc(func(err error) {
			wg.Done()
		}),
		CompleteFunc(func() {
			assert.Fail(s.T(), "executed complete")

		}),
	)
	observable.Subscribe(
		NextFunc(func(item interface{}) {
			assert.Equal(s.T(), item, "hello world")
			wg.Done()
		}),
		ErrorFunc(func(err error) {
			wg.Done()
		}),
		CompleteFunc(func() {
			assert.Fail(s.T(), "executed complete")
		}),
	)

	observable.input() <- errors.New("failure")
	wg.Wait()
}

func (s *BehaviorObservableSuite) TestNewBehaviorObservable_UnsubscribedAfterError() {
	wg := sync.WaitGroup{}

	wg.Add(4)
	observable := NewBehaviorObservable(nil, "hello world")
	observable.Subscribe(
		NextFunc(func(item interface{}) {
			assert.Equal(s.T(), item, "hello world")
			wg.Done()
		}),
		ErrorFunc(func(err error) {
			wg.Done()
		}),
		CompleteFunc(func() {
			assert.Fail(s.T(), "executed complete")

		}),
	)
	observable.input() <- errors.New("failure")
	// First subscriber unsubscribes automatically at this point.

	observable.Subscribe(
		NextFunc(func(item interface{}) {
			assert.Equal(s.T(), item, "hello world")
			wg.Done()
		}),
		ErrorFunc(func(err error) {
			wg.Done()
		}),
		CompleteFunc(func() {
			assert.Fail(s.T(), "executed complete")
		}),
	)
	observable.input() <- errors.New("failure")
	wg.Wait()
}

func (s *BehaviorObservableSuite) TestBehaviorObservable_Pipe() {
	observable := NewBehaviorObservable(nil, "hello world")
	future := observable.Pipe(
		Filter(func(value interface{}) bool {
			return strings.Compare(value.(string), "unittest") != 0
		}),
		Take(1),
	).ToPromise()
	observable.input() <- "unittest"
	observable.input() <- "unittest1"
	assert.Equal(s.T(), "hello world", <-future())
	observable.input() <- "unittest2"
	assert.Equal(s.T(), nil, <-future())
	observable.input() <- "hello world"
	assert.Equal(s.T(), nil, <-future())
}

func (s *BehaviorObservableSuite) TestBehaviorObservable_Pipe2() {
	observable := NewBehaviorObservable(nil, "unittest")
	future := observable.Pipe(
		Filter(func(value interface{}) bool {
			return strings.Compare(value.(string), "unittest") != 0
		}),
		Take(1),
	).ToPromise()
	observable.input() <- "unittest"
	observable.input() <- "unittest1"
	assert.Equal(s.T(), "unittest1", <-future())
}

func (s *BehaviorObservableSuite) TestBehaviorObservable_PipeMultipleTake() {
	observable := NewBehaviorObservable(nil, "hello world")
	future := observable.Pipe(
		Filter(func(value interface{}) bool {
			return strings.Compare(value.(string), "unittest") != 0
		}),
		Take(2),
	).ToPromise()
	observable.input() <- "unittest"
	observable.input() <- "unittest1"
	assert.Equal(s.T(), "hello world", <-future())
	observable.input() <- "unittest2"
	assert.Equal(s.T(), nil, <-future())
	observable.input() <- "hello world"
	assert.Equal(s.T(), nil, <-future())
}

func (s *SimpleObservableSuite) TestBehaviorObservable_PipeSubscribe() {
	wg := sync.WaitGroup{}

	wg.Add(2)
	observable := NewBehaviorObservable(nil, "hello world")
	sub := observable.Pipe(
		Filter(func(value interface{}) bool {
			return value != "unittest"
		}),
		Take(2),
	).Subscribe(
		NextFunc(func(item interface{}) {
			wg.Done()
		}),
		ErrorFunc(func(err error) {
			assert.Fail(s.T(), "executed error")

		}),
		CompleteFunc(func() {
			wg.Done()
		}),
	)

	observable.input() <- "unittest"
	observable.input() <- "unittest1"
	observable.input() <- "unittest2"
	observable.input() <- "unittest3"
	wg.Wait()

	wg.Add(1)
	sub.Unsubscribe()
	wg.Wait()
}

func TestBehaviorObservableSuite(t *testing.T) {
	suite.Run(t, new(BehaviorObservableSuite))
}
