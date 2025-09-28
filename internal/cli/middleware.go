package cli

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/mcoluomo/RSS-Aggregator/internal/config"
	"github.com/mcoluomo/RSS-Aggregator/internal/database"
)

func MiddlewareLoggedIn(handler func(s *config.State, cmd Command, user database.User) error) func(*config.State, Command) error {
	return func(s *config.State, cmd Command) error {
		ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)

		defer cancel()

		if s.StConfig.Current_user_name == "[None]" {
			return fmt.Errorf("no users in registered in database")
		}

		user, err := s.Db.GetUser(ctx, s.StConfig.Current_user_name)
		if err != nil {
			return fmt.Errorf("%w: failed fetching user: 【%s】", err, s.StConfig.Current_user_name)
		}

		if err = handler(s, cmd, user); err != nil {
			return fmt.Errorf("failed calling handler")
		}

		return nil
	}
}

func isValidUrl(str string) bool {
	u, err := url.Parse(str)

	return err == nil && u.Scheme != "" && u.Host != ""
}
