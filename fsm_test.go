package yafsm_test

import (
	"fmt"
	"testing"

	"github.com/bcen/yafsm"
	"github.com/stretchr/testify/assert"
)

const (
	StateIdle        yafsm.State = "idle"
	StateConnect     yafsm.State = "connect"
	StateActive      yafsm.State = "active"
	StateOpenSent    yafsm.State = "opensent"
	StateOpenConfirm yafsm.State = "openconfirm"
	StateEsablished  yafsm.State = "esablished"

	Foo yafsm.State = "foo"
)

var (
	BGPStates = yafsm.NewStates(
		StateIdle,
		StateConnect,
		StateActive,
		StateOpenSent,
		StateOpenConfirm,
		StateEsablished,
	)
	BGPTransitions = []yafsm.Transition{
		yafsm.NewTransition(
			yafsm.NewStates(StateIdle, StateEsablished, StateOpenSent, StateConnect, StateActive, StateOpenConfirm),
			StateIdle,
			yafsm.WithName("Reset"),
		),
		yafsm.NewTransition(
			yafsm.NewStates(StateIdle, StateConnect, StateActive),
			StateConnect,
			yafsm.WithName("Connect"),
		),
		yafsm.NewTransition(
			yafsm.NewStates(StateConnect, StateActive, StateOpenSent),
			StateActive,
			yafsm.WithName("Establish Connection"),
		),
		yafsm.NewTransition(
			yafsm.NewStates(StateConnect, StateActive),
			StateOpenSent,
			yafsm.WithName("Send"),
		),
		yafsm.NewTransition(
			yafsm.NewStates(StateOpenSent, StateOpenConfirm),
			StateOpenConfirm,
			yafsm.WithName("Confirm"),
		),
		yafsm.NewTransition(
			yafsm.NewStates(StateEsablished, StateOpenConfirm),
			StateEsablished,
			yafsm.WithName("Done"),
		),
	}
	BGPTransitionHandler = yafsm.CreateTransitionHandler(BGPTransitions)
)

func TestStates(t *testing.T) {
	testCases := []struct {
		s        yafsm.State
		expected bool
	}{
		{StateIdle, true},
		{StateActive, true},
		{Foo, false},
	}

	for _, testCase := range testCases {
		t.Run(string(testCase.s), func(t *testing.T) {
			actual := BGPStates.Has(testCase.s)
			assert.Equal(t, actual, testCase.expected)
		})
	}
}

func TestTransition(t *testing.T) {
	testCases := []struct {
		from  yafsm.State
		to    yafsm.State
		valid bool
	}{
		{StateIdle, StateConnect, true},
		{StateConnect, StateConnect, true},
		{StateConnect, StateOpenSent, true},
		{StateIdle, StateEsablished, false},
		{StateEsablished, StateConnect, false},
		{StateOpenSent, StateOpenSent, false},
		{StateIdle, Foo, false},
	}

	for _, testCase := range testCases {
		subtest := fmt.Sprintf("(%s,%s)", testCase.from, testCase.to)
		t.Run(subtest, func(t *testing.T) {
			err := BGPTransitionHandler(testCase.from, testCase.to)
			if testCase.valid {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestCreateTransitionHandlerPanicFromDupe(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r)
	}()
	yafsm.CreateTransitionHandler([]yafsm.Transition{
		yafsm.NewTransition(yafsm.NewStates("1", "2"), "3"),
		yafsm.NewTransition(yafsm.NewStates("2"), "3"),
	})
}

func TestCreateTransitionHandlerPanicFromEmptyFromStates(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r)
	}()
	yafsm.CreateTransitionHandler([]yafsm.Transition{
		yafsm.NewTransition(yafsm.NewStates(), "2"),
	})
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

func TestTransitionCallbackOverride(t *testing.T) {
	id := 0
	expected := 2

	cb1 := func(tran yafsm.Transition, from, to yafsm.State) error {
		id = 1
		return nil
	}
	cb2 := func(tran yafsm.Transition, from, to yafsm.State) error {
		id = expected
		return nil
	}

	beGood := yafsm.NewTransition(yafsm.NewStates("bad"), "good", yafsm.WithCallback(cb1))
	beGood.TransitionFrom("bad", yafsm.WithCallback(cb2))
	assert.Equal(t, id, expected)
}
