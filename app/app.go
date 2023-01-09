package app

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"go_async_shops_products/helper"
	"log"
	"reflect"
	"strings"
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
	flagAll       *bool
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

func (app *app) callAction(action string, threads int) error {
	m, ok := reflect.TypeOf(app).MethodByName("Action" + action)

	if !ok {
		return errors.New("action doesn't exists")
	}

	arguments := []reflect.Value{reflect.ValueOf(app)}
	m.Func.Call(arguments)

	return nil
}

func (app *app) callAllActions() {
	t := reflect.TypeOf(app)
	arguments := []reflect.Value{reflect.ValueOf(app)}

	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)

		if strings.Contains(m.Name, "Action") {
			startTime := time.Now()
			log.Printf("Start Action '%v'\n", m.Name)
			m.Func.Call(arguments)
			log.Printf("Action '%v' Finished after: %v\n\n", m.Name, time.Since(startTime))
		}
	}
}

func (app *app) Run() error {
	startTime := time.Now()
	defer log.Printf("Action Finished after: %v\n", time.Since(startTime))

	action := app.flagSet.Arg(0)

	if *app.flagAll {
		app.callAllActions()
	} else {
		return app.callAction(action, 10)
	}

	return nil
}

func Main(args []string) error {
	app := NewApp()
	defer app.db.Close()
	app.flagSet = helper.InitFlagSetCallback(args, Usage, nil, func(flagSet *flag.FlagSet, flagHelp *bool, allowNoArgs *bool) {
		*allowNoArgs = true

		app.flagHelp = flagHelp
		app.flagAll = flagSet.Bool("all", false, "Start all actions")
		app.flagThreads = flagSet.Int("threads", 10, "Count od threads")
		app.flagShops = flagSet.Int("shops", 0, "Count of shops to process")
		app.flagProducts = flagSet.Int("products", 0, "Count of products to process")
	})

	if err := app.Run(); err != nil {
		return errors.New(fmt.Sprintf("app: main: %v", err))
	}

	if err := app.db.Close(); err != nil {
		return errors.New(fmt.Sprintf("app: main: closing database %v", err))
	}

	return nil
}
