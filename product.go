package main

import (
	"encoding/xml"
	"log"
	"os"
)

type Products struct {
	XMLName  xml.Name   `xml:"offers"`
	Products []*Product `xml:"item"`
}

func (p *Products) Render(file *os.File) {
	data, err := xml.MarshalIndent(p, " ", "    ")

	if err != nil {
		log.Fatalf("Error while xml.MarshalIndent: %v", err)
	}

	_, err = file.Write(data)

	if err != nil {
		log.Fatalf("Error while write to file: %v", err)
	}
}

type Product struct {
	XMLName     xml.Name `xml:"item"`
	Id          int      `xml:"id,attr"`
	ShopId      int      `xml:"shop-id"`
	Price       float32  `xml:"price"`
	Name        string   `xml:"name"`
	Description string   `xml:"description"`
}

func (p *Product) Render(file *os.File) {
	file, err := os.OpenFile("products_bundle.xml", os.O_APPEND|os.O_WRONLY, os.ModePerm)

	data, err := xml.MarshalIndent(p, " ", "    ")

	if err != nil {
		log.Fatalf("Error while xml.MarshalIndent: %v", err)
	}

	_, err = file.Write(data)

	if err != nil {
		log.Fatalf("Error while write to file: %v", err)
	}

	file.Close()
}
