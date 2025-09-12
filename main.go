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
	fmt.Printf("Reading config with the contents: %+v\n", cfg)

	cfgState := &config.State{StConfig: &cfg, Db: dbQueries}

	cmds := &cli.Commands{Handlers: map[string]func(*config.State, cli.Command) error{}}
	cmds.Register("login", cli.LoginHandler)
	cmds.Register("register", cli.RegisterHandler)
	cmds.Register("reset", cli.ResetHandler)
	cmds.Register("users", cli.UserHandler)
	cmds.Register("agg", cli.AggHandler)
	cmds.Register("addfeed", cli.AddFeedHandler)
	cmds.Register("feeds", cli.PrintFeedsHandler)

	if len(os.Args) < 2 {
		log.Fatalf("Please provide <command> [arg]")
	}

	cmd := cli.Command{Name: os.Args[1], Args: os.Args[2:]}
	if err := cmds.Run(cfgState, cmd); err != nil {
		log.Fatalf("\nerror running command: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("\nReading config again: %+v\n", cfg)
}
