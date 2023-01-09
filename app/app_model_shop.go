package app

import (
	"encoding/xml"
	"log"
)

type Shop struct {
	XMLName  xml.Name   `xml:"shop"`
	Id       int        `xml:"id,attr"`
	Name     string     `xml:"name"`
	Url      string     `xml:"url"`
	OpensAt  string     `xml:"working_time>open"`
	ClosesAt string     `xml:"working_time>close"`
	Products []*Product `xml:"offers>item"`
}

func (app *app) getShops(limit int) []*Shop {
	query := `
		SELECT id, name, url, opens_at, closes_at
		FROM shops 
		WHERE opens_at <= ? AND closes_at >= ? 
	`
	args := []any{app.timeFormatted, app.timeFormatted}

	if limit != 0 {
		query += "LIMIT ?"
		args = append(args, limit)
	}

	rows, err := app.db.Query(query, args...)

	if err != nil {
		log.Fatalf("Error while selecting shops: %v", err)
	}

	defer rows.Close()

	var shops []*Shop

	for rows.Next() {
		item := &Shop{}
		if err := rows.Scan(&item.Id, &item.Name, &item.Url, &item.OpensAt, &item.ClosesAt); err != nil {
			log.Fatalf("Error while scanning row: %v", err)
		}
		shops = append(shops, item)
	}

	err = rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	return shops
}
