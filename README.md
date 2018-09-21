# Yet Another Finite State Machine

## A toy, DO NOT USE IN PRODUCTION.


### Examples:

```go
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
if err != nil {
    fmt.Println(err)
}

err := handler(green, red)
if err == nil {
    fmt.Println("Yay~")
}
```
