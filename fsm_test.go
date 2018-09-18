package yafsm_test

import (
	"fmt"
	"testing"

	"github.com/bcen/yafsm"
	"github.com/stretchr/testify/assert"
)

const (
	Todo       yafsm.State = "todo"
	InProgress yafsm.State = "in progress"
	Verify     yafsm.State = "verify"
	Done       yafsm.State = "done"

	Foo yafsm.State = "foo"
)

var (
	AllStates        = yafsm.NewStates(Todo, InProgress, Verify, Done)
	KanbanTransition = yafsm.CreateTransitionHandler(
		[]yafsm.Transition{
			yafsm.NewTransition(yafsm.NewStates(Todo, InProgress, Verify), Todo),
			yafsm.NewTransition(yafsm.NewStates(Todo, InProgress, Verify), InProgress),
			yafsm.NewTransition(yafsm.NewStates(InProgress, Verify), Verify),
			yafsm.NewTransition(yafsm.NewStates(Verify), Done),
		},
	)
)

func TestStates(t *testing.T) {
	testCases := []struct {
		s        yafsm.State
		expected bool
	}{
		{Todo, true},
		{Done, true},
		{Foo, false},
	}

	for _, testCase := range testCases {
		actual := AllStates.Has(testCase.s)
		assert.Equal(t, actual, testCase.expected)
	}
}

func TestTransition(t *testing.T) {
	testCases := []struct {
		from  yafsm.State
		to    yafsm.State
		valid bool
	}{
		{Todo, InProgress, true},
		{InProgress, Todo, true},
		{Verify, Done, true},
		{Verify, Verify, true},
		{Todo, Done, false},
		{Done, Verify, false},
		{Done, Done, false},
		{Todo, Foo, false},
	}

	for _, testCase := range testCases {
		subtest := fmt.Sprintf("(%s,%s)", testCase.from, testCase.to)
		t.Run(subtest, func(t *testing.T) {
			err := KanbanTransition(testCase.from, testCase.to)
			if testCase.valid {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestSingleTransition(t *testing.T) {
	var called bool
	cb := func(tran yafsm.Transition, from, to yafsm.State) error {
		called = true
		return nil
	}
	beGood := yafsm.NewTransition(
		yafsm.NewStates("bad"),
		"good",
		yafsm.WithCallback(cb),
	)

	// valid transition
	assert.Nil(t, beGood.TransitionFrom("bad"))
	assert.True(t, called)

	// invalid transition
	assert.NotNil(t, beGood.TransitionFrom("wut"))
}

func TestSingleTransitionCallbackError(t *testing.T) {
	cbErr := fmt.Errorf("error")
	cb := func(tran yafsm.Transition, from, to yafsm.State) error {
		return cbErr
	}
	beGood := yafsm.NewTransition(yafsm.NewStates("bad"), "good")

	err := beGood.TransitionFrom("bad", yafsm.WithCallback(cb))
	assert.Equal(t, err, cbErr)
}