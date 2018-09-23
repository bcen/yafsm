package yafsm_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/bcen/yafsm"
	"github.com/stretchr/testify/assert"
)

func getDOT(path string) string {
	content, err := ioutil.ReadFile("testdata/" + path)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func TestCreateTransitionsFromDOTError(t *testing.T) {
	dot := "digraph G {"
	_, _, err := yafsm.CreateTransitionsFromDOT(dot)
	assert.NotNil(t, err)
}

func TestCreateTransitionsFromDOT(t *testing.T) {
	dot := getDOT("bgp.dot")

	states, trans, err := yafsm.CreateTransitionsFromDOT(dot)
	assert.Nil(t, err)
	for _, s := range BGPStates {
		assert.True(t, states.Has(s))
	}

	handler := yafsm.CreateTransitionHandler(trans)

	err = handler(StateConnect, StateOpenSent)
	assert.Nil(t, err)

	err = handler(StateIdle, StateEsablished)
	assert.NotNil(t, err)
}

func TestCreateDOT(t *testing.T) {
	dot := yafsm.CreateDOT(BGPTransitions)
	expected := getDOT("bgp.dot")
	assert.Equal(t, strings.TrimSpace(dot), strings.TrimSpace(expected))
}
