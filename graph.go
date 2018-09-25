package yafsm

import (
	"github.com/awalterschulze/gographviz"
)

func CreateDOT(transitions []Transition) string {
	root := "G"
	graph := gographviz.NewEscape()
	graph.SetDir(true)
	graph.Attrs.Add("rankdir", "LR")
	for _, t := range transitions {
		var attrs map[string]string

		if t.Name() != "" {
			attrs = map[string]string{"label": t.Name()}
		}

		graph.AddNode(root, string(t.To()), nil)
		for _, f := range t.From() {
			graph.AddNode(root, string(f), nil)
			graph.AddEdge(string(f), string(t.To()), true, attrs)
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
		var name string
		states = append(states, State(key))
		from := make([]State, 0, len(val))
		for srcKey, edges := range val {
			from = append(from, State(srcKey))

			if len(edges) != 0 {
				name = edges[0].Attrs["label"]
			}
		}

		options := make([]TransitionConfig, 0, 1)
		if name != "" {
			options = append(options, WithName(name))
		}
		trans = append(trans, NewTransition(NewStates(from...), State(key), options...))
	}

	return states, trans, nil
}
