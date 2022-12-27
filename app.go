package main

import (
	"database/sql"
	"go-mysql-test/helper"
	"log"
	"reflect"
	"time"
)

type app struct {
	db            *sql.DB
	time          time.Time
	timeFormatted string
	viewCreated   bool
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
		viewCreated:   true,
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
