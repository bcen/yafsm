package yafsm

import "github.com/awalterschulze/gographviz"

func CreateDOTString(transitions []Transition) string {
	root := "G"
	graph := gographviz.NewEscape()
	graph.SetDir(true)
	for _, t := range transitions {
		graph.AddNode(root, string(t.To()), nil)
		for _, f := range t.From() {
			graph.AddNode(root, string(f), nil)
			graph.AddEdge(string(f), string(t.To()), true, nil)
		}
	}
	return graph.String()
}

func CreateTransitionsFromDOT(dot string) (States, []Transition, error) {
	graph, err := gographviz.Read([]byte(dot))
	if err != nil {
		return nil, nil, err
	}

	trans := make([]Transition, 0, len(graph.Edges.DstToSrcs))
	states := make(States, 0, len(graph.Edges.DstToSrcs))

	for key, val := range graph.Edges.DstToSrcs {
		states = append(states, State(key))
		from := make([]State, 0, len(val))
		for srcKey, _ := range val {
			from = append(from, State(srcKey))
		}

		trans = append(trans, NewTransition(NewStates(from...), State(key)))
	}

	return states, trans, nil
}
