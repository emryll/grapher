package main

import (
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
		IsElevated:  IsElevated(pid),
		IsSigned:    IsSigned(pid),
		ProcessId:   pid,
	}

	path, err := GetProcessExecutable(pid)
	if err == nil {
		entry.Name = path
	}

	if len(parent) > 0 {
		entry.ParentPid = parent[0]
	} else {
		handle, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION, false, pid)
		if err == nil {
			defer windows.CloseHandle(handle)
			ppid, err := GetParentPid(handle)
			if err == nil {
				entry.ParentPid = ppid
			}
		}
	}

	parentName, err := GetProcessExecutable(entry.ParentPid)
	if err == nil {
		entry.ParentName = parentName
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
