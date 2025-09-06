package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mcoluomo/RSS-Aggregator/internal/cli"
	"github.com/mcoluomo/RSS-Aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Reading config with the contents: %+v\n", cfg)

	cfgState := &config.State{StConfig: &cfg}

	cmds := &cli.Commands{Handlers: map[string]func(*config.State, cli.Command) error{}}
	cmds.Register("login", cli.LoginHandler)

	if len(os.Args) < 2 {
		log.Fatalf("\nUsing rss-aggregator...\nPlease provide <command> [args]")
	}

	cmd := cli.Command{Name: os.Args[1], Args: os.Args[2:]}
	if err := cmds.Run(cfgState, cmd); err != nil {
		log.Fatalf("\nerror running command: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Reading config again: %+v\n", cfg)
}
