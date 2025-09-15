package cli

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mcoluomo/RSS-Aggregator/internal/config"
	"github.com/mcoluomo/RSS-Aggregator/internal/database"
	"github.com/mcoluomo/RSS-Aggregator/internal/rss"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Handlers map[string]func(*config.State, Command) error
}

const feedUrl string = "https://www.wagslane.dev/index.xml"

func (c *Commands) Run(s *config.State, cmd Command) error {
	handler, exist := c.Handlers[cmd.Name]
	if !exist {
		return fmt.Errorf("\nunknown command")
	}

	return handler(s, cmd)
}

func (c *Commands) Register(name string, f func(*config.State, Command) error) {
	c.Handlers[name] = f
}

func LoginHandler(s *config.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no argument was given")
	}

	if len(cmd.Args) > 1 {
		return fmt.Errorf("command only takes one argumeant: <command> [userName]")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	defer cancel()
	exists, err := s.Db.UserExists(ctx, cmd.Args[0])
	if err != nil {
		return fmt.Errorf("failed checking if user exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("User 【%s】 is not registered", cmd.Args[0])
	}
	// add an else clause that tells the user that they are already logged in

	s.StConfig.SetUser(cmd.Args[0])
	fmt.Printf("login with 【%s】 was successful\n", cmd.Args[0])

	return nil
}

func RegisterHandler(s *config.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no argument was given")
	}

	if len(cmd.Args) > 1 {
		return fmt.Errorf("command only takes one argumeant: <command> [userName]")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	defer cancel()
	exists, err := s.Db.UserExists(ctx, cmd.Args[0])
	if err != nil {
		return fmt.Errorf("failed checking if user exists: %w", err)
	}
	if exists {
		return fmt.Errorf("user 【%s】is already registered", cmd.Args[0])
	}

	createUserParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Name:      cmd.Args[0],
	}

	_, err = s.Db.CreateUser(ctx, createUserParams)
	if err != nil {
		return fmt.Errorf("failed creating a user: %w", err)
	}

	s.StConfig.SetUser(cmd.Args[0])
	fmt.Printf("user 【%s】 was successfully registered", cmd.Args[0])

	return nil
}

func ResetHandler(s *config.State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("command does not accept any arguments")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	defer cancel()
	if err := s.Db.DeleteAllUsers(ctx); err != nil {
		return fmt.Errorf("failed deleting all users: %w", err)
	}

	s.StConfig.SetUser("[None]")
	fmt.Println("All users have been deleted!")

	return nil
}

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

	for _, user := range users {
		if user.Name == s.StConfig.Current_user_name {
			fmt.Println("* " + s.StConfig.Current_user_name + " (current)")
		} else {
			fmt.Println("* " + user.Name)
		}
	}

	return nil
}

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

func AddFeedHandler(s *config.State, cmd Command) error {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	defer cancel()

	users, err := s.Db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("AddFeedHandler failed fetching users: %w", err)
	}

	if len(users) == 0 {
		return fmt.Errorf("no users in database to add a feed")
	}

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

	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    fetchUserId(users, s),
	}

	feed, err := s.Db.CreateFeed(ctx, newFeed)
	if err != nil {
		return fmt.Errorf("failed creating feed")
	}

	fmt.Println("successfully created feed")
	feedFollowParams := database.CreateFeedFollowParams{
		UserID:    feed.UserID,
		FeedID:    feed.ID,
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	fmt.Println("inserted feed follow record")
	feedFollow, err := s.Db.CreateFeedFollow(ctx, feedFollowParams)
	if err != nil {
		return fmt.Errorf("failed creating feed follow record: %w", err)
	}
	fmt.Printf("	* FeedName:      %s\n", feedFollow.FeedName)
	fmt.Printf("	* CurrentUser:   %s\n", s.StConfig.Current_user_name)

	fmt.Printf("successfuly added feed for user 【%s】\n", s.StConfig.Current_user_name)
	fmt.Printf("%v\n", feed.CreatedAt.Time)

	return nil
}

func fetchUserId(users []database.User, s *config.State) uuid.UUID {
	for _, user := range users {
		if user.Name == s.StConfig.Current_user_name {
			return user.ID
		}
	}
	return uuid.Nil
}

func isValidUrl(str string) bool {
	u, err := url.Parse(str)

	return err == nil && u.Scheme != "" && u.Host != ""
}

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

func FollowHandler(s *config.State, cmd Command) error {
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

func FeedFollowingHandler(s *config.State, cmd Command) error {
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
