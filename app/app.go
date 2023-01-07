package app

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"go_async_shops_products/helper"
	"log"
	"reflect"
	"time"
)

const Usage = `start                               Start application
	start [ -help ]`

type app struct {
	db            *sql.DB
	time          time.Time
	timeFormatted string
	flagSet       *flag.FlagSet
	flagHelp      *bool
	flagThreads   *int
	flagShops     *int
	flagProducts  *int
}

func NewApp() *app {
	db := helper.ConnectDb()

	time := time.Now()
	// timeFormatted := time.Format("15:04:05")
	timeFormatted := "12:00:00"

	return &app{
		db:            db,
		time:          time,
		timeFormatted: timeFormatted,
	}
}

func (app *app) CallAction(method string, threads int) {
	startTime := time.Now()

	// arguments := []reflect.Value{reflect.ValueOf(threads)}

	reflect.ValueOf(app).MethodByName("Action" + method).Call(nil)

	log.Printf("%s Finished after: %v\n", method, time.Since(startTime))
}

func (app *app) Run() {
	app.CallAction("Index", 10)
	app.CallAction("All", 10)
}

func Main(args []string) error {
	app := NewApp()
	defer app.db.Close()
	app.flagSet = helper.InitFlagSetCallback(args, Usage, nil, func(flagSet *flag.FlagSet) {
		app.flagHelp = flagSet.Bool("help", false, "Print help information")
		app.flagThreads = flagSet.Int("threads", 10, "Count od threads")
		app.flagShops = flagSet.Int("shops", 0, "Count of shops to process")
		app.flagProducts = flagSet.Int("products", 0, "Count of products to process")
	})
	app.Run()

	if err := app.db.Close(); err != nil {
		errors.New(fmt.Sprintf("error start: error while closing database %v", err))
	}

	return nil
}
