package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/mcoluomo/RSS-Aggregator/internal/config"
)

func UserHandler(s *config.State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("command does not accept any arguments")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	defer cancel()
	users, err := s.Db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed fetching all users: %w", err)
	}
	fmt.Println("Listing all users...")
	if len(users) == 0 {
		fmt.Println("No users to list")
	}

	fmt.Println("---------------------------------")
	for _, user := range users {
		if user.Name == s.StConfig.Current_user_name {
			fmt.Println("* " + s.StConfig.Current_user_name + " (current)")
		} else {
			fmt.Println("* " + user.Name)
		}
	}
	fmt.Println("---------------------------------")
	return nil
}
