#ifndef UTILS_H
#define UTILS_H

#include <windows.h>

typedef enum {
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
} PARAMETER_TYPE;

typedef enum {
	TYPE_UNKNOWN
	TYPE_PROCESS
	TYPE_THREAD
	TYPE_FILE
	TYPE_SEMAPHORE
	TYPE_EVENT
	TYPE_MUTEX
	TYPE_SYMLINK
	TYPE_PIPE
	TYPE_SECTION
	TYPE_TOKEN
	TYPE_DIRECTORY
	TYPE_DBG_OBJECT
	TYPE_DEVICE
	TYPE_DRIVER
	TYPE_DESKTOP
	TYPE_WORKER_FACTORY
} OBJECT_TYPE;

typedef struct {
    DWORD Pid; // owner pid
    DWORD Type; // object type
    DWORD Access; 
    DWORD Handle; // value only valid in that process (used for unique id)
    BYTE* Params; // contains object specific info
} HANDLE_ENTRY;

HANDLE_ENTRY* GetGlobalHandleTable(size_t*);

// internal helpers
DWORD GetHandleObjectType(HANDLE);
BYTE* GetHandleParameters(HANDLE, DWORD, size_t*);
BYTE* CreateParameter(char*, DWORD, DWORD, size_t*);
BYTE* BuildParameter(size_t*, DWORD, const char*, ...);

// function prototype for cast to call it
typedef NTSTATUS (NTAPI *NQSI)(ULONG, PVOID, ULONG, PULONG);
typedef NTSTATUS (NTAPI *NQO)(HANDLE, DWORD, PVOID, ULONG, PULONG);
static NQO NtQueryObject;

#endif