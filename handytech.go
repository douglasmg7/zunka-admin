package main

import (
	// "database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/douglasmg7/currency"
)

// var db *sql.DB

///////////////////////////////////////////////////////////////////////////////////////////////////
// PRODUCT
///////////////////////////////////////////////////////////////////////////////////////////////////
// Aldo product.
type HandytechProduct struct {
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
func (p *HandytechProduct) Status(validDate time.Time) string {
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
func (p *HandytechProduct) Diff(pn *HandytechProduct) bool {
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
type HandytechCategory struct {
	Name        string `db:"name"`
	ProductsQty int    `db:"products_qty"`
	Selected    bool   `db:"selected"`
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// MAKERS
///////////////////////////////////////////////////////////////////////////////////////////////////
type HandytechMaker struct {
	Name        string `db:"name"`
	ProductsQty int    `db:"products_qty"`
	Selected    bool   `db:"selected"`
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// FILTERS
///////////////////////////////////////////////////////////////////////////////////////////////////
type HandytechFilters struct {
	OnlyActive       bool
	OnlyAvailability bool
	MinPrice         string
	MaxPrice         string
	PathData         string
	SqlFilter        string
}

// Load handytech filters.
func LoadHandytechFilters(path string) *HandytechFilters {
	handytechFilters := HandytechFilters{}
	handytechFilters.PathData = path
	// Read handytech filters.
	err = readGob(handytechFilters.PathData, &handytechFilters)
	if err != nil {
		log.Printf("[warn] Not found Handytech filters data")
		handytechFilters.MinPrice = "100"
		handytechFilters.MaxPrice = "100000000"
	}
	// log.Printf("Handytech filters: %v", handytechFilters)
	return &handytechFilters
}

// Save filters.
func (f *HandytechFilters) Save() error {
	f.UpdateSqlFilter()
	return writeGob(f.PathData, f)
}

// Update sql filter.
func (f *HandytechFilters) UpdateSqlFilter() {
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
type HandytechFiltersValidation struct {
	MinPrice string
	MaxPrice string
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// SELECTED CATEGORIES
///////////////////////////////////////////////////////////////////////////////////////////////////
// Handytech Filters.
type HandytechSelectedCategories struct {
	Categories    []string
	SqlCategories string
	PathData      string
}

// Load handytech filters.
func LoadHandytechSelectedCategories(path string) *HandytechSelectedCategories {
	handytechSelectedCategories := HandytechSelectedCategories{}
	handytechSelectedCategories.PathData = path
	// Read handytech selected categories.
	err = readGob(handytechSelectedCategories.PathData, &handytechSelectedCategories)
	if err != nil {
		log.Printf("[warn] Not found Handytech selected categories data")
	}
	// log.Printf("Handytech selected categories: %v", handytechSelectedCategories)
	return &handytechSelectedCategories
}

// Save filters.
func (f *HandytechSelectedCategories) Save() error {
	f.UpdateSqlCategories()
	return writeGob(f.PathData, f)
}

// Update sql filter.
func (f *HandytechSelectedCategories) UpdateSqlCategories() {
	categories := []string{}
	for _, category := range f.Categories {
		categories = append(categories, fmt.Sprintf("\"%s\"", category))
	}
	f.SqlCategories = strings.Join(categories, ", ")
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// SELECTED MAKERS
///////////////////////////////////////////////////////////////////////////////////////////////////
// Handytech Filters.
type HandytechSelectedMakers struct {
	Makers    []string
	SqlMakers string
	PathData  string
}

// Load handytech filters.
func LoadHandytechSelectedMakers(path string) *HandytechSelectedMakers {
	handytechSelectedMakers := HandytechSelectedMakers{}
	handytechSelectedMakers.PathData = path
	// Read handytech selected categories.
	err = readGob(handytechSelectedMakers.PathData, &handytechSelectedMakers)
	if err != nil {
		log.Printf("[warn] Not found Handytech selected makers data")
	}
	// log.Printf("Handytech selected makers: %v", handytechSelectedMakers)
	return &handytechSelectedMakers
}

// Save filters.
func (f *HandytechSelectedMakers) Save() error {
	f.UpdateSqlMakers()
	return writeGob(f.PathData, f)
}

// Update sql Makers.
func (f *HandytechSelectedMakers) UpdateSqlMakers() {
	makers := []string{}
	for _, maker := range f.Makers {
		makers = append(makers, fmt.Sprintf("\"%s\"", maker))
	}
	f.SqlMakers = strings.Join(makers, ", ")
}
