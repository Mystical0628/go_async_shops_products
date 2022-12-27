package main

import (
	"log"
	"math"
	"os"
)

//func (app *app) ActionSimple(bundleSize int) {
//	file, err := os.OpenFile("products_simple.xml", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
//	if err != nil {
//		log.Fatalf("Error while opening file: %v", err)
//	}
//
//	productsItems := app.getProducts(0, app.productsTotal)
//	products := &Products{Products: productsItems}
//	products.Render(file)
//
//	app.bar.Add(app.productsTotal)
//}

func (app *app) ActionBundles(bundleSize int) {
	file, err := os.OpenFile("products_bundle.xml", os.O_CREATE|os.O_APPEND|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("Error while opening file: %v", err)
	}

	file.Close()

	bundleCount := int(math.Ceil(float64(app.productsTotal) / float64(bundleSize)))

	for bundleNum := 0; bundleNum < bundleCount; bundleNum++ {
		if bundleNum+1 == bundleCount { // If it`s last loop
			bundleSize = app.productsTotal - bundleNum*bundleSize
		}

		products := app.getProducts(bundleNum*bundleSize, bundleSize)
		for _, item := range products {
			item.Render(file)
		}
		app.bar.Add(bundleSize)
	}
}

//func runner(productsChan chan []Product) {
//	for _, item := range products {
//		item.Render(file)
//	}
//}

//func (app *app) ActionBundlesAsync(bundleSize int) {
//	productsTotal := app.getProductsTotal()
//
//	var products [][5]string
//	bundleCount := int(math.Ceil(float64(productsTotal) / float64(bundleSize)))
//
//	productsChan := make(chan [][5]string)
//
//	for bundleNum := 0; bundleNum < bundleCount; bundleNum++ {
//		if bundleNum+1 == bundleCount { // If it`s last loop
//			bundleSize = productsTotal - bundleNum*bundleSize
//		}
//
//		app.wg.Add(1)
//		go func() {
//			productsChan <- app.getProducts(bundleNum*bundleSize, bundleSize)
//			app.bar.Add(bundleSize)
//		}()
//
//		go func() {
//			defer app.wg.Done()
//			products = append(products, <-productsChan...)
//		}()
//
//		// products = append(products, <-productsChan...)
//	}
//
//	app.wg.Wait()
//
//	log.Println(len(products))
//}
