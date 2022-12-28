package migration

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"go_async_shops_products/cli"
	"go_async_shops_products/helper"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const Usage = `migration OPTIONS COMMAND [arg...]  Manage Migrations
	migration [ -help ]`
var CommandsUsage = []string{CommandCreateUsage, CommandDeleteUsage, CommandUpUsage, CommandDownUsage, CommandForceUsage}

func initMigrater() *migrate.Migrate {
	db := helper.ConnectDb()

	driver, _ := mysql.WithInstance(db, &mysql.Config{})
	migrater, err := migrate.NewWithDatabaseInstance(
		"file://"+os.Getenv("MIGRATION_DIR"),
		os.Getenv("DB_DATABASE"),
		driver,
	)

	if err != nil {
		log.Fatal(err)
	}

	if err == nil {
		migrater.PrefetchMigrations = 10
		migrater.LockTimeout = time.Duration(15) * time.Second

		// handle Ctrl+c
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT)
		go func() {
			for range signals {
				log.Println("Stopping after this running migration ...")
				migrater.GracefulStop <- true
				return
			}
		}()
	}

	return migrater
}

func createCommand(args []string, migrater *migrate.Migrate) (cli.Commander, error) {
	switch args[0] {
	case "create":
		return NewCommandCreate(args[1:], migrater), nil

	case "delete":
		helper.ConfirmAction("Are you sure you want to DELETE migrations? [y/N]")

		return NewCommandDelete(args[1:], migrater), nil

	case "up":
		return NewCommandUp(args[1:], migrater), nil

	case "down":
		helper.ConfirmAction("Are you sure you want to DOWN migrations? [y/N]")

		return NewCommandDown(args[1:], migrater), nil

	case "force":
		return NewCommandForce(args[1:], migrater), nil

	default:
		return nil, errors.New("error createCommand: command does`t exists")
	}
}

func processArgs(args []string, migrater *migrate.Migrate) error {
	var cmd cli.Commander

	cmd, err := createCommand(args, migrater)

	if err == nil {
		err = cli.RunCommand(cmd)
	}

	if err != nil {
		return errors.New(fmt.Sprintf("error processArgs: " + err.Error()))
	}

	return nil
}

func Main(args []string) error {
	flagSet := helper.InitFlagSet(args, Usage, CommandsUsage)
	migrater := initMigrater()

	startTime := time.Now()
	err := processArgs(flagSet.Args(), migrater)
	log.Println("Finished after", time.Since(startTime))

	if err != nil {
		return errors.New(fmt.Sprintf("error cli/migration: " + err.Error()))
	}
	return nil
}