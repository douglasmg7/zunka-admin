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
// Aldo product.
type HandytechProduct struct {
	ZunkaProductId         sql.NullString `db:"zunka_product_id"`
	ItCodigo               sql.NullString `db:"it_codigo"` // From dealer.
	DescItem               sql.NullString `db:"desc_item"`
	DescItemEc             sql.NullString `db:"desc_item_ec"`
	NarrativaEc            sql.NullString `db:"narrativa_ec"`
	VlItem                 sql.NullInt64  `db:"vl_item"`
	VlItemSdesc            sql.NullInt64  `db:"vl_item_sdesc"`
	VlIpi                  sql.NullInt64  `db:"vl_ipi"`
	PercPrecoSugeridoSolar sql.NullInt64  `db:"perc_preco_sugerido_solar"`
	PrecoSugerido          sql.NullInt64  `db:"preco_sugerido"`
	PrecoMaximo            sql.NullInt64  `db:"preco_maximo"`
	Categoria              sql.NullString `db:"categoria"`
	SubCategoria           sql.NullString `db:"sub_categoria"`
	Peso                   sql.NullInt64  `db:"peso"`
	CodigoRefer            sql.NullString `db:"codigo_refer"`
	Fabricante             sql.NullString `db:"fabricante"`
	Saldos                 sql.NullInt64  `db:"saldos"`
	ArquivoImagem          sql.NullString `db:"arquivo_imagem"`
	ImagesUrl              []string
	CreatedAt              time.Time `db:"created_at"`
	// ChangedAt              time.Time      `db:"changed_at"`
}

// Process Br Currency.
func (p *HandytechProduct) ProcessArquivoImagem() {
	if p.ArquivoImagem.Valid {
		p.ImagesUrl = strings.Split(p.ArquivoImagem.String, "\uffff")
		for _, s := range p.ImagesUrl {
			log.Println(s)
		}
	}
}

// Define product status.
func (p *HandytechProduct) Status() string {
	if p.ZunkaProductId.String == "" {
		return "no-registered"
	}
	return "registered"
}

// Process Br Currency.
func (p *HandytechProduct) ProcessBrCurrency(val sql.NullInt64) string {
	if val.Valid {
		printer := message.NewPrinter(language.Portuguese)
		return printer.Sprintf("R$ %.2f", float64(val.Int64)/100)
	} else {
		return "NULL"
	}
}

// Process Br Currency.
func (p *HandytechProduct) ProcessWight(val sql.NullInt64) string {
	if val.Valid {
		printer := message.NewPrinter(language.Portuguese)
		return printer.Sprintf("%.3f kg", float64(val.Int64)/1000)
	} else {
		return "NULL"
	}
}

// Process string, show NULL.
func (p *HandytechProduct) ProcessString(val sql.NullString) string {
	if val.Valid {
		return val.String
	} else {
		return "NULL"
	}
}

// Check if some valid url image.
func (p *HandytechProduct) HasUrlImage() bool {
	if p.ArquivoImagem.Valid {
		if len(p.ArquivoImagem.String) > 0 {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
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
	MinPrice  string
	MaxPrice  string
	PathData  string
	SqlFilter string
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
