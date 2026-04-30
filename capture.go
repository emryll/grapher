package main

import (
	"bufio"
	"fmt"
	"os"
)

func BeginCapture() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter path to save capture in (ENTER for default): ")
	path, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	//TODO: check if path is valid
	//TODO: if it doesnt exist, ask to create

	//TODO: ask description
	fmt.Print("Set description for capture: ")
	description, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	var session Session
	session.Description = description

}

func LoadSession(dir string) (Session, error) {

}

func (gr GraphRegistry) TakeSnapshot(interval int) Snapshot {
	var snap Snapshot
	for _, graph := range gr {
		var graphSnap GraphSnapshot
		for pid, p := range graph.Members {
			psSnap := &ProcessSnapshot{
				ProcessId:   p.ProcessId,
				Connections: make(map[uint32]*Connection),
			}
			if p.Process != nil {
				psSnap.Name = p.Process.Name
				psSnap.ParentPid = p.Process.ParentPid
				psSnap.IsElevated = p.Process.IsElevated
				psSnap.IsSigned = p.Process.IsSigned
			}
			graphSnap[pid] = psSnap
		}
		snap.Graphs = append(snap.Graphs, graphSnap)
	}

	snap.Interval = uint32(interval)
	return snap
}
