package cli

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/mcoluomo/RSS-Aggregator/internal/cli"
	"github.com/mcoluomo/RSS-Aggregator/internal/config"
	"github.com/mcoluomo/RSS-Aggregator/internal/database"
)

func MiddlewareLoggedIn(handler func(s *config.State, cmd Command, user database.User) error) func(*config.State, cli.Command) error {
	return func(s *config.State, cmd cli.Command) error {
		ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

		defer cancel()

		users, err := s.Db.GetUsers(ctx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// change this message tomorrow
				return fmt.Errorf("failed getting users: %s %w", s.StConfig.Current_user_name, err)
			}
		}

		// see if we still need to use this
		if len(users) == 0 || s.StConfig.Current_user_name == "[None]" {
			return fmt.Errorf("no users in database to create a feed follow record")
		}
		return nil
	}
}

func isValidUrl(str string) bool {
	u, err := url.Parse(str)

	return err == nil && u.Scheme != "" && u.Host != ""
}
