package app

import (
	"database/sql"
	"encoding/xml"
	"log"
)

type Product struct {
	XMLName     xml.Name `xml:"item"`
	Id          int      `xml:"id,attr"`
	Name        string   `xml:"name"`
	Description string   `xml:"description"`
	Price       float32  `xml:"price"`
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