// Package yafsm is a simple library for building finite state machine.
package yafsm

import (
	"fmt"
	"sort"
	"strings"
)

type State string
type States []State
type Callback func(t Transition, from State, to State) error
type TransitionConfig func(*config)

type config struct {
	name     string
	callback Callback
}

func getConfig(options ...TransitionConfig) *config {
	c := &config{}
	for _, opt := range options {
		opt(c)
	}
	return c
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

func createEdges(trans []Transition) map[string]Transition {
	edges := make(map[string]Transition)
	for _, t := range trans {
		if len(t.From()) == 0 {
			panic("Transition must have at least one from states")
		}

		for _, f := range t.From() {
			key := fmt.Sprintf("(%s,%s)", string(f), string(t.To()))
			if _, ok := edges[key]; ok {
				panic(fmt.Sprintf("%s is a duplicate transition", key))
			}

			edges[key] = t
		}
	}
	return edges
}

func (s State) String() string {
	return strings.Title(string(s))
}

// NewStates creates a list of State from input.
func NewStates(states ...State) States {
	ret := make(States, len(states))
	for i, s := range states {
		ret[i] = s
	}
	return ret
}

// Has checks existence of a given State.
func (states States) Has(s State) bool {
	for _, state := range states {
		if state == s {
			return true
		}
	}
	return false
}

func (states States) AsSortedStrings() []string {
	ret := make([]string, len(states))
	for i, s := range states {
		ret[i] = string(s)
	}
	sort.Strings(ret)
	return ret
}

type Transition interface {
	Name() string
	From() States
	To() State
	TransitionFrom(State, ...TransitionConfig) error
	GetCallback() Callback
}

type transition struct {
	name string
	from States
	to   State

	callback Callback
}

func (t transition) Name() string {
	return t.name
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
	return doTransition(createEdges([]Transition{t}), from, t.To(), options...)
}

// WithCallback sets a callback for a given transition.
func WithCallback(cb Callback) TransitionConfig {
	return func(c *config) {
		c.callback = cb
	}
}

// WithName sets an optional name for the transition.
func WithName(name string) TransitionConfig {
	return func(c *config) {
		c.name = name
	}
}

// NewTransition creates a new transition.
func NewTransition(from States, to State, options ...TransitionConfig) Transition {
	c := getConfig(options...)
	return &transition{c.name, from, to, c.callback}
}

// CreateTransitionHandler binds and returns an action handler for the given transitions.
func CreateTransitionHandler(trans []Transition) func(State, State, ...TransitionConfig) error {
	edges := createEdges(trans)
	return func(from, to State, options ...TransitionConfig) error {
		return doTransition(edges, from, to, options...)
	}
}

func doTransition(edges map[string]Transition, from, to State, options ...TransitionConfig) error {
	c := getConfig(options...)

	key := fmt.Sprintf("(%s,%s)", string(from), string(to))
	tran := edges[key]
	if tran == nil {
		return fmt.Errorf(`"%s" -> "%s" is not a valid transition`, from, to)
	}

	cb := getCallback(tran, c)
	return cb(tran, from, to)
}
