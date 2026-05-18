package main

import (
	"sync"
	"time"
)

var DEFAULT_BANNER = 0

// describes a capture
type Session struct {
	Timestamp   time.Time
	Description string
	Snapshots   []Snapshot
}

type Bitmask uint32
type GraphSnapshot map[uint32]*ProcessSnapshot

type Snapshot struct {
	Graphs []GraphSnapshot
	//TODO: object access registry
	Interval uint32 // offset from timestamp in seconds
}

// Describes a process (node)
type ProcessSnapshot struct {
	Connections map[uint32]*Connection
	Name        string
	ProcessId   uint32
	ParentName  string
	ParentId    uint32
	IsSigned    bool
	IsElevated  bool
}

//*=========================[ Graphing ]===========================

type Pool map[uint32]*ProcessSnapshot

type Graph struct {
	mu      sync.RWMutex
	Members map[uint32]*ProcessNode
}

type GraphRegistry []*Graph

type ProcessNode struct {
	ProcessId   uint32
	Process     *Process
	Connections map[uint32]*Connection
}

type Connection struct {
	Type   Bitmask
	Weight int
}

type Traversal struct {
	flags  Bitmask
	weight int
}
