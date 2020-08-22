package main

// Aldo product.
type AllnationsCategory struct {
	Name        string `db:"name"`
	ProductsQty int    `db:"products_qty"`
	Selected    bool   `db:"selected"`
}
