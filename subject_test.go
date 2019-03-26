package rx

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SubjectSuite struct {
	suite.Suite
}

func (s *SubjectSuite) TestNewSubject() {
	wg := sync.WaitGroup{}

	wg.Add(4)
	onEvent := NewSubject()
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

func TestSubjectSuite(t *testing.T) {
	suite.Run(t, new(SubjectSuite))
}
