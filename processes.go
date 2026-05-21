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

func RegisterProcess(entry *windows.ProcessEntry32) {
	// entry.ProcessID , entry.ParentProcessID , entry.ExeFile
}
