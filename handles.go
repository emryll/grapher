package main

import (
	"encoding/binary"
	"fmt"

	"golang.org/x/sys/windows"
)

func ParseHandleTable() {
	//TODO: get handle table
	//TODO: iterate handle table
	//TODO: for each handle call RegisterHandle()
}

//TODO: is there any point in having a cache? i dont think so

func (h HandleEntry) ConvertToAccessEntry() AccessEntry {
	var entry AccessEntry
	entry.Object = h.ObjType
	entry.Handle = h.Handle
	entry.Pid = h.Pid

	if entry.Object == OBJECT_TYPE_PROCESS {
		pidParam := h.GetParameter("Pid")
		if pidParam != nil {
			entry.Name = fmt.Sprintf("%d", binary.LittleEndian.Uint32(pidParam.Buffer))
		}
	} else {
		nameParam := h.GetParameter("Name")
		if nameParam != nil {
			entry.Name = GetAnsiValue(nameParam.Buffer)
		}
	}

	//* get interaction type
	switch entry.Object {
	case OBJECT_TYPE_PROCESS:
		entry.Type |= ANY_ACCESS
		if h.Access.HasFlags(windows.PROCESS_VM_OPERATION) {
			entry.Type |= PS_READ_MEM
		}
		if h.Access.HasFlags(windows.PROCESS_CREATE_THREAD) {
			entry.Type |= PS_CREATE_THREAD
		}
	case OBJECT_TYPE_THREAD:
	case OBJECT_TYPE_FILE:
		entry.Type |= ANY_ACCESS
		if h.Access.HasFlags(windows.FILE_GENERIC_READ) {
			entry.Type |= FILE_READ
		}
		if h.Access.HasFlags(windows.FILE_GENERIC_WRITE) {
			entry.Type |= FILE_WRITE
		}
	case OBJECT_TYPE_EVENT:
		entry.Type = ANY_ACCESS
	case OBJECT_TYPE_SEMAPHORE:
		entry.Type = ANY_ACCESS
	case OBJECT_TYPE_MUTEX:
		entry.Type = ANY_ACCESS
	}
	return entry
}
