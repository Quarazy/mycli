package cmd

import "log"

// Command executes an action
type Runnable interface {
	Run(args []string)
}

// NoArgCommand is a Command accepting no args
type NoArgCommand struct {
	run func()
}

func (c *NoArgCommand) Run(args []string) {
	c.run()
}

// OneArgCommand is a Command accepting only 1 arg
type OneArgCommand struct {
	run func(arg string)
}

func (c *OneArgCommand) Run(args []string) {
	c.run(args[0])
}

// TwoArgCommand is a Command accepting only 2 args
type TwoArgCommand struct {
	run func(arg1, arg2 string)
}

func (c *TwoArgCommand) Run(args []string) {
	c.run(args[0], args[1])
}

// Commands are a mapping of all the commands
type Commands struct {
	commands map[string]Runnable
}

// New creates a new Commands struct
func New() *Commands {
	return &Commands{
		commands: make(map[string]Runnable),
	}
}

// Register adds commands
func (c *Commands) Register(name string, r Runnable) {
	c.commands[name] = r
}

func (c *Commands) Run(args []string) {
	if len(args) < 2 {
		log.Fatal("No command provided and no default")
	}

	// Arguments ignoring the bin name
	name := args[1]

	if r, ok := c.commands[name]; ok {
		r.Run(args[2:])
	} else {
		// TODO(quarazy): Show help mesasge or something.
		// Would be cool to offer a suggestion if we detect a typo
		log.Fatalf("Command: %s doesn't exist", name)
	}
}
