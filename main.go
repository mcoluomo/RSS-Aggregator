package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/mcoluomo/RSS-Aggregator/internal/cli"
	"github.com/mcoluomo/RSS-Aggregator/internal/config"
	"github.com/mcoluomo/RSS-Aggregator/internal/database"
)

func main() {
	fmt.Println("\nUsing rss-aggregator...")
	db, err := sql.Open("pgx", "postgres://postgres:olu@localhost:5432/gator?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("current user: %v\n", cfg.Current_user_name)

	cfgState := &config.State{StConfig: &cfg, Db: dbQueries}

	cmds := &cli.Commands{Handlers: map[string]func(*config.State, cli.Command) error{}}
	cmds.Register("login", cli.LoginHandler)
	cmds.Register("register", cli.RegisterHandler)
	cmds.Register("reset", cli.ResetHandler)
	cmds.Register("users", cli.UserHandler)
	cmds.Register("agg", cli.AggHandler)
	cmds.Register("addfeed", cli.MiddlewareLoggedIn(cli.AddFeedHandler))
	cmds.Register("feeds", cli.PrintFeedsHandler)
	cmds.Register("follow", cli.MiddlewareLoggedIn(cli.FollowHandler))
	cmds.Register("following", cli.MiddlewareLoggedIn(cli.FeedFollowingHandler))
	cmds.Register("unfollow", cli.MiddlewareLoggedIn(cli.UnfollowFeedFollow))
	cmds.Register("browse", cli.BrowseFeedsHandler)

	if len(os.Args) < 2 {
		log.Fatalf("\n---------------------------------\nPlease provide <command> [arg]\n---------------------------------\n")
	}

	cmd := cli.Command{Name: os.Args[1], Args: os.Args[2:]}
	if err = cmds.Run(cfgState, cmd); err != nil {
		log.Fatalf("\nerror running command: %v", err)
	}
	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	fmt.Printf("current user: %v\n", cfg.Current_user_name)
}
