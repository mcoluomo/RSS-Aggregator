package cli

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mcoluomo/RSS-Aggregator/internal/config"
	"github.com/mcoluomo/RSS-Aggregator/internal/database"
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

func AddFeedHandler(s *config.State, cmd Command, user database.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	defer cancel()

	if len(cmd.Args) > 2 {
		return fmt.Errorf("command only takes two argumeants: <command> 【[feedName]】 【[url]】")
	}

	if len(cmd.Args) < 2 {
		return fmt.Errorf("Please provide valid arguments for this command: <command> 【[feedName]】 【[url]】")
	}

	if strings.TrimSpace(cmd.Args[0]) == "" {
		return fmt.Errorf("Please provide valid feedName argument for this command: <command> 【[feedName]】 [url]")
	}

	if !isValidUrl(cmd.Args[1]) {
		return fmt.Errorf("Please provide valid url: <command> [feedName] 【[url]】")
	}

	user, _ = s.Db.GetUser(ctx, s.StConfig.Current_user_name)

	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
	}

	feed, err := s.Db.CreateFeed(ctx, newFeed)
	if err != nil {
		return fmt.Errorf("failed creating feed")
	}

	feedFollowParams := database.CreateFeedFollowParams{
		UserID:    feed.UserID,
		FeedID:    feed.ID,
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	fmt.Println("---------------------------------")
	feedFollow, err := s.Db.CreateFeedFollow(ctx, feedFollowParams)
	if err != nil {
		return fmt.Errorf("failed creating feed follow record: %w", err)
	}
	fmt.Printf("	* FeedName:      %s\n", feedFollow.FeedName)
	fmt.Printf("	* CurrentUser:   %s\n", s.StConfig.Current_user_name)

	fmt.Println("---------------------------------")

	fmt.Printf("successfuly added feed for user: 【%s】\n", s.StConfig.Current_user_name)
	fmt.Printf("%v\n", feed.CreatedAt.Time)

	return nil
}
