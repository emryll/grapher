package main

import (
	"fmt"
	"strings"
)

// @return     exit
func ParseCommand(tokens []string, session *Session) bool {
	if len(tokens) == 0 {
		return false
	}
	switch strings.ToLower(tokens[0]) {
	case "help":
		CliPrintHelp()
	case "overview":
		CliOverview(session)
	case "state":
		CliGetState(session)
	case "select":
		if len(tokens) < 2 {
			fmt.Printf("\tNot enough args. Usage: select <id>")
			return false
		}
		CliSelectSnap(session, tokens[1])
	case "show":
	case "graphs":
		if session.Selected == nil {
			fmt.Printf("\tNo snapshot selected.\n\tSelect a snapshot with:\n\t\tselect <name>\n")
			return false
		}
		CliGetGraphs(session.Selected)
	}
	return false
}

func CliPrintHelp() {

}

// TODO: overview command
func CliOverview(session *Session) {
	session.PrintDescription()
	if session.Selected == nil {
		fmt.Printf("\tNo snapshot selected. %d available\n", len(session.Snapshots))
		fmt.Printf("\tSelect a snapshot with:\n\t\tselect <name>\n")
		return
	}

	fmt.Printf("\t%s\n", line)
	nodes := GetMostWideReaching(3)
	fmt.Printf("\tMost wide-reaching processes:\n")
	for _, node := range nodes {
		fmt.Printf("\t- %s (PID %d)  : %d connections\n",
			node.Process.Name, node.ProcessId, len(node.Process.Connections))
	}
}

func CliSelectSnap(session *Session, name string) {
	for i, snap := range session.Snapshots {
		if snap.Name == name {
			session.Selected = &session.Snapshots[i]
			return
		}
	}
	fmt.Printf("\tCurrently selected session has no snapshot %s\n", name)
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
}

func CliGetPools(snap Snapshot, rule Traversal) {
	pools := snap.CreatePools()
	avg, median := GetAvgAndMedian(pools)
	fmt.Printf("\t[*] Created %d pools (%.2f avg size, %.2f median)\n", len(pools), avg, median)
	for _, pool := range pools {
		fmt.Printf("\n%s\n\n\tPool %d:\n")
		for pid, p := range pool {
			fmt.Printf("\t*\t%s (PID %d)", p.Name, pid)
		}
	}
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
