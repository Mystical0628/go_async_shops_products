package app

import (
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"log"
	"os"
	"sync"
)

func (app *app) ActionIndex() {
	var err error
	enc := app.initEncoder("products_all.xml")

	startElement := xml.StartElement{
		Name: xml.Name{Local: "shops"},
	}

	if err := enc.EncodeToken(startElement); err != nil {
		log.Fatalf("Error ActionIndex: %v", err)
	}

	shops := app.getShops(*app.flagShops)

	for _, shop := range shops {
		startElement := xml.StartElement{
			Name: xml.Name{Local: "shop"},
			Attr:	[]xml.Attr{
				{Name: xml.Name{Local: "id"}, Value: string(shop.Id)},
			},
		}

		err = enc.EncodeToken(startElement)
		if err == nil {
			err = app.processShopProducts(shop, enc)
		}
		if err == nil {
			err = enc.EncodeToken(startElement.End())
		}
		if err != nil {
			log.Fatalf("Error ActionIndex: %v", err)
		}
	}

	if err := enc.EncodeToken(startElement.End()); err != nil {
		log.Fatalf("Error ActionIndex: %v", err)
	}

	if err := enc.Flush(); err != nil {
		log.Fatalf("Error ActionIndex: %v", err)
	}
}

func (app *app) processShopProducts(shop *Shop, enc *xml.Encoder) error {
	return app.process(enc, app.getShopProductsTotal(shop.Id), app.getShopProductsRows(shop.Id, *app.flagProducts),
		fmt.Sprintf("Shop %d:", shop.Id))
}

func (app *app) ActionAll() {
	enc := app.initEncoder("products_all.xml")

	startElement := xml.StartElement{
		Name: xml.Name{Local: "items"},
	}

	err := enc.EncodeToken(startElement)

	if err == nil {
		err = app.processProducts(enc)
	}

	if err == nil {
		err = enc.EncodeToken(startElement.End())
	}

	if err == nil {
		err = enc.Flush()
	}

	if err != nil {
		log.Fatalf("Error ActionAll: %v", err)
	}
}

func (app *app) processProducts(enc *xml.Encoder) error {
	return app.process(enc, app.getProductsTotal(), app.getProductsRows(*app.flagProducts), "All Products:")
}

func (app *app) initBarWgMutexChan(barMax int, barName string) (
	*progressbar.ProgressBar,
	*sync.WaitGroup,
	*sync.Mutex,
	chan *Product,
) {
	threads := *app.flagThreads
	wg := &sync.WaitGroup{}
	wg.Add(threads)

	return progressbar.Default(int64(barMax), barName), wg, &sync.Mutex{}, make(chan *Product, threads)
}

func (app *app) process(enc *xml.Encoder, productsTotal int, productsRows *sql.Rows, barName string) error {
	defer productsRows.Close()

	bar, wg, mutex, productsChan := app.initBarWgMutexChan(productsTotal, barName)

	for i := 0; i < *app.flagThreads; i++ {
		go productWorker(bar, wg, mutex, productsChan, enc)
	}

	for productsRows.Next() {
		item := &Product{}
		if err := productsRows.Scan(&item.Id, &item.Name, &item.Description, &item.Price); err != nil {
			return errors.New(fmt.Sprintf("processProducts: scan productsRows: %v", err))
		}

		productsChan <- item
	}

	close(productsChan)
	wg.Wait()

	if err := productsRows.Close(); err != nil {
		return errors.New(fmt.Sprintf("processProducts: close productsRows: %v", err))
	}

	return nil
}

func productWorker(
	bar *progressbar.ProgressBar,
	wg *sync.WaitGroup,
	mutex *sync.Mutex,
	productsChan chan *Product,
	enc *xml.Encoder,
	) {
	defer wg.Done()

	for item := range productsChan {
		mutex.Lock()
		enc.Encode(item)
		mutex.Unlock()
		bar.Add(1)
	}
}

func (app *app) initEncoder(filename string) *xml.Encoder {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("Error while opening file: %v", err)
	}
	enc := xml.NewEncoder(file)
	enc.Indent(" ", "    ")

	return enc
}