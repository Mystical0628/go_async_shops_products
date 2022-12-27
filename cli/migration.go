package main

import (
	"go-mysql-test/cli/migration"
	"go-mysql-test/helper"
)

func init() { helper.LoadEnv() }

func main() {
	migration.Main()
}
