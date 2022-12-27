package migration

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"go_async_shops_products/helper"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const CommandCreateUsage = `create NAME     Create a set of timestamped up/down migrations titled NAME.`

type commandCreate struct {
	*command
	name string
}

func (cmd *commandCreate) Call() error {
	if err := cmd.command.Exec(); err != nil {
		return err
	}

	var version string
	var err error

	dir := os.Getenv("MIGRATION_DIR")
	ext := "." + os.Getenv("MIGRATION_EXT")

	seqDigits, err := strconv.Atoi(os.Getenv("MIGRATION_SEQ_DIGITS"))
	if err != nil {
		return err
	}

	seq, err := strconv.ParseBool(os.Getenv("MIGRATION_SEQ"))
	if err != nil {
		return err
	}

	if seq {
		fileMatches, err := filepath.Glob(filepath.Join(dir, "*"+ext))
		if err != nil {
			return err
		}

		version, err = nextSeqVersion(fileMatches, seqDigits)
		if err != nil {
			return err
		}
	} else {
		version = strconv.FormatInt(time.Now().Unix(), 10)
	}

	matches, err := filepath.Glob(filepath.Join(dir, version+"_*"+ext))
	if err != nil {
		return err
	}

	if len(matches) > 0 {
		return fmt.Errorf("duplicate migration version: %s", version)
	}

	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	for _, direction := range []string{"up", "down"} {
		basename := fmt.Sprintf("%s_%s.%s%s", version, cmd.name, direction, ext)
		filename := filepath.Join(dir, basename)

		if err = helper.CreateFile(filename); err != nil {
			return err
		}

		log.Println("Create: " + filename)
	}

	return nil
}

func (cmd *commandCreate) Validate() error {
	err := cmd.command.Validate()

	return err
}

func (cmd *commandCreate) Parse() error {
	err := cmd.command.Parse()

	cmd.name = cmd.GetFlagSet().Arg(0)

	return err
}

func NewCommandCreate(args []string, migrater *migrate.Migrate) *commandCreate {
	return &commandCreate{
		command: NewCommand("create", args, CommandCreateUsage, migrater),
		name:    "",
	}
}

func nextSeqVersion(fileMatches []string, seqDigits int) (string, error) {
	if seqDigits <= 0 {
		return "", errors.New("Digits must be positive")
	}

	nextSeq := uint64(1)

	if len(fileMatches) > 0 {
		filename := fileMatches[len(fileMatches)-1]
		matchSeqStr := filepath.Base(filename)
		idx := strings.Index(matchSeqStr, "_")

		if idx < 1 { // Using 1 instead of 0 since there should be at least 1 digit
			return "", fmt.Errorf("Malformed migration filename: %s", filename)
		}

		var err error
		matchSeqStr = matchSeqStr[0:idx]
		nextSeq, err = strconv.ParseUint(matchSeqStr, 10, 64)

		if err != nil {
			return "", err
		}

		nextSeq++
	}

	version := fmt.Sprintf("%0[2]*[1]d", nextSeq, seqDigits)

	if len(version) > seqDigits {
		return "", fmt.Errorf("Next sequence number %s too large. At most %d digits are allowed", version, seqDigits)
	}

	return version, nil
}
