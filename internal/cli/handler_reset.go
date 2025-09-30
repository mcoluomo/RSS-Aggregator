package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/mcoluomo/RSS-Aggregator/internal/config"
)

func ResetHandler(s *config.State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("command does not accept any arguments")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	defer cancel()
	if err := s.Db.DeleteAllUsers(ctx); err != nil {
		return fmt.Errorf("failed deleting all users: %w", err)
	}

	s.StConfig.SetUser("[None]")

	fmt.Println("----------------------------------------")
	fmt.Println("Removing...")
	fmt.Println("	All users have been deleted!")
	fmt.Println("============= ALL REMOVED ==============")

	return nil
}
