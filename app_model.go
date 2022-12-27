package main

import (
	"database/sql"
	"log"
)

func (app *app) getShops() []*Shop {
	rows, err := app.db.Query(`
		SELECT id, name, url, opens_at, closes_at
		FROM shops 
		WHERE opens_at <= ? AND closes_at >= ? 
		`, app.timeFormatted, app.timeFormatted)

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

func (app *app) getProductsTotal() int {
	var total int

	app.db.QueryRow(`
		SELECT COUNT(*) AS total 
		FROM products p
		    LEFT JOIN shops s ON (p.shop_id = s.id)
		WHERE s.opens_at <= ? AND s.closes_at >= ?`, app.timeFormatted, app.timeFormatted).Scan(&total)

	return total
}

func (app *app) getProducts(start int, limit int) []*Product {
	rows, err := app.db.Query(`
		SELECT p.id, p.shop_id, p.name, p.description, p.price
		FROM products p 
		    LEFT JOIN shops s ON (p.shop_id = s.id)
		WHERE s.opens_at <= ? AND s.closes_at >= ? 
		LIMIT ?, ?`, app.timeFormatted, app.timeFormatted, start, limit)

	if err != nil {
		log.Fatalf("Error while selecting shops: %v", err)
	}

	defer rows.Close()

	var products []*Product

	for rows.Next() {
		item := &Product{}
		if err := rows.Scan(&item.Id, &item.ShopId, &item.Name, &item.Description, &item.Price); err != nil {
			log.Fatalf("Error while scanning row: %v", err)
		}
		products = append(products, item)
	}

	err = rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	return products
}

func (app *app) getProductsRows(start int, limit int) *sql.Rows {
	rows, err := app.db.Query(`
		SELECT p.id, p.shop_id, p.name, p.description, p.price
		FROM products p 
		    LEFT JOIN shops s ON (p.shop_id = s.id)
		WHERE s.opens_at <= ? AND s.closes_at >= ? 
		LIMIT ?, ?`, app.timeFormatted, app.timeFormatted, start, limit)

	if err != nil {
		log.Fatalf("Error while selecting shops: %v", err)
	}

	return rows
}

func (app *app) getShopProductsRows(shopId int) *sql.Rows {
	rows, err := app.db.Query(`
		SELECT p.id, p.name, p.description, p.price
		FROM products p
		WHERE p.shop_id = ?`, shopId)

	if err != nil {
		log.Fatalf("Error while selecting shops: %v", err)
	}

	return rows
}
