package yafsm_test

import (
	"strings"
	"testing"

	"github.com/bcen/yafsm"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransitionsFromDOTError(t *testing.T) {
	dot := "digraph G {"
	_, _, err := yafsm.CreateTransitionsFromDOT(dot)
	assert.NotNil(t, err)
}

func TestCreateTransitionsFromDOT(t *testing.T) {
	dot := `
	digraph G {
		todo -> todo;
		inprogress -> todo;
		verify -> todo;
		todo -> inprogress;
		inprogress -> inprogress;
		verify -> inprogress;
		inprogress -> verify;
		verify -> verify;
		verify -> done;
	}
	`

	states, trans, err := yafsm.CreateTransitionsFromDOT(dot)
	assert.Nil(t, err)
	assert.True(t, states.Has("todo"))
	assert.True(t, states.Has("inprogress"))
	assert.True(t, states.Has("verify"))
	assert.True(t, states.Has("done"))

	handler := yafsm.CreateTransitionHandler(trans)

	err = handler("todo", "inprogress")
	assert.Nil(t, err)

	err = handler("todo", "done")
	assert.NotNil(t, err)
}

func TestCreateDOTString(t *testing.T) {
	dot := yafsm.CreateDOTString(KanbanTransitions)
	expected := `digraph  {
	todo->todo;
	"in progress"->todo;
	verify->todo;
	todo->"in progress";
	"in progress"->"in progress";
	verify->"in progress";
	"in progress"->verify;
	verify->verify;
	verify->done;
	"in progress";
	done;
	todo;
	verify;

}`
	assert.Equal(t, strings.TrimSpace(dot), strings.TrimSpace(expected))
}
