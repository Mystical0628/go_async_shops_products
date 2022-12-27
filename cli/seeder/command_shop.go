package seeder

import (
	"bytes"
	"database/sql"
	"errors"
	"github.com/schollz/progressbar/v3"
	"go_async_shops_products/fakes"
	"strconv"
)

const CommandShopUsage = `shop [N]        Sow N shops`
const ShopBundleSize = 1000

type commandShop struct {
	*commandSeed
	n int
}

func (cmd *commandShop) Exec() error {
	err := cmd.commandSeed.Exec()

	if err == nil {
		bar := progressbar.Default(int64(cmd.n))

		err = fakes.GenerateBundlesShop(cmd.n,
			func(bundle []fakes.Shop, bundleNum int, bundleSize int) error {
				var valuesBuf bytes.Buffer

				valuesBuf.WriteString("INSERT INTO shops(name, url, opens_at, closes_at) VALUES ")

				for i := 0; i < bundleSize-1; i++ {
					valuesBuf.WriteString(bundle[i].InsertString() + ",\n")
				}

				valuesBuf.WriteString(bundle[bundleSize-1].InsertString())

				if _, err := cmd.db.Exec(valuesBuf.String()); err != nil {
					return err
				}

				return bar.Add(bundleSize)
			})
	}

	return err
}

func (cmd *commandShop) Validate() error {
	err := cmd.commandSeed.Validate()

	if err == nil && cmd.n < 1 {
		err = errors.New("limit argument N must be >= 1")
	}

	return err
}

func (cmd *commandShop) Parse() error {
	err := cmd.commandSeed.Parse()

	if err == nil {
		cmd.n, err = strconv.Atoi(cmd.GetFlagSet().Arg(0))

		if err != nil {
			err = errors.New("can't read limit argument N")
		}
	}

	return err
}

func NewCommandShop(args []string, db *sql.DB) *commandShop {
	return &commandShop{
		NewCommandSeed("product", args, CommandShopUsage, db),
		0,
	}
}
