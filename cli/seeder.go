package main

import (
	"go-mysql-test/cli/seeder"
	"go-mysql-test/helper"
)

func init() { helper.LoadEnv() }

func main() {
	seeder.Main()
}
