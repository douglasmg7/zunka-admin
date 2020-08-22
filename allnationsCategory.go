package main

// Aldo product.
type AllnationsCategory struct {
	Name        string `db:"name"`
	Text        string `db:"text"`
	ProductsQty int    `db:"products_qty"`
	Selected    bool   `db:"selected"`
}
