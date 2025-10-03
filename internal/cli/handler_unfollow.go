package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/mcoluomo/RSS-Aggregator/internal/config"
	"github.com/mcoluomo/RSS-Aggregator/internal/database"
)

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
