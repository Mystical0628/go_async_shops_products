package app

import (
	"database/sql"
	"errors"
	"fmt"
	"go_async_shops_products/helper"
	"log"
	"reflect"
	"time"
)

const Usage = "start                               Start application"

type app struct {
	db            *sql.DB
	time          time.Time
	timeFormatted string
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

	arguments := []reflect.Value{reflect.ValueOf(threads)}

	reflect.ValueOf(app).MethodByName("Action" + method).Call(arguments)

	log.Printf("%s Finished after: %v\n", method, time.Since(startTime))
}

func (app *app) Run() {
	app.CallAction("Index", 10)
}

func Main(args []string) error {
	app := NewApp()
	defer app.db.Close()
	app.Run()

	if err := app.db.Close(); err != nil {
		errors.New(fmt.Sprintf("error start: error while closing database %v", err))
	}

	return nil
}