package cli

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/mcoluomo/RSS-Aggregator/internal/config"
	"github.com/mcoluomo/RSS-Aggregator/internal/database"
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

	if err = s.Db.MarkFeedFetched(ctx, nextFeed.ID); err != nil {
		return fmt.Errorf("%w: failed trying to mark feed: %v", err, nextFeed.Name)
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] Fetching feed: %s (%s)\n", now, nextFeed.Name, nextFeed.Url)
	fmt.Printf("[%s] Marked feed as fetched.\n", now)

	feedDate, err := rss.FetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return err
	}
	fmt.Printf("[%s] Fetched %d items from feed.\n", now, len(feedDate.Channel.Item))

	log.Println("Fechted FEED")
	var saved, skipped int
	for _, feedItem := range feedDate.Channel.Item {

		fmt.Println("-------------------------------------------")
		if strings.TrimSpace(feedItem.Title) == "" {
			skipped += 1
			continue
		}
		publicationTime, err := time.Parse(time.RFC1123Z, feedItem.PubDate)
		if err != nil {
			log.Panicf("%v: failed converting %s into time.time", err, feedItem.PubDate)
		}
		log.Println("converted date string into time.time")
		feedItemId := nextFeed.ID
		log.Println("got feed id")

		var createPostParams database.CreatePostParams = database.CreatePostParams{
			CreatedAt:   sql.NullTime{Time: time.Now(), Valid: true},
			UpdatedAt:   sql.NullTime{Time: time.Now(), Valid: true},
			Title:       feedItem.Title,
			Url:         feedItem.Link,
			Description: sql.NullString{String: feedItem.Description, Valid: true},
			PublishedAt: sql.NullTime{Time: publicationTime, Valid: true},
			FeedID:      feedItemId,
		}
		log.Println("created post params")

		_, err = s.Db.CreatePost(ctx, createPostParams)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				log.Printf("[%s] Skipped duplicate post: %s\n", now, feedItem.Title)
				skipped++
				continue
			}
			fmt.Printf("[%s] ERROR: failed creating post: %v\n", now, err)
			continue
		}
		fmt.Printf("[%s] Saved post: %s\n", now, feedItem.Title)
		saved++
	}
	fmt.Printf("[%s] Finished processing feed: %s. %d new posts saved, %d duplicates skipped.\n", now, nextFeed.Name, saved, skipped)
	return nil
}

func BrowseFeedsHandler(s *config.State, cmd Command) error {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	defer cancel()

	if len(cmd.Args) == 0 {
		cmd.Args = append(cmd.Args, "2")
	}
	if len(cmd.Args) > 1 {
		return fmt.Errorf("command only takes one argumeant: <command> [limit_if_posts_as_number]")
	}

	if !containsOnlyNumericDigits(cmd.Args[0]) {
		return fmt.Errorf("invalid numberic argument: <command> [limit_if_posts_as_number]")
	}

	user, _ := s.Db.GetUser(ctx, s.StConfig.Current_user_name)
	numOfPosts, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		log.Panicf("%v: failed converting %s into integer", err, cmd.Args[0])
	}

	getUserPostParams := database.GetUserPostsParams{
		UserID: user.ID,
		Limit:  int32(numOfPosts),
	}
	userPosts, err := s.Db.GetUserPosts(ctx, getUserPostParams)
	if err != nil {
		return fmt.Errorf("%w: failed fetching user 【%s】 posts", err, s.StConfig.Current_user_name)
	}
	for i, postRow := range userPosts {
		fmt.Printf("Post #%d\n", i+1)
		fmt.Printf("Feed:        %s\n", postRow.Name)
		fmt.Printf("Title:       %s\n", postRow.Title)
		if postRow.PublishedAt.Valid {
			fmt.Printf("Published:   %s\n", postRow.PublishedAt.Time.Format("2006-01-02"))
		}
		fmt.Printf("URL:         %s\n", postRow.Url)
		if postRow.Description.Valid && len(postRow.Description.String) > 0 {
			desc := postRow.Description.String
			if len(desc) > 100 {
				desc = desc[:100] + "..."
			}
			fmt.Printf("Description: %s\n", desc)
		}
		fmt.Println("------------------------------------------------------------")
	}

	return nil
}

func containsOnlyNumericDigits(numericStr string) bool {
	for _, char := range numericStr {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}
