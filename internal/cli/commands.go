package cli

import (
	"fmt"
	"os"

	"github.com/mcoluomo/RSS-Aggregator/internal/config"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Handlers map[string]func(*config.State, Command) error
}

func (c *Commands) Run(s *config.State, cmd Command) error {
	handler, exist := c.Handlers[cmd.Name]
	if !exist {
		return fmt.Errorf("\nunknown command")
	}

	return handler(s, cmd)
}

func (c *Commands) Register(name string, f func(*config.State, Command) error) {
	c.Handlers[name] = f
}

func LoginHandler(s *config.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no arguments were given")
	}
	s.StConfig.SetUser(os.Args[2])
	s.StConfig.Current_user_name = cmd.Args[0]
	fmt.Printf("⇒ %s\nlogin with the given username was successful\n", cmd.Args[0])

	return nil
}
