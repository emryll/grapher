package main

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// main function which will walk all processes
// and add any new connections to the graph.
func ScanProcesses() error {
	//* get a process snapshot
	handle, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(handle)
	entry := windows.ProcessEntry32{Size: uint32(unsafe.Sizeof(windows.ProcessEntry32{}))}

	//* walk the processes
	for {
		err = windows.Process32Next(handle, &entry)
		if err != nil {
			if err.Error() == "There are no more files." {
				return nil
			}
			return err
		}

		RegisterProcess(&entry)
	}
}

func CreateProcessEntry(pid uint32, parent ...uint32) ProcessSnapshot {
	entry := ProcessSnapshot{
		Connections: make(map[uint32]*Connection),
		IsSigned:    IsSigned(pid),
		ProcessId:   pid,
	}

	var handleFound bool
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		fmt.Printf("[WARNING] Failed to open process %d: %v\n", pid, err)
	} else {
		handleFound = true
		defer windows.CloseHandle(handle)
	}

	if handleFound {
		elevated, err := IsProcessElevated(handle)
		if err == nil {
			entry.IsElevated = elevated
		}
		path, err := GetProcessExecutable(handle)
		if err == nil {
			entry.Name = path
		}
	}

	var parentFound bool
	if len(parent) > 0 {
		entry.ParentPid = parent[0]
		parentFound = true
	} else if handleFound {
		ppid, err := GetParentPid(handle)
		if err == nil {
			entry.ParentPid = ppid
			parentFound = true
		}
	}

	if parentFound {
		parentHandle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, entry.ParentPid)
		if err == nil {
			parentName, err := GetProcessExecutable(parentHandle)
			if err == nil {
				entry.ParentName = parentName
			}
		}
	}
	return entry
}

func CreateProcessTable() *ProcessTable {
	return &ProcessTable{Table: make(map[uint32]*ProcessSnapshot)}
}

func (ps *ProcessTable) LookupProcess(pid uint32) *ProcessSnapshot {
	if ps.Table == nil {
		return nil
	}
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	if process, exists := ps.Table[pid]; exists {
		return process
	}
	return nil
}

func (ps *ProcessTable) AddProcess(process *ProcessSnapshot) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.Table == nil {
		ps.Table = make(map[uint32]*ProcessSnapshot)
	}
	ps.Table[process.ProcessId] = process
}

func (ps *ProcessTable) RemoveProcess(pid uint32) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.Table == nil {
		return
	}
	if _, exists := ps.Table[pid]; exists {
		delete(ps.Table, pid)
	}
}
