package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// sum / n
func GetAvgAndMedian(pools []Pool) (float64, float64) {
	var (
		total  int
		median float64
		avg    float64
	)
	//TODO: sort pool by length
	for i, pool := range pools {

		//TODO: check if i is in middle
	}
	avg = float64(total / len(pools))

	return avg, median
}

func (s Session) PrintDescription() {
	fmt.Printf("\tRun description: ")
	if s.Description == "" {
		fmt.Printf("N/A\n")
	} else {
		fmt.Printf("%s\n", s.Description)
	}
}

func (s Session) PrintSelected() {
	if s.Selected == nil {
		fmt.Printf("\t[*] No snapshot selected (%d available)\n", len(s.Snapshots))
		for _, snap := range s.Snapshots {
			fmt.Printf("\t\t- Snapshot %d (%d nodes, +%dms)\n", snap.SnapId, snap.GetNodeCount(), snap.Interval)
		}
		fmt.Println()
		return
	}
	fmt.Printf("\t[*] Snapshot %d selected\n", s.Selected.SnapId)
	fmt.Printf("\t\t: %d separate graphs\n", len(s.Selected.Graphs))
	fmt.Printf("\t\t: %d nodes total\n", s.Selected.GetNodeCount())
	fmt.Printf("\t\t: %d connections total\n", s.Selected.GetTotalConnections())
	fmt.Printf("\t\t: +%dms\n\n", s.Selected.Interval)
}

func (s Snapshot) GetNodeCount() int {
	var count int
	for _, graph := range s.Graphs {
		count += len(graph)
	}
	return count
}

func (gs GraphSnapshot) GetTotalConnections() int {
	var connections int
	for _, p := range gs {
		connections += len(p.Connections)
	}
	return connections
}

func (s Snapshot) GetTotalConnections() int {
	var count int
	for _, graph := range s.Graphs {
		count += graph.GetTotalConnections()
	}
	return count
}

func GetInput(reader *bufio.Reader, msg ...string) string {
	if len(msg) > 0 {
		fmt.Printf("%s: ", msg)
	}
	if reader == nil {
		reader = bufio.NewReader(os.Stdin)
	}

	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
