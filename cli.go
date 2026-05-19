package main

import (
	"fmt"
	"strconv"
	"strings"
)

// @return     exit
func ParseCommand(tokens []string, session *Session) bool {
	if len(tokens) == 0 {
		return false
	}
	switch strings.ToLower(tokens[0]) {
	case "exit", "quit":
		return true

	case "help", "?":
		CliPrintHelp()

	case "overview":
		CliOverview(session)

	case "state":
		CliGetState(session)

	case "select":
		if len(tokens) < 2 {
			fmt.Printf("\tNot enough args. Usage: select <snap>\n\n")
			return false
		}
		CliSelectSnap(session, tokens[1])

	case "graphs":
		if session.Selected == nil {
			fmt.Printf("\tNo snapshot selected.\n\tSelect a snapshot with:\n\t\tselect <name>\n\n")
			return false
		}
		CliGetGraphs(session.Selected)

	case "show", "view":
		if len(tokens) < 3 {
			fmt.Printf("\tNot enough args. Usage: show <type> <id>\n\n")
			return false
		}
		switch tokens[1] {
		case "snap", "snapshot":
			CliViewSnap(tokens[2], session)

		case "graph":
			if session.Selected == nil {
				fmt.Printf("\tNo snapshot selected\n\tSelect a snapshot with:\n\t\tselect <name>\n\n")
				return false
			}
			index, err := strconv.Atoi(tokens[2])
			if err != nil {
				fmt.Printf("\tFailed to convert \"%s\" to number\n\t\tError: %v\n\n", tokens[2], err)
				return false
			}
			CliViewGraph(session.Selected, index)

		case "process", "ps":
			if session.Selected == nil {
				fmt.Printf("\tNo snapshot selected\n\tSelect a snapshot with:\n\t\tselect <name>\n\n")
				return false
			}
			pid, err := strconv.Atoi(tokens[2])
			if err != nil {
				fmt.Printf("\tFailed to convert \"%s\" to number\n\t\tError: %v\n\n", tokens[2], err)
				return false
			}
			CliViewProcess(session.Selected, uint32(pid))
		}

	case "find":
		if len(tokens) < 2 {
			fmt.Printf("\tNot enough args. Usage: find <min>\n\n")
			return false
		}
		if session.Selected == nil {
			fmt.Printf("\tNo snapshot selected\n\tSelect a snapshot with:\n\t\tselect <name>\n\n")
			return false
		}
		min, err := strconv.Atoi(tokens[1])
		if err != nil {
			fmt.Printf("\tFailed to convert \"%s\" to number\n\t\tError: %v\n\n", tokens[1], err)
			return false
		}
		CliGetByConnection(*session.Selected, min)
	}
	return false
}

func CliPrintHelp() {
	fmt.Printf("\tThis is a tool for capturing and analyzing process relationship graphs.\n")
	fmt.Printf("\t-----------------------------------------------------------------------\n", line)
	fmt.Println("\tAvailable commands:")
	fmt.Println("\t\thelp [command]  Show this help message.")
	fmt.Println("\t\texit            Exit the command line.")
	fmt.Println("\t\tstate           Show the current state, in regards to snaps.")
	fmt.Println("\t\tselect <snap>   Select a snapshot for analysis (by name).")
	fmt.Println("\t\toverview        Get a quick overview about session and selected snap.")
	fmt.Println("\t\tgraphs          View the graphs in the currently selected snap.")
	//fmt.Println("\t\tpools           ")
	fmt.Println("\t\tfind <min>      Find all processes with more than min connections.")
	fmt.Println()
}

func CliViewSnap(name string, session *Session) {
	var snap *Snapshot
	for _, s := range session.Snapshots {
		if s.Name == name {
			snap = &s
			break
		}
	}
	if snap == nil {
		fmt.Printf("\tNo snapshot \"%s\" found in current session.\n\n", name)
		return
	}

	fmt.Printf("\t\t: %d separate graphs\n", len(snap.Graphs))
	fmt.Printf("\t\t: %d nodes total\n", snap.GetNodeCount())
	fmt.Printf("\t\t: %d connections total\n", snap.GetTotalConnections())
	fmt.Printf("\t\t: +%dms\n\n", snap.Interval)
}

