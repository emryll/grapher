# grapher
This is a tool I am writing as part of my research on multi-process correlation for detecting multi-process malware. This tool allows you to capture process relationship graph data over a period of time on Windows systems, and to later view previously collected process relationship data, using an interactive command-line interface or exporting the data.

> **Not ready to be used at this point in time!**

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
To enter the interactive commandline interface used to view captured data, run
```
grapher.exe <path>
```
, where `path` is the path to the capture folder.

Available commands in the interactive CLI are:
```
help [command]  Show information about commands.
exit            Exit the commandline interface.
state           Show the current state, in regards to snaps.
select <snap>   Select a snapshot for analysis (by name).
overview        Get a quick overview about session and selected snap.
graphs          View the graphs in the currently selected snap.
pools [flags]   Split graphs into subsets based on a traversal rule.
find <min>      Find all processes with more than min connections.
```
