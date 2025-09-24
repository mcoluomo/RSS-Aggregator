package cli

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mcoluomo/RSS-Aggregator/internal/config"
	"github.com/mcoluomo/RSS-Aggregator/internal/database"
)

func FollowHandler(s *config.State, cmd Command, user database.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	defer cancel()

	users, err := s.Db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("FollowHandler failed fetching users: %w", err)
	}

	if len(users) == 0 {
		return fmt.Errorf("no users in database to create a feed follow record")
	}

	if len(cmd.Args) == 0 {
		return fmt.Errorf("no argument was given")
	}

	if len(cmd.Args) > 1 {
		return fmt.Errorf("command only takes one argumeant: <command> [url]")
	}

	if !isValidUrl(cmd.Args[0]) {
		return fmt.Errorf("Please provide valid url: <command> 【[url]】")
	}

	feedId, err := s.Db.GetFeedId(ctx, cmd.Args[0])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed fetching feed id: %w", err)
		}
	}

	feedFollowParams := database.CreateFeedFollowParams{
		UserID:    fetchUserId(users, s),
		FeedID:    feedId,
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	feedFollow, err := s.Db.CreateFeedFollow(ctx, feedFollowParams)
	if err != nil {
		return fmt.Errorf("failed creating feed follow record: %w", err)
	}

	fmt.Printf("* FeedName:      %s\n", feedFollow.FeedName)
	fmt.Printf("* FollowerName:   %s\n", s.StConfig.Current_user_name)
	fmt.Printf("* Created:       %v\n", feedFollow.CreatedAt.Time)
	fmt.Printf("* Updated:       %v\n\n", feedFollow.UpdatedAt.Time)
	fmt.Printf("following feed %s", feedFollow.FeedName)

	return nil
}

func FeedFollowingHandler(s *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("command does not accept any arguments")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	defer cancel()

	users, err := s.Db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed getting user: %s %w", s.StConfig.Current_user_name, err)
	}

	// see if we still need to use this
	if len(users) == 0 || s.StConfig.Current_user_name == "[None]" {
		return fmt.Errorf("no users in database to create a feed follow record")
	}

	feedFollows, err := s.Db.GetFeedFollowsForUser(ctx, fetchUserId(users, s))
	if err != nil {
		return fmt.Errorf("failed getting feed follows: %w", err)
	}

	fmt.Printf("Getting all feeds that 【%s】 is following...\n", s.StConfig.Current_user_name)
	for _, feedFollow := range feedFollows {
		fmt.Printf("* Feed: 【%s】\n", feedFollow.FeedName)
	}

	return nil
}
