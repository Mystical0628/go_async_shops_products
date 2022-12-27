package app

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

func (app *app) getShopProductsTotal(shopId int) int {
	var total int

	app.db.QueryRow(`
		SELECT COUNT(*) AS total
		FROM products
		WHERE shop_id = ?`, shopId).Scan(&total)

	return total
}

func (app *app) getShopProductsRows(shopId int) *sql.Rows {
	rows, err := app.db.Query(`
		SELECT p.id, p.name, p.description, p.price
		FROM products p
		WHERE p.shop_id = ? LIMIT 0, 10`, shopId)

	if err != nil {
		log.Fatalf("Error while selecting shops: %v", err)
	}

	return rows
}
