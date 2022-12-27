package cli

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

type Commander interface {
	GetCommandName() string
	GetFlagSet() *flag.FlagSet
	IsAllowNoArgs() bool
	AllowNoArgs()
	DisallowNoArgs()
	Exec() error
	Parse() error
	Validate() error
}

type command struct {
	commandName string
	allowNoArgs bool
	flagSet     *flag.FlagSet
	flagHelp    *bool
	args        []string
	usage       string
}

func (cmd *command) GetCommandName() string {
	return cmd.commandName
}

func (cmd *command) IsAllowNoArgs() bool {
	return cmd.allowNoArgs
}

func (cmd *command) AllowNoArgs() {
	cmd.allowNoArgs = true
}

func (cmd *command) DisallowNoArgs() {
	cmd.allowNoArgs = false
}

func (cmd *command) GetFlagSet() *flag.FlagSet {
	return cmd.flagSet
}

func (cmd *command) GetFlagHelp() *bool {
	return cmd.flagHelp
}

func (cmd *command) GetArgs() []string {
	return cmd.args
}

func (cmd *command) GetUsage() string {
	return cmd.usage
}

func (cmd *command) Exec() error {
	return nil
}

func (cmd *command) Parse() error {
	if err := cmd.flagSet.Parse(cmd.args); err != nil {
		log.Fatal(err)
	}

	if *cmd.flagHelp {
		fmt.Fprintln(os.Stderr, cmd.usage)
		cmd.flagSet.PrintDefaults()
		os.Exit(0)
	}

	if !cmd.allowNoArgs && cmd.flagSet.NArg() == 0 {
		return errors.New("please specify arguments")
	}

	return nil
}

func (cmd *command) Validate() error {
	return nil
}

func NewCommand(cmdName string, args []string, usage string) *command {
	flagSet := flag.NewFlagSet(cmdName, flag.ExitOnError)
	flagHelp := flagSet.Bool("help", false, "Print help information")

	return &command{
		commandName: cmdName,
		allowNoArgs: false,
		flagSet:     flagSet,
		flagHelp:    flagHelp,
		args:        args,
		usage:       usage,
	}
}

func RunCommand(cmd Commander) error {
	err := cmd.Parse()

	if err == nil {
		err = cmd.Validate()
	}

	if err == nil {
		err = cmd.Exec()
	}

	return err
}
