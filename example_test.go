package yafsm_test

import (
	"fmt"

	"github.com/bcen/yafsm"
)

func ExampleCreateTransitionHandler() {
	const (
		green  yafsm.State = "green"
		yellow yafsm.State = "yellow"
		red    yafsm.State = "red"
	)
	handler := yafsm.CreateTransitionHandler([]yafsm.Transition{
		yafsm.NewTransition(yafsm.NewStates(red), green),
		yafsm.NewTransition(yafsm.NewStates(green), yellow),
		yafsm.NewTransition(yafsm.NewStates(yellow), red),
	})

	err := handler(green, red)
	fmt.Println(err)

	err = handler(green, yellow)
	fmt.Println(err)
	// Output:
	// "Green" -> "Red" is not a valid transition
	// <nil>
}

func ExampleCreateDOT() {
	const (
		todo       yafsm.State = "todo"
		inprogress yafsm.State = "inprogress"
		verify     yafsm.State = "verify"
		done       yafsm.State = "done"
	)
	transitions := []yafsm.Transition{
		yafsm.NewTransition(yafsm.NewStates(todo, inprogress, verify), todo),
		yafsm.NewTransition(yafsm.NewStates(todo, inprogress, verify), inprogress),
		yafsm.NewTransition(yafsm.NewStates(inprogress, verify), verify),
		yafsm.NewTransition(yafsm.NewStates(verify), done),
	}

	dot := yafsm.CreateDOT(transitions)
	fmt.Println(dot)
	// Output:
	// digraph  {
	//        todo->todo;
	//        inprogress->todo;
	//        verify->todo;
	//        todo->inprogress;
	//        inprogress->inprogress;
	//        verify->inprogress;
	//        inprogress->verify;
	//        verify->verify;
	//        verify->done;
	//        done;
	//        inprogress;
	//        todo;
	//        verify;

	// }
}
