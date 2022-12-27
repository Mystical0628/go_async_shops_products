package fakes

import (
	"fmt"
	"github.com/bxcodec/faker/v4"
)

const ShopBundleSize = 10000

type Shop struct {
	Name     string `faker:"name"`
	Url      string `faker:"url"`
	OpensAt  string `faker:"fakeOpensAt"`
	ClosesAt string `faker:"fakeClosesAt"`
}

func (s Shop) InsertString() string {
	return fmt.Sprintf("(\"%s\", \"%s\", \"%s\", \"%s\")", s.Name, s.Url, s.OpensAt, s.ClosesAt)
}

func NewShop() Shop {
	shop := Shop{}
	faker.FakeData(&shop)
	return shop
}

func GenerateBundlesShop(count int, callback func(bundle []Shop, bundleNum int, bundleSize int) error) error {
	return GenerateBundles[Shop](count, ShopBundleSize, NewShop, callback)
}
