package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func BeginCapture(max int) error {
	reader := bufio.NewReader(os.Stdin)
	path := GetInput(reader, "Enter path to save capture in (skip for default)")

	if !DirectoryExists(path) && !DirectoryExists(filepath.Join(path, BASE_DIR_NAME)) {
		answer := GetInput(reader, "The directory does not exist, do you want to create it? (y/n) ")
		if answer == strings.ToLower("y") || answer == "" {
			//TODO: create directory
		} else {
			return fmt.Errorf("filepath does not exist")
		}
	}

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

func (s Session) WriteToDisk(path string) error {
	//
	//TODO: write object access
	//TODO: write graphsnapshots
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
