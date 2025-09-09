package cli

import (
	"context"
	"database/sql"
	"fmt"
	"os"
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
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	defer cancel()
	exists, err := s.Db.UserExists(ctx, cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error checking if user exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("User %s is not registered", cmd.Args[0])
	}
	// add an else clause that tells the user that they are already logged in

	s.StConfig.SetUser(os.Args[2])
	s.StConfig.Current_user_name = cmd.Args[0]
	fmt.Printf("⇒ %s\nlogin with %s was successful\n", cmd.Args[0], cmd.Args[0])

	return nil
}

func RegisterHandler(s *config.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no argument was given")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	defer cancel()
	exists, err := s.Db.UserExists(ctx, cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error checking if user exists: %w", err)
	}
	if exists {
		return fmt.Errorf("user %s is already registered", cmd.Args[0])
	}

	createUserParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Name:      cmd.Args[0],
	}

	_, err = s.Db.CreateUser(ctx, createUserParams)
	if err != nil {
		return fmt.Errorf("error encountered creating a user: %w", err)
	}

	s.StConfig.SetUser(os.Args[2]) // migth remove this
	s.StConfig.Current_user_name = cmd.Args[0]
	fmt.Printf("user %s was successfully registered", cmd.Args[0])

	return nil
}

func ResetHandler(s *config.State, cmd Command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("command does not accept any arguments")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

	defer cancel()
	if err := s.Db.DeleteAll(ctx); err != nil {
		return fmt.Errorf("error encountered deleting all users: %w", err)
	}
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
		return fmt.Errorf("error encountered fetching all users: %w", err)
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
