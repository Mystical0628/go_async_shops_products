package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"go_async_shops_products/app"
	cliMigration "go_async_shops_products/cli/migration"
	cliSeeder "go_async_shops_products/cli/seeder"
	"go_async_shops_products/helper"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	helper.LoadEnv()
}

func initFlagUsage() {
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: OPTIONS COMMAND [arg...]
Options:
  -help            Print usage
Commands:
  %s
  %s
  %s
`, cliSeeder.Usage, cliMigration.Usage, app.Usage)
	}
}

func runCommand(args []string) error {
	switch args[0] {
	case "seeder":
		return cliSeeder.Main(args[1:])

	case "migration":
		return cliMigration.Main(args[1:])

	case "start":
		return app.Main(args[1:])

	default:
		return errors.New(fmt.Sprintf("error runCommand: command %s does`t exists", args[0]))
	}
}

func main() {
	initFlagUsage()
	flag.Parse()

	if len(flag.Args()) < 1 {
		flag.Usage()
		flag.PrintDefaults()
		os.Exit(2)
	}

	args := flag.Args()

	if err := runCommand(args); err != nil {
		log.Fatal(err)
	}
}