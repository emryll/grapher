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

func CliGetByConnection(snap Snapshot, min int) {
	nodes := GetNodesByNumConnections(snap, min)
	if len(nodes) == 0 {
		fmt.Printf("\t[!] Found no processes with %d or more connections.\n", min)
		return
	}
	fmt.Printf("\t[*] Found %d processes with %d or more connections:\n", min)
	for _, p := range nodes {
		fmt.Printf("\t*\t%s (PID %d)  :  %d connections\n", p.Name, p.ProcessId, len(p.Connections))
	}
}

func GetNodesByNumConnections(snap Snapshot, min int) []*ProcessSnapshot {
	var nodes []*ProcessSnapshot
	for _, g := range snap.Graphs {
		for _, p := range g {
			if len(p.Connections) >= min {
				nodes = append(nodes, p)
			}
		}
	}
	return nodes
}
