package main

import "fmt"

// @return     exit, err
func ParseCommand(tokens []string) (bool, error) {
	//TODO:

	return false, nil
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

func PrintBanner(num ...int) {
	var choice int
	if len(num) == 0 {
		choice = DEFAULT_BANNER
	} else {
		choice = num[0]
	}
	switch choice {
	case 0:
		fmt.Println("		                                 d8b 					   ")
		fmt.Println("                                        ?88 					   ")
		fmt.Println("                                         88b 					   ")
		fmt.Println(" d888b8b    88bd88b d888b8b  ?88,.d88b,  888888b  d8888b  88bd88b ")
		fmt.Println("d8P' ?88    88P'  `d8P' ?88  `?88'  ?88  88P `?8bd8b_,dP  88P'  ` ")
		fmt.Println("88b  ,88b  d88     88b  ,88b   88b  d8P d88   88P88b     d88      ")
		fmt.Println("`?88P'`88bd88'     `?88P'`88b  888888P'd88'   88b`?888P'd88'      ")
		fmt.Println("       )88                     88P'                               ")
		fmt.Println("      ,88P                    d88                                 ")
		fmt.Println("  `?8888P                     ?8P                                 ")
	}
	fmt.Printf("\t(%s)\n\n", VERSION)
}