func CliViewGraph(snap *Snapshot, index int) {
	if index >= len(snap.Graphs) {
		fmt.Printf("\tThe highest graph index in this snapshot is %d\n\n", len(snap.Graphs)-1)
		return
	}

	graph := snap.Graphs[index]
	fmt.Printf("\t\t: %d nodes total\n", len(graph))
	fmt.Printf("\t\t: %d connections total\n", graph.GetTotalConnections())

	avg, median := graph.GetAvgAndMedianConnections()
	fmt.Printf("\t\t: %.1f median connections\n", median)
	fmt.Printf("\t\t: %.1f avg connections\n\n", avg)
}

func CliViewProcess(snap *Snapshot, pid uint32) {
	var process *ProcessSnapshot
	for _, graph := range snap.Graphs {
		if ps, exists := graph[pid]; exists {
			process = ps
		}
	}

	if process == nil {
		fmt.Printf("\tNo process with pid %d found in current snapshot\n\n", pid)
		return
	}

	fmt.Printf("\t[*] Process %d (%d connections)\n",
		process.ProcessId, len(process.Connections))
	fmt.Printf("\t\t* Path: %s\n", process.Name)
	fmt.Printf("\t\t* Parent: %s (PID %d)\n", process.ParentName, process.ParentPid)
	fmt.Printf("\t\t* Signed: ")
	if process.IsSigned {
		fmt.Printf("TRUE\n")
	} else {
		fmt.Printf("FALSE\n")
	}
	fmt.Printf("\t\t* Elevated: ")
	if process.IsElevated {
		fmt.Printf("TRUE\n")
	} else {
		fmt.Printf("FALSE\n")
	}
	fmt.Println()
}

func CliOverview(session *Session) {
	session.PrintDescription()
	if session.Selected == nil {
		fmt.Printf("\tNo snapshot selected. %d available\n", len(session.Snapshots))
		fmt.Printf("\tSelect a snapshot with:\n\t\tselect <name>\n\n")
		return
	}

	fmt.Printf("\t%s\n", line)
	fmt.Printf("\t[ %d graphs, %d total nodes, %d total connections ]\n\n",
		len(session.Selected.Graphs), session.Selected.GetNodeCount(), session.Selected.GetTotalConnections())

	nodes := session.Selected.GetMostWideReaching(3)
	fmt.Printf("\tMost wide-reaching processes:\n")
	for _, node := range nodes {
		fmt.Printf("\t- %s (PID %d)  : %d connections\n",
			node.Name, node.ProcessId, len(node.Connections))
	}
	fmt.Println()
}

func CliSelectSnap(session *Session, name string) {
	for i, snap := range session.Snapshots {
		if snap.Name == name {
			session.Selected = &session.Snapshots[i]
			return
		}
	}
	fmt.Printf("\tCurrently selected session has no snapshot \"%s\"\n\n", name)
}

// Print the current state.
func CliGetState(session *Session) {
	session.PrintDescription()
	fmt.Printf("\t%s\n", line)
	session.PrintSelected()
}

func CliGetGraphs(snap *Snapshot) {
	for i, g := range snap.Graphs {
		fmt.Printf("\t%d) %d nodes, %d connections\n", i, len(g), g.GetTotalConnections())
	}
	fmt.Println()
}

func CliGetPools(snap Snapshot, rule Traversal) {
	pools := snap.CreatePools()
	avg, median := GetAvgAndMedian(pools)
	for i, pool := range pools {
		fmt.Printf("\n%s\n\n\tPool %d:\n", line, i)
		for pid, p := range pool {
			fmt.Printf("\t*\t%s (PID %d)", p.Name, pid)
		}
	}
	fmt.Printf("\n%s\n\n", line)
	fmt.Printf("\t[*] Created %d pools (%.2f avg size, %.2f median)\n\n", len(pools), avg, median)
}

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
