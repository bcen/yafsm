package yafsm

import (
	"fmt"
	"strings"
)

type State string
type States []State
type Callback func(t Transition, from State, to State) error
type TransitionOption func(Transition)

func (s State) String() string {
	return strings.Title(string(s))
}

func NewStates(states ...State) States {
	ret := make(States, len(states))
	for i, s := range states {
		ret[i] = s
	}
	return ret
}

func (states States) Has(s State) bool {
	for _, state := range states {
		if state == s {
			return true
		}
	}
	return false
}

type Transition interface {
	From() States
	To() State
	TransitionFrom(State, ...TransitionOption) error
	GetCallback() Callback
	SetCallback(cb Callback)
}

type transition struct {
	from States
	to   State

	callback Callback
}

func (t *transition) From() States {
	return t.from
}

func (t *transition) To() State {
	return t.to
}

func (t *transition) GetCallback() Callback {
	return t.callback
}

func (t *transition) SetCallback(cb Callback) {
	t.callback = cb
}

func (t *transition) TransitionFrom(from State, options ...TransitionOption) error {
	return doTransition([]Transition{t}, from, t.To(), options...)
}

func WithCallback(cb Callback) TransitionOption {
	return func(t Transition) {
		t.SetCallback(cb)
	}
}

func NewTransition(from States, to State, options ...TransitionOption) Transition {
	t := &transition{from, to, nil}
	for _, opt := range options {
		opt(t)
	}
	return t
}

func CreateTransitionHandler(trans []Transition) func(State, State, ...TransitionOption) error {
	return func(from, to State, options ...TransitionOption) error {
		return doTransition(trans, from, to, options...)
	}
}

func doTransition(trans []Transition, from, to State, options ...TransitionOption) error {
	var tran Transition

Loop:
	for _, t := range trans {
		// we do not check for duplicates

		if t.To() != to {
			continue
		}

		for _, f := range t.From() {
			if f == from {
				tran = NewTransition(t.From(), t.To(), WithCallback(t.GetCallback()))
				break Loop
			}
		}
	}

	if tran == nil {
		return fmt.Errorf("\"%s\" -> \"%s\" is not a valid transition", from, to)
	}

	for _, opt := range options {
		opt(tran)
	}

	if cb := tran.GetCallback(); cb != nil {
		if err := cb(tran, from, to); err != nil {
			return err
		}
	}

	return nil
}
