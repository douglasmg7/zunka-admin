package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/douglasmg7/currency"
)

var db *sql.DB

///////////////////////////////////////////////////////////////////////////////////////////////////
// PRODUCT
///////////////////////////////////////////////////////////////////////////////////////////////////
// Aldo product.
type AllnationsProduct struct {
	ZunkaProductId       string            `db:"zunka_product_id"`
	Code                 string            `db:"code"` // From dealer.
	Description          string            `db:"description"`
	Timestamp            string            `db:"timestamp"`
	Department           string            `db:"department"`
	Category             string            `db:"category"`
	SubCategory          string            `db:"sub_category"`
	Maker                string            `db:"maker"`
	TechnicalDescription string            `db:"technical_description"`
	UrlImage             string            `db:"url_image"`
	PartNumber           string            `db:"part_number"`
	Ean                  string            `db:"ean"`
	Ncm                  string            `db:"ncm"`
	PriceSale            currency.Currency `db:"price_sale"`
	PriceWithoutSt       currency.Currency `db:"price_without_st"`
	IcmsStTaxation       bool              `db:"icms_st_taxation"`
	WarrantyMonth        int               `db:"warranty_month"`
	LengthMm             int               `db:"length_mm"`
	WidthMm              int               `db:"width_mm"`
	HeightMm             int               `db:"height_mm"`
	WeightG              int               `db:"weight_g"`
	Active               bool              `db:"active"`
	Availability         bool              `db:"availability"` // Months.
	Origin               string            `db:"origin"`
	StockOrigin          string            `db:"stock_origin"`
	StockQty             int               `db:"stock_qty"`
	CreatedAt            time.Time         `db:"created_at"`
	ChangedAt            time.Time         `db:"changed_at"`
	CheckedAt            time.Time         `db:"checked_at"`
	RemovedAt            time.Time         `db:"removed_at"`
}

// Define product status.
func (p *AllnationsProduct) Status(validDate time.Time) string {
	if !p.RemovedAt.IsZero() {
		return "removed"
	}
	if !p.Availability {
		return "unavailable"
	}
	// Status have a valid time for not created products at zunkasite.
	if p.ZunkaProductId == "" && p.ChangedAt.Before(validDate) {
		return ""
	}
	// For created products at zunkasite, only clean status by user.
	if !p.CheckedAt.IsZero() && p.CheckedAt.After(p.ChangedAt) {
		return ""
	}
	// New.
	if p.ChangedAt.Equal(p.CreatedAt) {
		return "new"
	} else {
		return "changed"
	}
}

