package app

import (
	"database/sql"
	"log"
)

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

func (app *app) getShopProductsTotal(shopId int) int {
	var total int

	app.db.QueryRow(`
		SELECT COUNT(*) AS total
		FROM products
		WHERE shop_id = ?`, shopId).Scan(&total)

	return total
}

func (app *app) getShopProductsRows(shopId int, limit int) *sql.Rows {
	query := `
		SELECT id, name, description, price
		FROM products 
		WHERE shop_id = ?
	`
	args := []any{shopId}

	if limit != 0 {
		query += "LIMIT ?"
		args = append(args, limit)
	}

	rows, err := app.db.Query(query, args...)

	if err != nil {
		log.Fatalf("Error getShopProductsRows: %v", err)
	}

	return rows
}

func (app *app) getProductsTotal() int {
	var total int

	app.db.QueryRow(`
		SELECT COUNT(*) AS total
		FROM products
	`).Scan(&total)

	return total
}

func (app *app) getProductsRows(limit int) *sql.Rows {
	query := `
		SELECT id, name, description, price
		FROM products 
		ORDER BY shop_id, id
	`
	var args []any

	if limit != 0 {
		query += "LIMIT ?"
		args = append(args, limit)
	}

	rows, err := app.db.Query(query, args...)

	if err != nil {
		log.Fatalf("Error getProductsRows: %v", err)
	}

	return rows
}
