package fakes

import (
	"errors"
	"fmt"
	"github.com/bxcodec/faker/v4"
	"go-mysql-test/helper"
	"log"
	"math"
	"math/rand"
	"reflect"
	"time"
)

func init() {
	helper.LoadEnv()
	db := helper.ConnectDb()

	rand.Seed(time.Now().UnixNano())

	var shopsCount int
	err := db.QueryRow("SELECT COUNT(*) as total FROM shops").Scan(&shopsCount)
	if err != nil {
		log.Fatal(err)
	}

	_ = faker.AddProvider("fakeShopId", func(v reflect.Value) (interface{}, error) {
		if shopsCount == 0 {
			return nil, errors.New("shops not found")
		}

		return 1 + rand.Intn(shopsCount-1+1), nil
	})

	_ = faker.AddProvider("fakePrice32", func(v reflect.Value) (interface{}, error) {
		return float32(100+rand.Intn(1000000-100+1)) / 100, nil
	})

	_ = faker.AddProvider("fakeOpensAt", func(v reflect.Value) (interface{}, error) {
		return fmt.Sprintf("%02d:%02d:00", 1+rand.Intn(12-1+1), 0+rand.Intn(5-0+1)*10), nil
	})

	_ = faker.AddProvider("fakeClosesAt", func(v reflect.Value) (interface{}, error) {
		return fmt.Sprintf("%02d:%02d:00", 13+rand.Intn(24-13+1), 0+rand.Intn(5-0+1)*10), nil
	})
}

type Faker interface{}

func GenerateBundle[F Faker](bundleSize int, creator func() F) []F {
	var bundle []F

	for i := 0; i < bundleSize; i++ {
		bundle = append(bundle, creator())
	}

	return bundle
}

func GenerateBundles[F Faker](count int, bundleSize int, creator func() F,
	callback func(bundle []F, bundleNum int, bundleSize int) error,
) error {
	bundleCount := int(math.Ceil(float64(count) / float64(bundleSize)))

	for bundleNum := 0; bundleNum < bundleCount; bundleNum++ {
		if bundleNum+1 == bundleCount { // If it`s last loop
			bundleSize = count - bundleNum*bundleSize
		}

		bundle := GenerateBundle(bundleSize, creator)

		if err := callback(bundle, bundleNum, bundleSize); err != nil {
			return err
		}
	}

	return nil
}
