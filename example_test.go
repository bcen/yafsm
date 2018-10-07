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

	cb := func(t yafsm.Transition, from, to yafsm.State) error {
		fmt.Printf("From callback: %s -> %s\n", from, to)
		return nil
	}

	handler := yafsm.CreateTransitionHandler([]yafsm.Transition{
		yafsm.NewTransition(yafsm.NewStates(red), green),
		yafsm.NewTransition(yafsm.NewStates(green), yellow, yafsm.WithCallback(cb)),
		yafsm.NewTransition(yafsm.NewStates(yellow), red),
	})

	err := handler(green, red)
	fmt.Println(err)

	err = handler(green, yellow)
	fmt.Println(err)
	// Output:
	// "Green" -> "Red" is not a valid transition
	// From callback: Green -> Yellow
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
		yafsm.NewTransition(yafsm.NewStates(verify), done, yafsm.WithName("Mark Done")),
	}

	dot := yafsm.CreateDOT(transitions)
	fmt.Println(dot)
	// Output:
	// digraph  {
	// 	rankdir=LR;
	// 	todo->todo;
	// 	inprogress->todo;
	// 	verify->todo;
	// 	todo->inprogress;
	// 	inprogress->inprogress;
	// 	verify->inprogress;
	// 	inprogress->verify;
	// 	verify->verify;
	// 	verify->done[ label="Mark Done" ];
	// 	done;
	// 	inprogress;
	// 	todo;
	// 	verify;
	//
	//}
}

func ExampleCreateTransitionsFromDOT() {
	dot := `
	digraph {
	    green -> yellow [label = "an optional label"];
	    yellow -> red;
	    red -> green;
	}
	`
	states, trans, _ := yafsm.CreateTransitionsFromDOT(dot)
	handler := yafsm.CreateTransitionHandler(trans)

	fmt.Println("States:")
	for _, s := range states.AsSortedStrings() {
		fmt.Println(s)
	}

	fmt.Println("Transitions:")

	if err := handler("green", "red"); err != nil {
		fmt.Println(err)
	}

	if err := handler("green", "yellow"); err == nil {
		fmt.Println(`"Green" -> "Yellow": OK`)
	}
	// Output:
	// States:
	// green
	// red
	// yellow
	// Transitions:
	// "Green" -> "Red" is not a valid transition
	// "Green" -> "Yellow": OK
}
