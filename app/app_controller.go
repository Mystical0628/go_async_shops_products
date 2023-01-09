package app

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/grokify/html-strip-tags-go"
	"github.com/schollz/progressbar/v3"
	"log"
	"os"
	"strconv"
	"sync"
)

func (app *app) ActionStream() {
	var err error
	enc := app.initFileEncoder("shops_products_stream.xml")

	shops := app.getShops(*app.flagShops)

	shopsElement := xml.StartElement{
		Name: xml.Name{Local: "shops"},
	}

	shopElement := xml.StartElement{
		Name: xml.Name{Local: "shop"},
		Attr: []xml.Attr{
			{Name: xml.Name{Local: "id"}, Value: ""},
		},
	}

	workingTimeElement := xml.StartElement{Name: xml.Name{Local: "working_time"}}
	offersElement := xml.StartElement{Name: xml.Name{Local: "offers"}}

	enc.EncodeToken(shopsElement)
	for _, shop := range shops {
		shopElement.Attr[0].Value = strconv.Itoa(shop.Id)
		enc.EncodeToken(shopElement)

		enc.EncodeElement(shop.Name, xml.StartElement{Name: xml.Name{Local: "name"}})
		enc.EncodeElement(shop.Url, xml.StartElement{Name: xml.Name{Local: "url"}})

		enc.EncodeToken(workingTimeElement)
		enc.EncodeElement(shop.OpensAt[:5], xml.StartElement{Name: xml.Name{Local: "open"}})
		enc.EncodeElement(shop.ClosesAt[:5], xml.StartElement{Name: xml.Name{Local: "close"}})
		enc.EncodeToken(workingTimeElement.End())

		enc.EncodeToken(offersElement)

		err = app.processIndex(enc, shop)

		enc.EncodeToken(offersElement.End())
		enc.EncodeToken(shopElement.End())

		if err != nil {
			log.Fatalf("Error ActionIndex: %v", err)
		}
	}
	enc.EncodeToken(shopsElement.End())

	if err := enc.Flush(); err != nil {
		log.Fatalf("Error ActionIndex: %v", err)
	}
}

func (app *app) ActionByShops() {
	enc := app.initFileEncoder("shops_products_by_shops.xml")

	shops := app.getShops(*app.flagShops)

	shopsElement := xml.StartElement{
		Name: xml.Name{Local: "shops"},
	}

	enc.EncodeToken(shopsElement)
	for _, shop := range shops {
		err := app.processAll(enc, shop)

		if err != nil {
			log.Fatalf("Error ActionIndex: %v", err)
		}

		enc.Encode(shop)
	}
	enc.EncodeToken(shopsElement.End())

	if err := enc.Flush(); err != nil {
		log.Fatalf("Error ActionIndex: %v", err)
	}

	// shops := app.getShops(*app.flagShops)
}

func (app *app) initByteEncoder(buffer *bytes.Buffer) *xml.Encoder {
	enc := xml.NewEncoder(bufio.NewWriter(buffer))
	enc.Indent(" ", "    ")

	return enc
}

func (app *app) initFileEncoder(filename string) *xml.Encoder {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("Error while opening file: %v", err)
	}
	enc := xml.NewEncoder(file)
	enc.Indent(" ", "    ")

	return enc
}

func (app *app) initBarWgMutexChan(barMax int, barName string) (
	*progressbar.ProgressBar,
	*sync.WaitGroup,
	*sync.Mutex,
	chan Product,
) {
	threads := *app.flagThreads
	wg := &sync.WaitGroup{}
	wg.Add(threads)

	return progressbar.Default(int64(barMax), barName), wg, &sync.Mutex{}, make(chan Product, threads)
}

func (app *app) processIndex(enc *xml.Encoder, shop *Shop) error {
	productsTotal := *app.flagProducts
	productsRows := app.getShopProductsRows(shop.Id, productsTotal)
	defer productsRows.Close()

	if productsTotal == 0 {
		productsTotal = app.getShopProductsTotal(shop.Id)
	}

	bar, wg, mutex, productsChan := app.initBarWgMutexChan(productsTotal, fmt.Sprintf("Shop %d:", shop.Id))

	for i := 0; i < *app.flagThreads; i++ {
		go productWorkerIndex(bar, wg, mutex, productsChan, enc)
	}

	for productsRows.Next() {
		item := Product{}
		if err := productsRows.Scan(&item.Id, &item.Name, &item.Description, &item.Price); err != nil {
			return errors.New(fmt.Sprintf("processProducts: scan productsRows: %v", err))
		}

		productsChan <- item
	}

	close(productsChan)

	if err := productsRows.Close(); err != nil {
		return errors.New(fmt.Sprintf("processProducts: close productsRows: %v", err))
	}

	wg.Wait()

	return nil
}

func (app *app) processAll(enc *xml.Encoder, shop *Shop) error {
	productsTotal := *app.flagProducts
	productsRows := app.getShopProductsRows(shop.Id, productsTotal)
	defer productsRows.Close()

	if productsTotal == 0 {
		productsTotal = app.getShopProductsTotal(shop.Id)
	}

	bar, wg, _, productsChan := app.initBarWgMutexChan(productsTotal, fmt.Sprintf("Shop %d:", shop.Id))

	for i := 0; i < *app.flagThreads; i++ {
		go productWorkerAll(bar, wg, productsChan, shop)
	}

	for productsRows.Next() {
		item := Product{}
		if err := productsRows.Scan(&item.Id, &item.Name, &item.Description, &item.Price); err != nil {
			return errors.New(fmt.Sprintf("processProducts: scan productsRows: %v", err))
		}

		productsChan <- item
	}

	close(productsChan)

	if err := productsRows.Close(); err != nil {
		return errors.New(fmt.Sprintf("processProducts: close productsRows: %v", err))
	}

	wg.Wait()

	return nil
}

func productWorkerIndex(
	bar *progressbar.ProgressBar,
	wg *sync.WaitGroup,
	mutex *sync.Mutex,
	productsChan chan Product,
	enc *xml.Encoder,
) {
	defer wg.Done()

	for item := range productsChan {
		item.Description = strip.StripTags(item.Description)
		mutex.Lock()
		enc.Encode(item)
		mutex.Unlock()
		bar.Add(1)
	}
}

func productWorkerAll(
	bar *progressbar.ProgressBar,
	wg *sync.WaitGroup,
	productsChan chan Product,
	shop *Shop,
) {
	defer wg.Done()

	for item := range productsChan {
		item.Description = strip.StripTags(item.Description)
		shop.Products = append(shop.Products, &item)
		bar.Add(1)
	}
}
