package main

import (
	"sync"
	"time"
	"unsafe"
)

var DEFAULT_BANNER = 0

type Bitmask uint32

//*===============[ Handle enumeration ]=================

type cHandleEntry struct {
	Params     unsafe.Pointer
	ParamsSize uint64
	Handle     uint32
	Access     uint32
	ObjType    uint32
	Pid        uint32
}

type HandleEntry struct {
	Params  map[string]Parameter
	ObjType uint32
	Handle  uint32
	Access  uint32
	Pid     uint32
}

type Parameter struct {
	Name      string
	Type      uint8
	Domain    uint8
	TimeStamp int64
}

//*===================[ Snapshots ]========================

// describes a capture
type Session struct {
	Timestamp   time.Time
	Description string
	Snapshots   []Snapshot
}

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
	ParentPid   uint32
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
