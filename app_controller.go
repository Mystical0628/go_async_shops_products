package main

import (
	"github.com/schollz/progressbar/v3"
	"log"
	"os"
	"sync"
)

func (app *app) ActionSimple(threads int) {
	file, err := os.OpenFile("products_simple.xml", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("Error while opening file: %v", err)
	}

	shops := app.getShops()

	for _, shop := range shops {
		shop.Render(file)
	}

	productsItems := app.getProducts(0, app.productsTotal)
	products := &Products{Products: productsItems}
	products.Render(file)

	app.bar.Add(app.productsTotal)
}

func (app *app) ActionSimpleAsync(threads int) {
	file, err := os.OpenFile("products_simple_async.xml", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("Error while opening file: %v", err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(threads)

	productsRows := app.getProductsRows(0, app.productsTotal)
	defer productsRows.Close()

	productsRowsChan := make(chan *Product, threads)

	for i := 0; i < threads; i++ {
		go runner(productsRowsChan, wg, app.bar, file, i)
	}

	for productsRows.Next() {
		item := &Product{}
		if err := productsRows.Scan(&item.Id, &item.ShopId, &item.Name, &item.Description, &item.Price); err != nil {
			log.Fatalf("Error while scanning row: %v", err)
		}

		productsRowsChan <- item
	}

	close(productsRowsChan)

	wg.Wait()

	err = productsRows.Close()
	if err != nil {
		log.Fatal(err)
	}
}

//func (app *app) ActionBundles(bundleSize int) {
//	file, err := os.OpenFile("products_bundle.xml", os.O_CREATE|os.O_APPEND|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
//	if err != nil {
//		log.Fatalf("Error while opening file: %v", err)
//	}
//
//	file.Close()
//
//	bundleCount := int(math.Ceil(float64(app.productsTotal) / float64(bundleSize)))
//
//	for bundleNum := 0; bundleNum < bundleCount; bundleNum++ {
//		if bundleNum+1 == bundleCount { // If it`s last loop
//			bundleSize = app.productsTotal - bundleNum*bundleSize
//		}
//
//		products := app.getProducts(bundleNum*bundleSize, bundleSize)
//		for _, item := range products {
//			item.Render(file)
//		}
//		app.bar.Add(bundleSize)
//	}
//}

func runner(productsChan chan *Product, wg *sync.WaitGroup, bar *progressbar.ProgressBar, file *os.File, n int) {
	defer wg.Done()

	for item := range productsChan {
		item.Render(file)
		bar.Add(1)
	}
}

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
