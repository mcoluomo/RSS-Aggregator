package main

import (
	"fmt"
	"log"

	"github.com/RSS-Aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	cfg.SetUser("olu")

	cfg, err = config.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)
}
