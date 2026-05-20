package main

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

// Add weight or new connection type to a connection.
func (c *Connection) Expand(flags Bitmask, weight int) {
	c.Weight += weight
	c.Type |= flags
}

func (c *Connection) Passes(rule Traversal) bool {
	return c.Type.HasFlags(rule.flags) && c.Weight >= rule.weight
}
