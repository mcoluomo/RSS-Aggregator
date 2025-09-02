package cli

import (
	"fmt"

	"github.com/mcoluomo/RSS-Aggregator/internal/config"
)

type Command struct {
	Name      string
	Arguments []string
}

type Commands struct {
	cmd map[string]func(*config.State, Command) error
}

func (c *Commands) run(s *config.State, cmd Command) error {
	return nil
}

func (c *Commands) register(name string, f func(*config.State, Command) error) {
	return
}

func SetUserHandler(s *config.State, cmd Command) error {
	if len(cmd.Arguments) == 0 {
		return fmt.Errorf("no arguments Error")
	}
	s.StConfig.Current_user_name = cmd.Arguments[0]
	fmt.Printf("current user was set ⇒ %s\n", cmd.Arguments[0])

	return nil
}
