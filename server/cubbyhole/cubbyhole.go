package cubbyhole

type Cubbyhole struct {
	Message string
}

func (c *Cubbyhole) Put(message string) {
	c.Message = message
}

func (c *Cubbyhole) Drop() {
	c.Message = ""
}

func (c *Cubbyhole) Look() string {
	return c.Message
}

func (c *Cubbyhole) Get() string {
	message := c.Look()
	c.Drop()
	return message
}
