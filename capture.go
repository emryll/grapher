package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func BeginCapture(max int) error {
	reader := bufio.NewReader(os.Stdin)
	path := GetInput(reader, "Enter path to save capture in (skip for default)")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0644)
		if err != nil {
			return err
		}
	}

	var session Session
	session.Description = GetInput(reader, "Set description for capture")
	session.Timestamp = time.Now()

	g_ProcessTable = CreateProcessTable()
	err := session.InitializeCapture(path)
	if err != nil {
		return fmt.Errorf("failed to initialize capture: %v\n", err)
	}
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
			snapshot := TakeSnapshot(relativeTime)
			err := snapshot.WriteToDisk(path)
			if err != nil {
				fmt.Printf("[ERROR] Failed to write snapshot to disk: %v\n", err)
			}
		case <-psTicker.C:
			ScanProcesses()
		case <-handleTicker.C:
			ParseHandleTable()
		}
	}
	session.WriteToDisk(path)
	return nil
}

func (s Snapshot) WriteToDisk(path string) error {
	graphs, err := json.MarshalIndent(s.Graphs, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(path, s.Name+"_graphs.json"), graphs, 0644)
	if err != nil {
		return err
	}

	oac, err := json.MarshalIndent(g_ObjectAccessRegistry, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(path, s.Name+"_objects.json"), oac, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s Session) InitializeCapture(path string) error {
	metadata, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(path, s.Name+".json"), metadata, 0644)
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
