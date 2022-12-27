package helper

import (
	"database/sql"
	"flag"
	"fmt"
	driverMysql "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
	"time"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func ConnectDb() *sql.DB {
	cfg := driverMysql.Config{
		User:                 os.Getenv("DB_USERNAME"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		Addr:                 os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		DBName:               os.Getenv("DB_DATABASE"),
		AllowNativePasswords: true,
		InterpolateParams:    true,
		MultiStatements:      true,
	}

	db, err := sql.Open(os.Getenv("DB_DRIVER"), cfg.FormatDSN())
	if err != nil {
		log.Fatalf("Error opening DB: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging DB: %v", err)
	}

	db.SetConnMaxLifetime(time.Second * 0)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	return db
}

func ConfirmAction(msg string) {
	log.Println(msg)
	var response string
	fmt.Scanln(&response)
	response = strings.ToLower(strings.TrimSpace(response))

	if response == "y" {
		log.Println("Confirm")
	} else {
		log.Fatal("Cancelled")
	}
}

func CreateFile(filename string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)

	if err != nil {
		return err
	}

	return f.Close()
}

func GetFlagArgs() []string {
	flag.Parse()

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(2)
	}

	return flag.Args()[1:]
}
