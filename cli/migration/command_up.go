package migration

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"strconv"
)

const CommandUpUsage = `up [N] [-all]   Apply all or N up migrations`

type commandUp struct {
	*command
	flagAll *bool
	n       int
}

func (cmd *commandUp) Exec() error {
	err := cmd.command.Exec()

	if err == nil {
		if cmd.n >= 0 {
			err = cmd.migrater.Steps(cmd.n)
		} else {
			err = cmd.migrater.Up()
		}
	}

	return err
}

func (cmd *commandUp) Validate() error {
	err := cmd.command.Validate()

	if err == nil && cmd.n < -1 {
		err = errors.New("limit argument N must be >= -1")
	}

	return err
}

func (cmd *commandUp) Parse() error {
	err := cmd.command.Parse()

	if err == nil {
		if cmd.GetFlagSet().NArg() > 0 {
			cmd.n, err = strconv.Atoi(cmd.GetFlagSet().Arg(0))
			if err != nil {
				err = errors.New("can't read limit argument N")
			}
		} else {
			cmd.n = 1
		}

		if *cmd.flagAll {
			cmd.n = -1
		}
	}

	return err
}

func NewCommandUp(args []string, migrater *migrate.Migrate) *commandUp {
	cmd := NewCommand("up", args, CommandUpUsage, migrater)
	cmd.AllowNoArgs()

	return &commandUp{
		cmd,
		cmd.GetFlagSet().Bool("all", false, "Apply all up migrations"),
		1,
	}
}
