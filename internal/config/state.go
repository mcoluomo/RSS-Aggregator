package config

import "github.com/mcoluomo/RSS-Aggregator/internal/database"

type State struct {
	Db       *database.Queries
	StConfig *Config
}
