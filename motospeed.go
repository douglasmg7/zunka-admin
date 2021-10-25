package main

import (
	"database/sql"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"log"
	"strconv"
	"strings"
	// "github.com/douglasmg7/currency"
)

///////////////////////////////////////////////////////////////////////////////////////////////////
// PRODUCT
///////////////////////////////////////////////////////////////////////////////////////////////////
// Motospeed product.
type MotospeedProduct struct {
	Sku            sql.NullString `db:"sku"` // From dealer.
	Title          sql.NullString `db:"title"`
	Description    sql.NullString `db:"description"`
	Ean            sql.NullString `db:"ean"`
	Model          sql.NullString `db:"model"`
	Connection     sql.NullString `db:"connection"`
	Compatibility  sql.NullString `db:"compatibility"`
	Curve          sql.NullString `db:"curve"`
	NCM            sql.NullString `db:"ncm"`
	MasterBox      sql.NullInt64  `db:"master_box"`
	WeightKG       sql.NullInt64  `db:"weight_kg"`
	LengthCM       sql.NullInt64  `db:"length_cm"`
	WidthCM        sql.NullInt64  `db:"width_cm"`
	DepthCM        sql.NullInt64  `db:"depth_cm"`
	IPI            sql.NullInt64  `db:"ipi"`
	Price100       sql.NullInt64  `db:"price_100"`
	PriceDist100   sql.NullInt64  `db:"price_dist_100"`
	PriceSell100   sql.NullInt64  `db:"price_sell_100"`
	Stock          sql.NullInt64  `db:"stock"`
	ZunkaProductId sql.NullString `db:"zunka_product_id"`
	CreatedAt      sql.NullTime   `db:"created_at"`
	ChangedAt      sql.NullTime   `db:"changed_at"`
	RemovedAt      sql.NullTime   `db:"removed_at"`
}

// Define product status.
func (p *MotospeedProduct) Status() string {
	if p.ZunkaProductId.String == "" {
		return "no-registered"
	}
	return "registered"
}

// Process Br Currency.
func (p *MotospeedProduct) ProcessBrCurrency(val sql.NullInt64) string {
	if val.Valid {
		printer := message.NewPrinter(language.Portuguese)
		return printer.Sprintf("R$ %.2f", float64(val.Int64)/100)
	} else {
		return "NULL"
	}
}

// Process string, show NULL.
func (p *MotospeedProduct) ProcessString(val sql.NullString) string {
	if val.Valid {
		return val.String
	} else {
		return "NULL"
	}
}

// Process int64.
func (p *MotospeedProduct) ProcessInt64(val sql.NullInt64) string {
	if val.Valid {
		return fmt.Sprintf("%d", val.Int64)
	} else {
		return "NULL"
	}
}

// Process time.
func (p *MotospeedProduct) ProcessTime(val sql.NullTime) string {
	if val.Valid {
		return fmt.Sprintf("%d", val.Time)
	} else {
		return "NULL"
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// CATEGORY
///////////////////////////////////////////////////////////////////////////////////////////////////
type MotospeedCategory struct {
	Name        string `db:"name"`
	ProductsQty int    `db:"products_qty"`
	Selected    bool   `db:"selected"`
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// FILTERS
///////////////////////////////////////////////////////////////////////////////////////////////////
type MotospeedFilters struct {
	MinPrice  string
	MaxPrice  string
	PathData  string
	SqlFilter string
}

// Load filters.
func LoadMotospeedFilters(path string) *MotospeedFilters {
	MotospeedFilters := MotospeedFilters{}
	MotospeedFilters.PathData = path
	// Read filters.
	err = readGob(MotospeedFilters.PathData, &MotospeedFilters)
	if err != nil {
		log.Printf("[warn] Not found Motospeed filters data")
		MotospeedFilters.MinPrice = "100"
		MotospeedFilters.MaxPrice = "100000000"
	}
	// log.Printf("Motospeed filters: %v", MotospeedFilters)
	return &MotospeedFilters
}

// Save filters.
func (f *MotospeedFilters) Save() error {
	f.UpdateSqlFilter()
	return writeGob(f.PathData, f)
}

// Update sql filter.
func (f *MotospeedFilters) UpdateSqlFilter() {
	// Filters.
	filtersArray := []string{}
	// Min price.
	minPrice, err := strconv.ParseUint(f.MinPrice, 10, 64)
	if err != nil {
		log.Panicf("[error] Updatindg sql filter: %s", err)
	}
	filtersArray = append(filtersArray, fmt.Sprintf(" vl_item >= %v ", minPrice))
	// Max price.
	maxPrice, err := strconv.ParseUint(f.MaxPrice, 10, 64)
	if err != nil {
		log.Panicf("[error] Updatindg sql filter: %s", err)
	}
	filtersArray = append(filtersArray, fmt.Sprintf(" vl_item <= %v ", maxPrice))

	f.SqlFilter = strings.Join(filtersArray, " AND ")
}

// Validation filters.
type MotospeedFiltersValidation struct {
	MinPrice string
	MaxPrice string
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// SELECTED CATEGORIES
///////////////////////////////////////////////////////////////////////////////////////////////////
// Filters.
type MotospeedSelectedCategories struct {
	Categories    []string
	SqlCategories string
	PathData      string
}

// Load filters.
func LoadMotospeedSelectedCategories(path string) *MotospeedSelectedCategories {
	motospeedSelectedCategories := MotospeedSelectedCategories{}
	motospeedSelectedCategories.PathData = path
	// Read selected categories.
	err = readGob(motospeedSelectedCategories.PathData, &motospeedSelectedCategories)
	if err != nil {
		// log.Printf("[warn] Not found Motospeed selected categories data")
	}
	// log.Printf("Motospeed selected categories: %v", motospeedSelectedCategories)
	return &motospeedSelectedCategories
}

// Save filters.
func (f *MotospeedSelectedCategories) Save() error {
	f.UpdateSqlCategories()
	return writeGob(f.PathData, f)
}

// Update sql filter.
func (f *MotospeedSelectedCategories) UpdateSqlCategories() {
	categories := []string{}
	for _, category := range f.Categories {
		categories = append(categories, fmt.Sprintf("\"%s\"", category))
	}
	f.SqlCategories = strings.Join(categories, ", ")
}
