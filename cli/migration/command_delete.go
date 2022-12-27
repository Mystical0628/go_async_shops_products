package migration

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

const CommandDeleteUsage = `delete V        Delete migration version V`

type commandDelete struct {
	*command
	v int
}

func (cmd *commandDelete) Call() error {
	err := cmd.command.Exec()

	var seqDigits int
	if err == nil {
		seqDigits, err = strconv.Atoi(os.Getenv("MIGRATION_SEQ_DIGITS"))
	}

	var seq bool
	if err == nil {
		seq, err = strconv.ParseBool(os.Getenv("MIGRATION_SEQ"))
	}

	var version string
	var matches []string
	dir := os.Getenv("MIGRATION_DIR")
	ext := "." + os.Getenv("MIGRATION_EXT")
	if err == nil {
		if seq {
			version = fmt.Sprintf("%0[2]*[1]d", cmd.v, seqDigits)
		} else {
			version = strconv.Itoa(cmd.v)
		}

		matches, err = filepath.Glob(filepath.Join(dir, version+"_*"+ext))
		if len(matches) == 0 {
			err = errors.New("not found migration version: " + version)
		}
	}

	if err == nil {
		for _, filename := range matches {
			if err = os.Remove(filename); err == nil {
				log.Println("Delete: " + filename)
			}
		}
	}

	return err
}

func (cmd *commandDelete) Validate() error {
	err := cmd.command.Validate()

	if err == nil && cmd.v < -1 {
		err = errors.New("argument V must be >= -1")
	}

	return err
}

func (cmd *commandDelete) Parse() error {
	err := cmd.command.Parse()

	if err == nil {
		cmd.v, err = strconv.Atoi(cmd.GetFlagSet().Arg(0))
		if err != nil {
			err = errors.New("can't read version argument V")
		}
	}

	return err
}

func NewCommandDelete(args []string, migrater *migrate.Migrate) *commandDelete {
	return &commandDelete{
		command: NewCommand("delete", args, CommandDeleteUsage, migrater),
		v:       -1,
	}
}
