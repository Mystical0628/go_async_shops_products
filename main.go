package main

import (
	_ "github.com/go-sql-driver/mysql"
	"go-mysql-test/helper"
	"log"
)

func init() {
	helper.LoadEnv()
}

func main() {
	app := NewApp()
	defer app.db.Close()
	app.Run()

	if err := app.db.Close(); err != nil {
		log.Fatalf("Error while closing database %v", err)
	}
}
