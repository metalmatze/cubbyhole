package cubbyhole

// Cubbyhole is a secret place where you can put a message
type Cubbyhole struct {
	Message string
}

// Put places your message into the cubbyhole
func (c *Cubbyhole) Put(message string) {
	c.Message = message
}

// Drop removes your message without returning it first
func (c *Cubbyhole) Drop() {
	c.Message = ""
}

// Look gets the message
func (c *Cubbyhole) Look() string {
	return c.Message
}

// Get returns your message and then removes it
func (c *Cubbyhole) Get() string {
	message := c.Look()
	c.Drop()
	return message
}
