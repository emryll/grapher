package main

//#include "utils.h"
import "C"

import (
	"encoding/binary"
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Parse the global handle table and register new handles.
func ParseHandleTable() {
	var handleCount C.size_t
	cHandleEntries := C.GetGlobalHandleTable(&handleCount)
	handleTable := unsafe.Slice((*HandleEntry)(unsafe.Pointer(cHandleEntries)), int(handleCount))

	for _, handle := range handleTable {
		handle.RegisterHandle()
	}
}

// Register a handle as an interaction, adding it
// to the object access registry and graphs if new.
func (h HandleEntry) RegisterHandle() {
	entry := h.ConvertToAccessEntry()
	entry.RegisterInteraction()
}

//TODO: is there any point in having a cache? i dont think so

func (h HandleEntry) ConvertToAccessEntry() AccessEntry {
	var entry AccessEntry
	entry.Object = h.ObjType
	entry.Handle = h.Handle
	entry.Pid = h.Pid

	if entry.Object == OBJECT_TYPE_PROCESS {
		pidParam := h.GetParameter("Pid")
		if !pidParam.Empty() {
			entry.Name = fmt.Sprintf("%d", binary.LittleEndian.Uint32(pidParam.Buffer))
		}
	} else {
		nameParam := h.GetParameter("Name")
		if !nameParam.Empty() {
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

func (h HandleEntry) GetParameter(name string) Parameter {
	if param, exists := h.Params[name]; exists {
		return param
	}
	return Parameter{}
}
