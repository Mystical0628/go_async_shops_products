package main

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"log"
	"os"
	"sync"
)

func (app *app) ActionIndex(threads int) {
	file, err := os.OpenFile("shops_products.xml", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("Error while opening file: %v", err)
	}

	shops := app.getShops()

	for _, shop := range shops {
		app.processShopProducts(shop, threads)
		log.Println(shop)
		shop.Render(file)
	}
}

func (app *app) processShopProducts(shop *Shop, threads int) {
	wg := &sync.WaitGroup{}
	productsTotal := app.getShopProductsTotal(shop.Id)
	bar := progressbar.Default(int64(productsTotal), fmt.Sprintf("Shop %d", shop.Id))
	productsRows := app.getShopProductsRows(shop.Id)

	defer productsRows.Close()

	wg.Add(threads)

	productsRowsChan := make(chan *Product, threads)

	for i := 0; i < threads; i++ {
		go processProduct(productsRowsChan, wg, shop.Products, bar)
	}

	for productsRows.Next() {
		item := &Product{}
		if err := productsRows.Scan(&item.Id, &item.Name, &item.Description, &item.Price); err != nil {
			log.Fatalf("Error while scanning row: %v", err)
		}

		productsRowsChan <- item
	}

	close(productsRowsChan)
	wg.Wait()

	if err := productsRows.Close(); err != nil {
		log.Fatal(err)
	}
}

func processProduct(productsChan chan *Product, wg *sync.WaitGroup, shopProducts []*Product, bar *progressbar.ProgressBar) {
	defer wg.Done()

	for item := range productsChan {
		shopProducts = append(shopProducts, item)
		bar.Add(1)
	}
}
