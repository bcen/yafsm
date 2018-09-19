package yafsm

import (
	"fmt"
	"strings"
)

type State string
type States []State
type Callback func(t Transition, from State, to State) error
type TransitionConfig func(*config)

type config struct {
	callback Callback
}

func getConfig(options ...TransitionConfig) *config {
	c := &config{}
	for _, opt := range options {
		opt(c)
	}
	return c
}

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
	TransitionFrom(State, ...TransitionConfig) error
	GetCallback() Callback
}

type transition struct {
	from States
	to   State

	callback Callback
}

func (t transition) From() States {
	return t.from
}

func (t transition) To() State {
	return t.to
}

func (t transition) GetCallback() Callback {
	return t.callback
}

func (t transition) TransitionFrom(from State, options ...TransitionConfig) error {
	return doTransition([]Transition{t}, from, t.To(), options...)
}

func WithCallback(cb Callback) TransitionConfig {
	return func(c *config) {
		c.callback = cb
	}
}

func NewTransition(from States, to State, options ...TransitionConfig) Transition {
	c := getConfig(options...)
	return &transition{from, to, c.callback}
}

func CreateTransitionHandler(trans []Transition) func(State, State, ...TransitionConfig) error {
	return func(from, to State, options ...TransitionConfig) error {
		return doTransition(trans, from, to, options...)
	}
}

func getCallback(t Transition, c *config) Callback {
	cb := c.callback
	if cb == nil {
		cb = t.GetCallback()
	}
	if cb == nil {
		// noop
		cb = func(t Transition, from, to State) error { return nil }
	}
	return cb
}

func doTransition(trans []Transition, from, to State, options ...TransitionConfig) error {
	var tran Transition
	c := getConfig(options...)

Loop:
	for _, t := range trans {
		// we do not check for duplicates

		if t.To() != to {
			continue
		}

		for _, f := range t.From() {
			if f == from {
				tran = t
				break Loop
			}
		}
	}

	if tran == nil {
		return fmt.Errorf("\"%s\" -> \"%s\" is not a valid transition", from, to)
	}

	cb := getCallback(tran, c)
	return cb(tran, from, to)
}
