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
	"github.com/mcoluomo/RSS-Aggregator/internal/rss"
)

func PrintFeedsHandler(s *config.State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("command does not accept any arguments")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	defer cancel()

	feeds, err := s.Db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("PrintFeedsHandler failed fetching feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}
	fmt.Printf("Found %d feeds:\n", len(feeds))

	fmt.Println("listing all feeds...")

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
	fmt.Println("----------------------------------------")
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name)
	fmt.Println("----------------------------------------")
}

func AddFeedHandler(s *config.State, cmd Command, user database.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	defer cancel()

	if len(cmd.Args) == 0 {
		return fmt.Errorf("Please provide the valid argument for this command: <command [url]")
	}

	if len(cmd.Args) > 2 {
		return fmt.Errorf("command only takes two argumeants: <command> 【[feedName]】 【[url]】")
	}

	if strings.TrimSpace(cmd.Args[0]) == "" {
		return fmt.Errorf("Please provide valid feedName argument for this command: <command> 【[feedName]】 [url]")
	}

	if len(cmd.Args) < 2 {
		return fmt.Errorf("Please provide all arguments for this command: <command> 【[feedName]】 [url]")
	}

	if !isValidUrl(cmd.Args[1]) {
		return fmt.Errorf("Please provide valid url: <command> [feedName] 【[url]】")
	}

	user, _ = s.Db.GetUser(ctx, s.StConfig.Current_user_name)

	newFeed := database.CreateFeedParams{
		ID:            uuid.New(),
		CreatedAt:     sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:     sql.NullTime{Time: time.Now(), Valid: true},
		Name:          cmd.Args[0],
		Url:           cmd.Args[1],
		UserID:        user.ID,
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	feed, err := s.Db.CreateFeed(ctx, newFeed)
	if err != nil {
		return fmt.Errorf("%w: failed creating feed", err)
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

func ScrapeFeedsHander(s *config.State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("command doesnt take any arguments")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	defer cancel()

	nextFeed, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("%w: failed fetching next feed", err)
	}
	if err = s.Db.MarkFeedFetched(ctx, nextFeed.ID); err != nil {
		return fmt.Errorf("%w: failed trying to mark feed", err)
	}

	feedDate, err := rss.FetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return fmt.Errorf("%w: ScrapeFeedsHander failed fetching feed", err)
	}
	for _, feedItem := range feedDate.Channel.Item {
		fmt.Printf("* FeedItemTitle: 【%s】\n", feedItem.Title)
	}
	return nil
}
