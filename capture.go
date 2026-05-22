package main

import (
	"bufio"
	"os"
	"time"
)

func BeginCapture(max int) error {
	reader := bufio.NewReader(os.Stdin)
	path := GetInput(reader, "Enter path to save capture in (skip for default)")

	//TODO: check if path is valid
	//TODO: if it doesnt exist, ask to create

	var session Session
	session.Description = GetInput(reader, "Set description for capture")
	session.Timestamp = time.Now()

	g_ProcessTable = CreateProcessTable()
	//TODO: start building graph

	//* capture loop
	snapTicker := time.NewTicker(time.Duration(SNAP_INTERVAL) * time.Minute)
	handleTicker := time.NewTicker(time.Duration(HANDLE_REFRESH_INTERVAL) * time.Second)
	psTicker := time.NewTicker(time.Duration(PS_REFRESH_INTERVAL) * time.Second)
	for {
		if len(session.Snapshots) >= max {
			break
		}
		select {
		case <-snapTicker.C:
			relativeTime := time.Now().UnixMilli() - session.Timestamp.UnixMilli()
			TakeSnapshot(relativeTime)
		case <-psTicker.C:
			ScanProcesses()
		case <-handleTicker.C:
			ParseHandleTable()
		}
	}
	session.WriteToDisk(path)
	return nil
}

func LoadSession(dir string) (Session, error) {
	// - traverse dir
	// -
}

// interval is the relative offset from beginning of capture (seconds)
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
