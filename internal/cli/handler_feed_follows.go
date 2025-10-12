package cli

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mcoluomo/RSS-Aggregator/internal/config"
	"github.com/mcoluomo/RSS-Aggregator/internal/database"
)

func FollowHandler(s *config.State, cmd Command, user database.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	defer cancel()

	if len(cmd.Args) == 0 {
		return fmt.Errorf("Please provide the valid argument for this command: <command> 【[url]】")
	}

	if len(cmd.Args) > 1 {
		return fmt.Errorf("command only takes one argumeant: <command> [url]")
	}

	if !isValidUrl(cmd.Args[0]) {
		return fmt.Errorf("Please provide valid url: <command> 【[url]】")
	}

	feedId, err := s.Db.GetFeedId(ctx, cmd.Args[0])
	if err != nil {
		return fmt.Errorf("%w: failed fetching feed id", err)
	}

	user, _ = s.Db.GetUser(ctx, s.StConfig.Current_user_name)

	feedFollowParams := database.CreateFeedFollowParams{
		UserID:    user.ID,
		FeedID:    feedId,
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	feedFollow, err := s.Db.CreateFeedFollow(ctx, feedFollowParams)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			fmt.Printf("You are already following this feed: %s\n", cmd.Args[0])
			return nil
		}
		return fmt.Errorf("failed creating feed follow record: %w", err)
	}

	fmt.Println("----------------------------------------")
	fmt.Printf("* FeedName:      %s\n", feedFollow.FeedName)
	fmt.Printf("* FollowerName:   %s\n", s.StConfig.Current_user_name)
	fmt.Printf("* Created:       %v\n", feedFollow.CreatedAt.Time)
	fmt.Printf("* Updated:       %v\n", feedFollow.UpdatedAt.Time)
	fmt.Println("following feed", feedFollow.FeedName)

	return nil
}

func FeedFollowingHandler(s *config.State, cmd Command, user database.User) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("command does not accept any arguments")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	defer cancel()

	user, _ = s.Db.GetUser(ctx, s.StConfig.Current_user_name)

	feedFollows, err := s.Db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed getting feed follows: %w", err)
	}

	fmt.Printf("Getting all feeds that 【%s】 is following...\n", s.StConfig.Current_user_name)
	fmt.Println("---------------------------------")
	for _, feedFollow := range feedFollows {
		fmt.Printf("* Feed: 【%s】\n", feedFollow.FeedName)
	}

	fmt.Println("---------------------------------")

	return nil
}

func UnfollowFeedFollow(s *config.State, cmd Command, user database.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	defer cancel()

	if len(cmd.Args) == 0 {
		return fmt.Errorf("Please provide a url argument for this command: <command [url]")
	}

	if len(cmd.Args) > 1 {
		return fmt.Errorf("command only takes argumeants: <command> 【[url]】")
	}

	if !isValidUrl(cmd.Args[0]) {
		return fmt.Errorf("Please provide valid url: <command> 【[url]】")
	}

	feedId, err := s.Db.GetFeedId(ctx, cmd.Args[0])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("No feed found with that URL. Please add the feed first.")
		}
		return fmt.Errorf("%w: failed fetching feed id", err)
	}

	feedFollowParams := database.DeleteFeedFollowRowParams{
		UserID: getUserId(s),
		FeedID: feedId,
	}

	if err = s.Db.DeleteFeedFollowRow(ctx, feedFollowParams); err != nil {
		return fmt.Errorf("%w: failed deleting feed follow row for 【%s】", err, cmd.Args[0])
	}

	return nil
}
