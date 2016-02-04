package cubbyhole

import "testing"

func TestCubbyholePut(t *testing.T) {
	c := Cubbyhole{}

	if c.Message != "" {
		t.Error(`Cubbyhole must be "" before the first put`)
	}

	c.Put("test!")
	if c.Message != "test!" {
		t.Errorf(`Cubbyhole should be "test!" but is: %s`, c.Message)
	}
}

func TestCubbyholeDrop(t *testing.T) {
	c := Cubbyhole{}

	c.Drop()
	if c.Message != "" {
		t.Errorf(`Cubbyhole.Drop didn't drop the message: %s`, c.Message)
	}

	c.Message = "foobar"

	c.Drop()
	if c.Message != "" {
		t.Errorf(`Cubbyhole.Drop didn't drop the message: %s`, c.Message)
	}
}

func TestCubbyholeLook(t *testing.T) {
	c := Cubbyhole{}

	if message := c.Look(); message != "" {
		t.Errorf(`Cubbyhole.Look didn't return an empty message: %s`, message)
	}

	c.Message = "NPA"
	if message := c.Look(); message != "NPA" {
		t.Errorf(`Cubbyhole.Look didn't return the correct message "NPA" != %s`, message)
	}
}

func TestCubbyholeGet(t *testing.T) {
	c := Cubbyhole{}

	if message := c.Get(); message != "" {
		t.Errorf(`Cubbyhole.Get didn't return an empty message: %s`, message)
	}

	c.Message = "NPA"
	if message := c.Get(); message != "NPA" {
		t.Errorf(`Cubbyhole.Look didn't return the correct message "NPA" != %s`, message)
	}

	if message := c.Get(); message != "" {
		t.Errorf(`Cubbyhole.Get didn't return an empty message: %s`, message)
	}
}
