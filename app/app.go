package app

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"go_async_shops_products/helper"
	"log"
	"os"
	"reflect"
	"time"
)

const Usage = `start                               Start application
	start [ -help ]`

type app struct {
	db            *sql.DB
	time          time.Time
	timeFormatted string
	flagSet 	*flag.FlagSet
	flagHelp 	*bool
	flagThreads 	*int
	flagShops 		*int
	flagProducts 	*int
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
}

func (app *app) initFlagSet(args []string) {
	app.flagSet = flag.NewFlagSet("start", flag.ExitOnError)

	app.flagHelp = app.flagSet.Bool("help", false, "Print help information")
	app.flagThreads = app.flagSet.Int("threads", 10, "Count od threads")
	app.flagShops = app.flagSet.Int("shops", 0, "Count of shops to process")
	app.flagProducts = app.flagSet.Int("products", 0, "Count of products to process")

	app.flagSet.Usage = func() {
		fmt.Println(Usage)
		//fmt.Println("Commands:")
		//for _, usage := range commandsUsage {
		//	fmt.Println("  " + usage)
		//}
		fmt.Println("Options:")
		app.flagSet.PrintDefaults()
	}

	if err := app.flagSet.Parse(args); err != nil {
		log.Fatal(err)
	}

	if *app.flagHelp {
		app.flagSet.Usage()
		os.Exit(2)
	}
}

func Main(args []string) error {
	app := NewApp()
	defer app.db.Close()
	app.initFlagSet(args)
	app.Run()

	if err := app.db.Close(); err != nil {
		errors.New(fmt.Sprintf("error start: error while closing database %v", err))
	}

	return nil
}