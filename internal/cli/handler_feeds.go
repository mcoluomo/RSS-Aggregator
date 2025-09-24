package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/mcoluomo/RSS-Aggregator/internal/config"
	"github.com/mcoluomo/RSS-Aggregator/internal/database"
)

func PrintFeedsHandler(s *config.State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("command does not accept any arguments")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	defer cancel()

	users, err := s.Db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("PrintFeedHandler failed fetching users: %w", err)
	}

	if len(users) == 0 {
		return fmt.Errorf("no users in database to list feeds from")
	}

	feeds, err := s.Db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("PrintFeedsHandler failed fetching feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}
	fmt.Printf("Found %d feeds:\n", len(feeds))

	fmt.Println("listing all feeds...\n\nuser_name | feed_name | feed_url")

	for _, feed := range feeds {
		user, err := s.Db.GetUserById(ctx, feed.UserID)
		if err != nil {
			return fmt.Errorf("failed getting user: %w", err)
		}
		printFeed(feed, user)
	}
	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name)
}
