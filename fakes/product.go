package fakes

import (
	"fmt"
	"github.com/bxcodec/faker/v4"
)

const ProductBundleSize = 10000

type Product struct {
	ShopId      int     `faker:"fakeShopId"`
	Price       float32 `faker:"fakePrice32"`
	Name        string  `faker:"name"`
	Description string  `faker:"html"`
}

func (p Product) InsertString() string {
	return fmt.Sprintf("(%v, \"%s\", \"%s\", %v)", p.ShopId, p.Name, p.Description, p.Price)
}

func NewProduct() Product {
	product := Product{}
	faker.FakeData(&product)
	return product
}

func GenerateBundlesProduct(count int, callback func(bundle []Product, bundleNum int, bundleSize int) error) error {
	return GenerateBundles[Product](count, ProductBundleSize, NewProduct, callback)
}
