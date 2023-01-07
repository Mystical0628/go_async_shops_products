package app

import (
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

func (app *app) initEncoder(filename string) *xml.Encoder {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("Error while opening file: %v", err)
	}
	enc := xml.NewEncoder(file)
	enc.Indent(" ", "    ")

	return enc
}

func (app *app) ActionIndex() {
	var err error
	enc := app.initEncoder("shops_products.xml")

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

		err = app.process(enc, shop)

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

func (app *app) ActionAll() {
	log.Println("Hello ActionAll")
	// enc := app.initEncoder("shops_products.xml")

	// shops := app.getShops(*app.flagShops)
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

func (app *app) process(enc *xml.Encoder, shop *Shop) error {
	productsTotal := *app.flagProducts
	productsRows := app.getShopProductsRows(shop.Id, productsTotal)
	defer productsRows.Close()

	if productsTotal == 0 {
		productsTotal = app.getShopProductsTotal(shop.Id)
	}

	bar, wg, mutex, productsChan := app.initBarWgMutexChan(productsTotal, fmt.Sprintf("Shop %d:", shop.Id))

	for i := 0; i < *app.flagThreads; i++ {
		go productWorker(bar, wg, mutex, productsChan, enc)
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

func productWorker(
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
