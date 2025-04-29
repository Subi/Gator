package main

const url = "https://www.wagslane.dev/index.xml"

type Command struct {
	name string
	args []string
}

type Commands struct {
	opts map[string]func(*state, Command) error
}

func (c *Commands) run(s *state, cmd Command) error {
	err := c.opts[cmd.name](s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *Commands) register(name string, f func(*state, Command) error) {
	c.opts[name] = f
}
