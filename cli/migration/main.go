package migration

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"go-mysql-test/cli/cli"
	"go-mysql-test/helper"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func initMigrater() *migrate.Migrate {
	db := helper.ConnectDb()

	driver, _ := mysql.WithInstance(db, &mysql.Config{})
	migrater, err := migrate.NewWithDatabaseInstance(
		"file://"+os.Getenv("MIGRATION_DIR"),
		os.Getenv("DB_DATABASE"),
		driver,
	)

	defer func() {
		if err == nil {
			// TODO Solve error
			//if _, closeErr := migrater.Close(); err != nil {
			//	log.Println(closeErr)
			//}
		}
	}()

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

func initFlagUsage() {
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: migration OPTIONS COMMAND [arg...]
       migration [ -help ]
Options:
  -help            Print usage
Commands:
  %s
  %s
  %s
  %s
  %s
`, CommandCreateUsage, CommandDeleteUsage, CommandUpUsage, CommandDownUsage, CommandForceUsage)
	}
}

func parseFlagArgs() []string {
	flag.Parse()

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(2)
	}

	return flag.Args()[1:]
}

func createCommand(cmdName string, args []string, migrater *migrate.Migrate) (cli.Commander, error) {
	switch cmdName {
	case "create":
		return NewCommandCreate(args, migrater), nil

	case "delete":
		helper.ConfirmAction("Are you sure you want to DELETE migrations? [y/N]")

		return NewCommandDelete(args, migrater), nil

	case "up":
		return NewCommandUp(args, migrater), nil

	case "down":
		helper.ConfirmAction("Are you sure you want to DOWN migrations? [y/N]")

		return NewCommandDown(args, migrater), nil

	case "force":
		return NewCommandForce(args, migrater), nil

	default:
		return nil, errors.New("command does`t exists")
	}
}

func processArgs(args []string, migrater *migrate.Migrate) {
	var cmd cli.Commander

	cmd, err := createCommand(flag.Arg(0), args, migrater)

	if err == nil {
		err = cli.RunCommand(cmd)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func Main() {
	initFlagUsage()

	migrater := initMigrater()
	args := helper.GetFlagArgs()

	startTime := time.Now()

	processArgs(args, migrater)

	log.Println("Finished after", time.Since(startTime))
}
