package migration

import (
	"github.com/golang-migrate/migrate/v4"
	"go_async_shops_products/cli"
)

//type Commander interface {
//	Exec() error
//	Parse() error
//	Validate() error
//}

type command struct {
	cli.Commander
	migrater *migrate.Migrate
}

func NewCommand(commandName string, args []string, usage string, migrater *migrate.Migrate) *command {
	cmd := cli.NewCommand(commandName, args, usage)
	return &command{
		Commander: cmd,
		migrater:  migrater,
	}
}
