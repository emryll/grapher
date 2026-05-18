#ifndef UTILS_H
#define UTILS_H

#include <windows.h>

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

#endif