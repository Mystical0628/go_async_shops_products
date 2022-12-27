package seeder

import (
	"database/sql"
	"go-mysql-test/cli/cli"
)

type commandSeed struct {
	cli.Commander
	db *sql.DB
}

func NewCommandSeed(commandSeedName string, args []string, usage string, db *sql.DB) *commandSeed {
	return &commandSeed{
		cli.NewCommand(commandSeedName, args, usage),
		db,
	}
}
