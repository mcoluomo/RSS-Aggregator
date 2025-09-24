package cli

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mcoluomo/RSS-Aggregator/internal/config"
	"github.com/mcoluomo/RSS-Aggregator/internal/database"
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
