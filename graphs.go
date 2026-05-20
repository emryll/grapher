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
