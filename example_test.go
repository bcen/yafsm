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
