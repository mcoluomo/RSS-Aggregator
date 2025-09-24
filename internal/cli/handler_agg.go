package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/mcoluomo/RSS-Aggregator/internal/config"
	"github.com/mcoluomo/RSS-Aggregator/internal/rss"
)

func AggHandler(s *config.State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("command does not accept any arguments")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	defer cancel()

	feedData, err := rss.FetchFeed(ctx, feedUrl)
	if err != nil {
		return fmt.Errorf("failed fetching url feed %w", err)
	}

	fmt.Printf("\nlisting feed data:\n %v", feedData)

	return nil
}
