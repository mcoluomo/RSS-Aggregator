package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/mcoluomo/RSS-Aggregator/internal/config"
	"github.com/mcoluomo/RSS-Aggregator/internal/rss"
)

const feedUrl string = "https://www.wagslane.dev/index.xml"

func AggHandler(s *config.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("Please provide the valid argument for this command: <command> [time_between_reqs]")
	}

	if len(cmd.Args) > 1 {
		return fmt.Errorf("command only takes one argumeant: <command> [time_between_reqs]")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	defer cancel()

	feedData, err := rss.FetchFeed(ctx, feedUrl)
	if err != nil {
		return fmt.Errorf("failed fetching url feed %w", err)
	}

	fmt.Println("----------------------------------------")
	for _, post := range feedData.Channel.Item {
		fmt.Printf("\n\nlisting feed post:\\n 【%v】", post)
	}

	fmt.Println("\n----------------------------------------")
	return nil
}
