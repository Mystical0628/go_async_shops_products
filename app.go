package main

import (
	"database/sql"
	"github.com/schollz/progressbar/v3"
	"go-mysql-test/helper"
	"log"
	"reflect"
	"sync"
	"time"
)

type app struct {
	db            *sql.DB
	wg            sync.WaitGroup
	time          time.Time
	timeFormatted string
	viewCreated   bool
	bar           *progressbar.ProgressBar
	productsTotal int
}

func NewApp() *app {
	db := helper.ConnectDb()

	db.SetConnMaxLifetime(time.Second * 1)
	db.SetMaxOpenConns(30)
	db.SetMaxIdleConns(10)

	time := time.Now()
	// timeFormatted := time.Format("15:04:05")
	timeFormatted := "12:00:00"

	return &app{
		db:            db,
		wg:            sync.WaitGroup{},
		time:          time,
		timeFormatted: timeFormatted,
		viewCreated:   true,
	}
}

func (app *app) RunMethod(method string, bundleSize int) {
	startTime := time.Now()
	app.bar = progressbar.Default(int64(app.productsTotal))

	arguments := []reflect.Value{reflect.ValueOf(bundleSize)}

	reflect.ValueOf(app).MethodByName("Action" + method).Call(arguments)

	log.Printf("%s Finished after: %v\n", method, time.Since(startTime))
}

func (app *app) Run() {
	// app.productsTotal := app.getProductsTotal()

	app.productsTotal = 1000000

	app.RunMethod("Simple", 0)
	app.RunMethod("Bundles", 50000)
	// app.RunMethod("BundlesAsync", 50000)

	app.wg.Wait()
}
