package main

import (
	"fmt"

	"github.com/bcen/yafsm"
)

const (
	Green  yafsm.State = "green"
	Yellow yafsm.State = "yellow"
	Red    yafsm.State = "red"
)

func printTransition(t yafsm.Transition, from, to yafsm.State) error {
	fmt.Printf("Changing '%s' to '%s'\n", from, to)
	return nil
}

var (
	GoGreen                = yafsm.NewTransition(yafsm.NewStates(Red), Green, yafsm.WithCallback(printTransition))
	GoYellow               = yafsm.NewTransition(yafsm.NewStates(Green), Yellow, yafsm.WithCallback(printTransition))
	GoRed                  = yafsm.NewTransition(yafsm.NewStates(Yellow), Red, yafsm.WithCallback(printTransition))
	TrafficLightTransition = yafsm.CreateTransitionHandler([]yafsm.Transition{
		GoGreen,
		GoYellow,
		GoRed,
	})
)

func main() {
	fmt.Println("Testing single transition")
	fmt.Println("----")
	GoGreen.TransitionFrom(Red)
	GoYellow.TransitionFrom(Green)
	GoRed.TransitionFrom(Yellow)

	fmt.Println("Testing transition error")
	fmt.Println("----")
	// invalid transition
	if err := GoGreen.TransitionFrom(Yellow); err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	fmt.Println("Testing transition handler")
	fmt.Println("----")
	TrafficLightTransition(Red, Green)
	if err := TrafficLightTransition(Red, Yellow); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
