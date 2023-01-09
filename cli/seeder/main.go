package seeder

import (
	"database/sql"
	"errors"
	"fmt"
	"go_async_shops_products/cli"
	"go_async_shops_products/helper"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const Usage = `seeder OPTIONS SEED [arg...]        Sow database
	seeder [ -help ]`

var CommandsUsage = []string{CommandProductUsage, CommandShopUsage, CommandTruncateUsage}

func createCommand(args []string, db *sql.DB) (cli.Commander, error) {
	log.Println(args[0])
	switch args[0] {
	case "product":
		return NewCommandProduct(args[1:], db), nil

	case "shop":
		return NewCommandShop(args[1:], db), nil

	case "truncate":
		helper.ConfirmAction("Are you sure you want to truncate TABLE? [y/N]")
		return NewCommandTruncate(args[1:], db), nil

	default:
		return nil, errors.New("error createCommand: command does`t exists")
	}
}

func processArgs(args []string, db *sql.DB) error {
	var cmd cli.Commander

	cmd, err := createCommand(args, db)

	if err == nil {
		err = cli.RunCommand(cmd)
	}

	if err != nil {
		return errors.New(fmt.Sprintf("error processArgs: " + err.Error()))
	}

	return err
}

func Main(args []string) error {
	flagSet := helper.InitFlagSet(args, Usage, CommandsUsage)
	db := helper.ConnectDb()

	startTime := time.Now()
	err := processArgs(flagSet.Args(), db)
	log.Println("Finished after", time.Since(startTime))

	if err != nil {
		return errors.New(fmt.Sprintf("error cli/seeder: " + err.Error()))
	}

	return nil
}
