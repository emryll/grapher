package main

import "fmt"

func ParseCommand(tokens []string) error {
	//TODO:

	return nil
}

// TODO: overview command
func CliOverview() {
	//TODO:
}

func CliGetState(session Session) {
	//TODO: print run description
	//TODO: print which snapshot is selected

	//TODO: if none are selected, list available snapshots
	//TODO: if one is selected, print graph count, node count, connection count + stamp
}

func (gs GraphSnapshot) GetTotalConnections() int {
	var connections int
	for _, p := range gs {
		connections += len(p.Connections)
	}
	return connections
}

func CliGetGraphs(snap Snapshot) {
	for i, g := range snap.Graphs {
		fmt.Printf("\t%d) %d nodes, %d connections\n", i, len(g), g.GetTotalConnections())
	}
}

//TODO: view pool command
//TODO: get nodes with more than n connections
