package main

import (
	"database/sql"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"log"
	"strconv"
	"strings"
	"time"
	// "github.com/douglasmg7/currency"
)

///////////////////////////////////////////////////////////////////////////////////////////////////
// PRODUCT
///////////////////////////////////////////////////////////////////////////////////////////////////
// Motospeed product.
type MotospeedProduct struct {
	Code           sql.NullString `db:"code"` // From dealer.
	Title          sql.NullString `db:"title"`
	Desc           sql.NullString `db:"desc"`
	Price          sql.NullInt64  `db:"price"`
	Stock          sql.NullInt64  `db:"stock"`
	ZunkaProductId sql.NullString `db:"zunka_product_id"`
	CreatedAt      time.Time      `db:"created_at"`
	// ChangedAt              time.Time      `db:"changed_at"`
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
