# grapher
This is a tool I am using for data collection and viewing as part of my research on detecting multi-process malware.
With this you can capture graph snapshots on a Windows x64 system, and then later you can load these collected snapshots and view/query them via the interactive command-line interface.

## Usage
To build the tool, simply run `go build`.
Note that you need to have both the compiler and a C compiler.
### Capturing data
To begin capturing data, run
```
grapher.exe capture [flags]
```

You can use flags to alter the config, which defines intervals and max snapshot count.
Available flags are:
```
-m   Maximum amount of snapshots to take.
-t   Automatic capture timeout (minutes).
-p   Process refresh interval (seconds).
-h   Handle refresh interval (seconds).
```

### Viewing data
To enter the interactive command-line interface used to view captured data, run
```
grapher.exe <path>
```
, where `path` is the path to the capture folder.

**Not ready to be used at this point in time!**