// Diff check if products are different.
func (p *AllnationsProduct) Diff(pn *AllnationsProduct) bool {
	if p.ZunkaProductId != pn.ZunkaProductId {
		return true
	}
	if p.Code != pn.Code {
		return true
	}
	if p.Description != pn.Description {
		return true
	}
	if p.Timestamp != pn.Timestamp {
		return true
	}
	if p.Department != pn.Department {
		return true
	}
	if p.Category != pn.Category {
		return true
	}
	if p.SubCategory != pn.SubCategory {
		return true
	}
	if p.Maker != pn.Maker {
		return true
	}
	if p.TechnicalDescription != pn.TechnicalDescription {
		return true
	}
	if p.UrlImage != pn.UrlImage {
		return true
	}
	if p.PartNumber != pn.PartNumber {
		return true
	}
	if p.Ean != pn.Ean {
		return true
	}
	if p.Ncm != pn.Ncm {
		return true
	}
	if p.PriceSale != pn.PriceSale {
		return true
	}
	if p.PriceWithoutSt != pn.PriceWithoutSt {
		return true
	}
	if p.IcmsStTaxation != pn.IcmsStTaxation {
		return true
	}
	if p.WarrantyMonth != pn.WarrantyMonth {
		return true
	}
	if p.LengthMm != pn.LengthMm {
		return true
	}
	if p.WidthMm != pn.WidthMm {
		return true
	}
	if p.HeightMm != pn.HeightMm {
		return true
	}
	if p.WeightG != pn.WeightG {
		return true
	}
	if p.Active != pn.Active {
		return true
	}
	if p.Availability != pn.Availability {
		return true
	}
	if p.Origin != pn.Origin {
		return true
	}
	if p.StockOrigin != pn.StockOrigin {
		return true
	}
	if p.StockQty != pn.StockQty {
		return true
	}
	if p.CreatedAt != pn.CreatedAt {
		return true
	}
	if p.ChangedAt != pn.ChangedAt {
		return true
	}
	if p.CheckedAt != pn.CheckedAt {
		return true
	}
	if p.RemovedAt != pn.RemovedAt {
		return true
	}
	return false
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// CATEGORY
///////////////////////////////////////////////////////////////////////////////////////////////////
type AllnationsCategory struct {
	Name        string `db:"name"`
	ProductsQty int    `db:"products_qty"`
	Selected    bool   `db:"selected"`
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// MAKERS
///////////////////////////////////////////////////////////////////////////////////////////////////
type AllnationsMaker struct {
	Name        string `db:"name"`
	ProductsQty int    `db:"products_qty"`
	Selected    bool   `db:"selected"`
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// FILTERS
///////////////////////////////////////////////////////////////////////////////////////////////////
type AllnationsFilters struct {
	OnlyActive       bool
	OnlyAvailability bool
	MinPrice         string
	MaxPrice         string
	PathData         string
	SqlFilter        string
}

// Load allnations filters.
func LoadAllnationsFilters(path string) *AllnationsFilters {
	allnationsFilters := AllnationsFilters{}
	allnationsFilters.PathData = path
	// Read allnations filters.
	err = readGob(allnationsFilters.PathData, &allnationsFilters)
	if err != nil {
		log.Printf("[warn] Not found Allnations filters data")
		allnationsFilters.MinPrice = "100"
		allnationsFilters.MaxPrice = "100000000"
	}
	// log.Printf("Allnations filters: %v", allnationsFilters)
	return &allnationsFilters
}

// Save filters.
func (f *AllnationsFilters) Save() error {
	f.UpdateSqlFilter()
	return writeGob(f.PathData, f)
}

// Update sql filter.
func (f *AllnationsFilters) UpdateSqlFilter() {
	// Filters.
	filtersArray := []string{}
	// Only active.
	if f.OnlyActive {
		filtersArray = append(filtersArray, " active = true ")
	}
	// Only availability.
	if f.OnlyAvailability {
		filtersArray = append(filtersArray, " availability = true ")
	}
	// Min price.
	minPrice, err := strconv.ParseUint(f.MinPrice, 10, 64)
	if err != nil {
		log.Panicf("[error] Updatindg sql filter: %s", err)
	}
	filtersArray = append(filtersArray, fmt.Sprintf(" price_sale >= %v ", minPrice))
	// Max price.
	maxPrice, err := strconv.ParseUint(f.MaxPrice, 10, 64)
	if err != nil {
		log.Panicf("[error] Updatindg sql filter: %s", err)
	}
	filtersArray = append(filtersArray, fmt.Sprintf(" price_sale <= %v ", maxPrice))

	f.SqlFilter = strings.Join(filtersArray, " AND ")
}

// Validation filters.
type AllnationsFiltersValidation struct {
	MinPrice string
	MaxPrice string
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// SELECTED CATEGORIES
///////////////////////////////////////////////////////////////////////////////////////////////////
// Allnations Filters.
type AllnationsSelectedCategories struct {
	Categories    []string
	SqlCategories string
	PathData      string
}

// Load allnations filters.
func LoadAllnationsSelectedCategories(path string) *AllnationsSelectedCategories {
	allnationsSelectedCategories := AllnationsSelectedCategories{}
	allnationsSelectedCategories.PathData = path
	// Read allnations selected categories.
	err = readGob(allnationsSelectedCategories.PathData, &allnationsSelectedCategories)
	if err != nil {
		log.Printf("[warn] Not found Allnations selected categories data")
	}
	// log.Printf("Allnations selected categories: %v", allnationsSelectedCategories)
	return &allnationsSelectedCategories
}

// Save filters.
func (f *AllnationsSelectedCategories) Save() error {
	f.UpdateSqlCategories()
	return writeGob(f.PathData, f)
}

// Update sql filter.
func (f *AllnationsSelectedCategories) UpdateSqlCategories() {
	categories := []string{}
	for _, category := range f.Categories {
		categories = append(categories, fmt.Sprintf("\"%s\"", category))
	}
	f.SqlCategories = strings.Join(categories, ", ")
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// SELECTED MAKERS
///////////////////////////////////////////////////////////////////////////////////////////////////
// Allnations Filters.
type AllnationsSelectedMakers struct {
	Makers    []string
	SqlMakers string
	PathData  string
}

// Load allnations filters.
func LoadAllnationsSelectedMakers(path string) *AllnationsSelectedMakers {
	allnationsSelectedMakers := AllnationsSelectedMakers{}
	allnationsSelectedMakers.PathData = path
	// Read allnations selected categories.
	err = readGob(allnationsSelectedMakers.PathData, &allnationsSelectedMakers)
	if err != nil {
		log.Printf("[warn] Not found Allnations selected makers data")
	}
	// log.Printf("Allnations selected makers: %v", allnationsSelectedMakers)
	return &allnationsSelectedMakers
}

// Save filters.
func (f *AllnationsSelectedMakers) Save() error {
	f.UpdateSqlMakers()
	return writeGob(f.PathData, f)
}

// Update sql Makers.
func (f *AllnationsSelectedMakers) UpdateSqlMakers() {
	makers := []string{}
	for _, maker := range f.Makers {
		makers = append(makers, fmt.Sprintf("\"%s\"", maker))
	}
	f.SqlMakers = strings.Join(makers, ", ")
}
