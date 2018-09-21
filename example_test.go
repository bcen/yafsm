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
	// Output: "Green" -> "Red" is not a valid transition
	// <nil>
}

func ExampleCreateDOTString() {
	const (
		todo       yafsm.State = "todo"
		inprogress yafsm.State = "inprogress"
		verify     yafsm.State = "verify"
		done       yafsm.State = "done"
	)
	transitions := []yafsm.Transition{
		yafsm.NewTransition(yafsm.NewStates(Todo, InProgress, Verify), Todo),
		yafsm.NewTransition(yafsm.NewStates(Todo, InProgress, Verify), InProgress),
		yafsm.NewTransition(yafsm.NewStates(InProgress, Verify), Verify),
		yafsm.NewTransition(yafsm.NewStates(Verify), Done),
	}

	dot := yafsm.CreateDOTString(transitions)
	fmt.Println(dot)
	// Output: digraph  {
	//        todo->todo;
	//        "in progress"->todo;
	//        verify->todo;
	//        todo->"in progress";
	//        "in progress"->"in progress";
	//        verify->"in progress";
	//        "in progress"->verify;
	//        verify->verify;
	//        verify->done;
	//        "in progress";
	//        done;
	//        todo;
	//        verify;

	//}
}
