package main

import (
	"sync"
	"time"
	"unsafe"
)

var DEFAULT_BANNER = 0

type Bitmask uint32

const (
	SNAP_INTERVAL           = 5  // minutes
	HANDLE_REFRESH_INTERVAL = 10 // seconds
	PS_REFRESH_INTERVAL     = 10 // seconds
)

const (
	OBJECT_TYPE_UNKNOWN   = 0
	OBJECT_TYPE_PROCESS   = 1
	OBJECT_TYPE_THREAD    = 2
	OBJECT_TYPE_FILE      = 3
	OBJECT_TYPE_SEMAPHORE = 4
	OBJECT_TYPE_EVENT     = 5
	OBJECT_TYPE_MUTEX     = 6
	OBJECT_TYPE_SYMLINK   = 7
	OBJECT_TYPE_PIPE      = 8
)

const (
	PARAMETER_ANSISTRING    = 1
	PARAMETER_ASTR_ARRAY    = 10
	PARAMETER_UINT32        = 2
	PARAMETER_UINT32_ARRAY  = 20
	PARAMETER_UINT64        = 3
	PARAMETER_UINT64_ARRAY  = 30
	PARAMETER_BOOLEAN       = 4
	PARAMETER_BOOLEAN_ARRAY = 40
	PARAMETER_POINTER       = 5
	PARAMETER_POINTER_ARRAY = 50
	PARAMETER_BYTES         = 7
)

const (
	ANY_ACCESS       = 1 << 0
	FILE_READ        = 1 << 1
	PS_READ_MEM      = 1 << 1
	FILE_WRITE       = 1 << 2
	PS_CREATE_THREAD = 1 << 2
)

// win32 struct for GetTokenInformation
type TOKEN_ELEVATION struct {
	TokenIsElevated uint32
}

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
	Access  Bitmask
	Pid     uint32
}

type Parameter struct {
	Name   string
	Type   uint8
	Domain uint8
	Buffer []byte
}

//*===================[ Snapshots ]========================

// describes a capture
type Session struct {
	Timestamp   time.Time  `json:"timestamp"`
	Description string     `json:"description"`
	Snapshots   []Snapshot `json:"-"`
	Selected    *Snapshot  `json:"-"` // currently selected snap
}

type GraphSnapshot map[uint32]*ProcessSnapshot

type Snapshot struct {
	Name   string          `json:"name"`
	Graphs []GraphSnapshot `json:"graphs"`
	//TODO: object access registry
	Interval uint32 `json:"interval"` // offset from timestamp in seconds
}

// Describes a process (node)
type ProcessSnapshot struct {
	Connections map[uint32]*Connection
	Name        string `json:"name"`
	ProcessId   uint32 `json:"pid"`
	ParentName  string `json:"parent"`
	ParentPid   uint32 `json:"ppid"`
	IsSigned    bool   `json:"is_signed"`
	IsElevated  bool   `json:"is_elevated"`
	Graph       *Graph `json:"-"`
}

type CaptureConfig struct {
	MaxSnaps       int
	SnapInterval   int // seconds
	HandleInterval int // seconds
	PsInterval     int // seconds
	Timeout        int // minutes
}

//*=========================[ Graphing ]===========================

const (
	PG_DIRECT_RELATIVE = 1 << 0
	PG_IPC_CHANNEL     = 1 << 1
	PG_SAME_PS_ACCESS  = 1 << 2
	PG_SAME_FILE_OP    = 1 << 4
	PG_SAME_FILE_WRITE = 1 << 5
	PG_SAME_FILE_READ  = 1 << 6
)

type Pool map[uint32]*ProcessSnapshot

type Graph struct {
	id      int
	mu      sync.RWMutex
	Members map[uint32]*ProcessNode
}

type GraphRegistry map[int]*Graph

type ProcessNode struct {
	ProcessId   uint32
	Process     *ProcessSnapshot
	Connections map[uint32]*Connection
}

type Connection struct {
	Type   Bitmask `json:"conn_type"`
	Weight int     `json:"weight"`
}

type Traversal struct {
	flags  Bitmask
	weight int
}

type AccessEntry struct {
	Object uint32  // type enum
	Name   string  // name of object
	Type   Bitmask // type of interaction
	Handle uint32
	Pid    uint32
	Params map[string]Parameter // extended object info
}

var g_ProcessTable *ProcessTable

type ProcessTable struct {
	mu    sync.RWMutex
	Table map[uint32]*ProcessSnapshot
}
