package seeder

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"go-mysql-test/cli/cli"
	"go-mysql-test/helper"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func initFlagUsage() {
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: seeder OPTIONS SEED [arg...]
       seeder [ -help ]
Options:
  -help            Print usage
Seeds:
  %s
  %s
  %s
`, CommandProductUsage, CommandShopUsage, CommandTruncateUsage)
	}
}

func createCommand(seedName string, args []string, db *sql.DB) (cli.Commander, error) {
	switch seedName {
	case "product":
		return NewCommandProduct(args, db), nil

	case "shop":
		return NewCommandShop(args, db), nil

	case "truncate":
		helper.ConfirmAction("Are you sure you want to truncate TABLE? [y/N]")

		return NewCommandTruncate(args, db), nil

	default:
		return nil, errors.New("command does`t exists")
	}
}

func processArgs(args []string, db *sql.DB) {
	var cmd cli.Commander

	cmd, err := createCommand(flag.Arg(0), args, db)

	if err == nil {
		err = cli.RunCommand(cmd)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func Main() {
	initFlagUsage()

	db := helper.ConnectDb()
	args := helper.GetFlagArgs()

	startTime := time.Now()

	processArgs(args, db)

	log.Println("Finished after", time.Since(startTime))
}
