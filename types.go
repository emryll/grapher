package main

import "time"

// describes a capture
type Session struct {
	Timestamp   time.Time
	Description string
	Snapshots   []Snapshot
}

//TODO: what kind of structure to save it in?
//TODO: should i have nodes and edges separate, or like it is in

type Snapshot struct {
	Processes map[uint32]*ProcessSnapshot
}

// Describes a process (node)
type ProcessSnapshot struct {
	Name       string
	ProcessId  uint32
	ParentName string
	ParentId   uint32
	IsSigned   bool
	IsElevated bool
}
