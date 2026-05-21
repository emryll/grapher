package main

import "sync"

var (
	ID_COUNTER             = 1
	g_ObjectAccessRegistry *ObjectAccessRegistry
	g_GraphRegistry        GraphRegistry
)

// Returns resulting bitmask after stripping
func (c *Connection) Strip(flags Bitmask, weight int) Bitmask {
	c.Weight -= weight
	c.Type &^= flags // strip flags
	return c.Type
}

// Add weight or new connection type to a connection
func (c *Connection) Expand(flags Bitmask, weight int) {
	c.Weight += weight
	c.Type |= flags
}

// Check if connection passes given traversal rule
func (c *Connection) Passes(rule Traversal) bool {
	return c.Type.HasFlags(rule.flags) && c.Weight >= rule.weight
}

func (snap Snapshot) CreatePools() []Pool {
	var pools []Pool
	for _, g := range snap.Graphs {
		pools = append(pools, g.CreatePools()...)
	}
	return pools
}

func (g *GraphSnapshot) CreatePools() []Pool {
	//TODO: traverse graph
}

//*====================[ Object Access Lookup ]===================

// Lookup table for object interactions
// 500 000 entries would be around 32MB
type ObjectAccessRegistry struct {
	mu sync.RWMutex // used internally in methods
	// process -> object type -> name -> entry
	ProcessLookup map[uint32]map[ProcessAccessKey][]*AccessEntry // array is for anon objects
	// object type -> name -> process -> entry
	ObjectLookup map[uint32]map[ObjectAccessKey][]*AccessEntry
}

// With the triple nested map, amount of maps grows very quickly.
// To fix this issue, the structure is partially flattened.
// Instead of a triple map its a double map with a struct key,
// which has a very big effect on the amount of maps created.

// This key struct is made to flatten ProcessLookup
type ProcessAccessKey struct {
	ObjType uint32
	Name    string
}

// This key struct is made to flatten ObjectLookup
type ObjectAccessKey struct {
	Pid  uint32
	Name string
}

// Register a tracked interaction to the registry and graph.
// Before running this you should check that the interaction
// is something used in the relationship graphing (for efficiency).
func (entry *AccessEntry) RegisterInteraction() {
	g_ObjectAccessRegistry.AddEntry(*entry, entry.Pid)

	connections := entry.GetNewConnections(entry.Pid)
	if len(connections) == 0 {
		return
	}

	graph := GetGraph(entry.Pid)
	for _, newConn := range connections {
		graph.AddConnection(entry.Type, entry.GetWeight(), entry.Pid, newConn)
	}
}

// Add an interaction to the registry. Updates existing if one exists.
func (reg *ObjectAccessRegistry) AddEntry(entry AccessEntry, pid uint32) {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	// check that maps are initialized (avoid panic)
	if reg.ProcessLookup[pid] == nil {
		reg.ProcessLookup[pid] = make(map[ProcessAccessKey][]*AccessEntry)
	}
	if reg.ObjectLookup[entry.Object] == nil {
		reg.ObjectLookup[entry.Object] = make(map[ObjectAccessKey][]*AccessEntry)
	}

	objectKey := entry.CreateObjectKey()
	processKey := entry.CreateProcessKey()

	// check if entry exists, update existing if does
	entries := reg.FindByProcess([]uint32{pid}, []uint32{entry.Object}, entry.Name)
	if len(entries) > 0 {
		for _, ent := range entries {
			if ent.Handle != entry.Handle {
				continue
			}
			ent.Type |= entry.Type
			return
		}
	}

	e := entry // just to be safe with uniqueness...
	reg.ProcessLookup[pid][processKey] = append(reg.ProcessLookup[pid][processKey], &e)
	reg.ObjectLookup[pid][objectKey] = append(reg.ObjectLookup[pid][objectKey], &e)
}

// Delete all interaction entries under a certain process.
// This function should be called when a process exits, to cleanup.
func (reg *ObjectAccessRegistry) RemoveEntriesByProcess(pid uint32) {
	reg.mu.Lock()
	defer reg.mu.Unlock()

	if len(reg.ProcessLookup[pid]) == 0 {
		return
	}

	// remove entries
	for psKey, entries := range reg.ProcessLookup[pid] {
		for _, entry := range entries {
			objKey := ObjectAccessKey{Name: psKey.Name, Pid: pid}
			if len(reg.ObjectLookup[uint32(entry.Type)]) > 0 {
				delete(reg.ObjectLookup[uint32(entry.Type)], objKey)
			}
		}
	}
	delete(reg.ProcessLookup, pid)
}

// Find all corresponding entries based on the acting process.
// @param  pids    The type of object to be accessed.
// @param  objs   (optional) Object type filter for entries.
// @param  names  (optional) Whitelist for object names.
// @return          All matching object access entries.
func (reg *ObjectAccessRegistry) FindByProcess(pids []uint32, objs []uint32, names ...string) []*AccessEntry {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	if len(pids) == 0 {
		return nil
	}

	var (
		entries    []*AccessEntry
		typeFilter = make(map[uint32]bool)
		nameFilter = make(map[string]bool)
		pidFilter  = make(map[uint32]bool)
	)

	for _, val := range pids {
		pidFilter[val] = true
	}
	for _, val := range objs {
		typeFilter[val] = true
	}
	for _, val := range names {
		nameFilter[val] = true
	}

	for _, pid := range pids {
		if len(reg.ProcessLookup[pid]) == 0 {
			continue
		}
		for objKey, accessEntries := range reg.ProcessLookup[pid] {
			if len(objs) > 0 && !typeFilter[objKey.ObjType] {
				continue
			}
			if len(names) > 0 && !nameFilter[objKey.Name] {
				continue
			}
			entries = append(entries, accessEntries...)
		}
	}
	return entries
}

// Find all corresponding entries based on object description.
// @param  objectType    The type of object to be accessed.
// @param  interaction   (optional) Bitmask describing type of interaction.
// @param  names         (optional) Whitelist for object names.
// @return               All matching object access entries.
func (reg *ObjectAccessRegistry) FindByObject(objectType Bitmask, interaction Bitmask, names ...string) []*AccessEntry {
	if len(reg.ObjectLookup[uint32(objectType)]) == 0 {
		return nil
	}
	var (
		result     []*AccessEntry
		nameFilter = make(map[string]bool)
	)

	for _, name := range names {
		nameFilter[name] = true
	}

	for key, entries := range reg.ObjectLookup[uint32(objectType)] {
		if len(names) > 0 && !nameFilter[key.Name] {
			continue
		}
		for _, entry := range entries {
			if entry.Type.HasFlags(interaction) {
				result = append(result, entry)
			}
		}
	}
	return result
}

func (entry *AccessEntry) GetNewConnections(pid uint32) []uint32 {
	if !entry.IsTrackedInteraction() {
		return nil
	}

	var connections []uint32
	graph := GetGraph(entry.Pid)
	graph.mu.RLock()
	defer graph.mu.RUnlock()

	entries := g_ObjectAccessRegistry.FindByObject((Bitmask)(entry.Object), entry.Type, entry.Name)
	for _, ent := range entries {
		if graph.Members[entry.Pid].Connections[ent.Pid] == nil {
			connections = append(connections, ent.Pid)
		}
	}
	return connections
}

func (entry *AccessEntry) CreateObjectKey() ObjectAccessKey {
	return ObjectAccessKey{Name: entry.Name, Pid: entry.Pid}
}

func (entry *AccessEntry) CreateProcessKey() ProcessAccessKey {
	return ProcessAccessKey{Name: entry.Name, ObjType: entry.Object}
}
