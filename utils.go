package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

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

func (g GraphSnapshot) GetAvgAndMedianConnections() (float32, float32) {
	if len(g) == 0 {
		return 0, 0
	}

	var (
		total  float32
		median float32
		avg    float32
		nodes  []*ProcessSnapshot
	)

	for _, node := range g {
		nodes = append(nodes, node)
		total += float32(len(node.Connections))
	}
	sort.Slice(nodes, func(i, j int) bool {
		return len(nodes[i].Connections) < len(nodes[j].Connections)
	})

	if len(nodes)%2 == 0 {
		upperMidIndex := len(nodes) / 2
		totalMiddle := len(nodes[upperMidIndex-1].Connections)
		totalMiddle += len(nodes[upperMidIndex].Connections)
		median = float32(totalMiddle) / 2
	} else {
		midIndex := len(nodes) / 2
		median = float32(len(nodes[midIndex].Connections))
	}
	avg = total / float32(len(nodes))

	return avg, median
}

func (s *Session) PrintDescription() {
	fmt.Printf("\tRun description: ")
	if s.Description == "" {
		fmt.Printf("N/A\n")
	} else {
		fmt.Printf("%s\n", s.Description)
	}
}

func (s *Session) PrintSelected() {
	if s.Selected == nil {
		fmt.Printf("\t[*] No snapshot selected (%d available)\n", len(s.Snapshots))
		for _, snap := range s.Snapshots {
			fmt.Printf("\t\t- Snapshot %d (%d nodes, +%dms)\n", snap.SnapId, snap.GetNodeCount(), snap.Interval)
		}
		fmt.Println()
		return
	}
	fmt.Printf("\t[*] Snapshot %s selected\n", s.Selected.Name)
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

// Return the n most wide reaching processes. If n is 0, then all, ranked.
func (s *Snapshot) GetMostWideReaching(n int) []*ProcessSnapshot {
	var nodes []*ProcessSnapshot
	for _, graph := range s.Graphs {
		for _, node := range graph {
			nodes = append(nodes, node)
		}
	}
	// sort them based on connections (descending)
	sort.Slice(nodes, func(i, j int) bool {
		return len(nodes[i].Connections) > len(nodes[j].Connections)
	})
	if n == 0 {
		return nodes
	}
	return nodes[:n]
}

func (b Bitmask) HasFlags(flags Bitmask) bool {
	return (b & flags) == flags
}

func GetProcessExecutable(pid uint32) (string, error) {
	hProcess, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return "", err
	}
	defer windows.CloseHandle(hProcess)

	var buf [windows.MAX_PATH]uint16
	size := uint32(len(buf))
	// flag 0 for Win32 path format
	err = windows.QueryFullProcessImageName(hProcess, 0, &buf[0], &size)
	if err != nil {
		return "", err
	}
	return windows.UTF16ToString(buf[:size]), nil
}

func GetParentPid(handle *windows.Handle) (uint32, error) {
	var (
		pbi    windows.PROCESS_BASIC_INFORMATION
		retLen uint32
	)
	err := windows.NtQueryInformationProcess(
		handle,
		windows.ProcessBasicInformation,
		unsafe.Pointer(&pbi),
		uint32(unsafe.Sizeof(pbi)),
		&retLen,
	)
	if err != nil {
		return 0, err
	}
	return pbi.InheritedFromUniqueProcessId, nil
}
