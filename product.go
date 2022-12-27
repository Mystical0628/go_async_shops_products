package main

import (
	"encoding/xml"
	"log"
	"os"
)

type Shop struct {
	XMLName  xml.Name   `xml:"shop"`
	Id       int        `xml:"id,attr"`
	Name     string     `xml:"name"`
	Url      string     `xml:"url"`
	OpensAt  string     `xml:"open"`
	ClosesAt string     `xml:"close"`
	Products []*Product `xml:"offers"`
}

func (p *Shop) Render(file *os.File) {
	data, err := xml.MarshalIndent(p, " ", "    ")

	if err != nil {
		log.Fatalf("Error while xml.MarshalIndent: %v", err)
	}

	_, err = file.Write(data)

	if err != nil {
		log.Fatalf("Error while write to file: %v", err)
	}
}

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
	data, err := xml.MarshalIndent(p, " ", "    ")

	if err != nil {
		log.Fatalf("Error while xml.MarshalIndent: %v", err)
	}

	_, err = file.Write(data)

	if err != nil {
		log.Fatalf("Error while write to file: %v", err)
	}
}
