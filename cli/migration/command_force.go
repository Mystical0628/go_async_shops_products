package migration

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"strconv"
)

const CommandForceUsage = `force V         Set version V but don't run migration (ignores dirty state)`

type commandForce struct {
	*command
	v int
}

func (cmd *commandForce) Call() error {
	err := cmd.command.Exec()

	if err == nil {
		err = cmd.migrater.Force(cmd.v)
	}

	return err
}

func (cmd *commandForce) Validate() error {
	err := cmd.command.Validate()

	if err == nil && cmd.v < -1 {
		err = errors.New("argument V must be >= -1")
	}

	return err
}

func (cmd *commandForce) Parse() error {
	err := cmd.command.Parse()

	if err == nil {
		cmd.v, err = strconv.Atoi(cmd.GetFlagSet().Arg(0))
		if err != nil {
			err = errors.New("can't read version argument V")
		}
	}

	return nil
}

func NewCommandForce(args []string, migrater *migrate.Migrate) *commandForce {
	return &commandForce{
		command: NewCommand("force", args, CommandForceUsage, migrater),
		v:       -1,
	}
}
