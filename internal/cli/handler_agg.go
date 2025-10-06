package cli

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mcoluomo/RSS-Aggregator/internal/config"
	"github.com/mcoluomo/RSS-Aggregator/internal/rss"
)

func AggHandler(s *config.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Please provide the valid argument for this command: <command> [time_between_reqs]")
	}

	if len(cmd.Args) > 1 {
		return fmt.Errorf("command only takes one argumeant: <command> [time_between_reqs]")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("%w: invalid duration argument", err)
	}

	fmt.Printf("Collecting feeds every %v...\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	defer ticker.Stop()

	for range ticker.C {
		ScrapeFeedsHander(s)
	}

	return nil
}

func ScrapeFeedsHander(s *config.State) error {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	defer cancel()

	nextFeed, err := s.Db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("%w: failed fetching next feed", err)
	}

	fmt.Println("-------------------------------------------")
	log.Println("Found a feed to fetch!")
	if err = s.Db.MarkFeedFetched(ctx, nextFeed.ID); err != nil {
		return fmt.Errorf("%w: failed trying to mark feed: %v", err, nextFeed.Name)
	}

	feedDate, err := rss.FetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return fmt.Errorf("%w: ScrapeFeedsHander failed fetching feed", err)
	}

	fmt.Println("-------------------------------------------")
	var countOfBlankPosts int
	for _, feedItem := range feedDate.Channel.Item {
		if strings.TrimSpace(feedItem.Title) == "" {
			countOfBlankPosts += 1
			continue
		}
		log.Printf("Found post: %s\n", feedItem.Title)
	}
	log.Printf("Feed 【%s】 collected, %v posts found", nextFeed.Name, len(feedDate.Channel.Item)-countOfBlankPosts)
	return nil
}
