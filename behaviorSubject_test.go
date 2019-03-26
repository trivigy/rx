package rx

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BehaviorSubjectSuite struct {
	suite.Suite
}

func (s *BehaviorSubjectSuite) TestNewBehaviorSubject() {
	wg := sync.WaitGroup{}

	wg.Add(8)
	onEvent := NewBehaviorSubject("hello world")
	observable1 := onEvent.Observable()
	sub1 := observable1.Subscribe(
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
	sub2 := observable1.Subscribe(
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

	observable2 := onEvent.Observable()
	sub3 := observable2.Subscribe(
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
	sub4 := observable2.Subscribe(
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

	onEvent.Next("hello world")
	wg.Wait()

	wg.Add(4)
	sub1.Unsubscribe()
	sub2.Unsubscribe()
	sub3.Unsubscribe()
	sub4.Unsubscribe()
	wg.Wait()
}

func TestBehaviorSubjectSuite(t *testing.T) {
	suite.Run(t, new(BehaviorSubjectSuite))
}
