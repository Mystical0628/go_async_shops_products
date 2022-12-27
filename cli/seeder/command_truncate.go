package seeder

import (
	"database/sql"
)

const CommandTruncateUsage = `truncate TABLE  Truncate the table TABPLE`

type commandTruncate struct {
	*commandSeed
	table string
}

func (cmd *commandTruncate) Exec() error {
	err := cmd.commandSeed.Exec()

	if err == nil {
		_, err = cmd.db.Exec(`
			SET FOREIGN_KEY_CHECKS = 0;
			TRUNCATE TABLE ` + cmd.table + `;
			SET FOREIGN_KEY_CHECKS = 1;
		`)
	}

	return err
}

func (cmd *commandTruncate) Validate() error {
	err := cmd.commandSeed.Validate()

	return err
}

func (cmd *commandTruncate) Parse() error {
	err := cmd.commandSeed.Parse()

	if err == nil {
		cmd.table = cmd.GetFlagSet().Arg(0)
	}

	return err
}

func NewCommandTruncate(args []string, db *sql.DB) *commandTruncate {
	return &commandTruncate{
		NewCommandSeed("product", args, CommandTruncateUsage, db),
		"",
	}
}
