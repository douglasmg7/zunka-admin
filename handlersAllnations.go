package main

import (
	// "bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

///////////////////////////////////////////////////////////////////////////////////////////////////
// ALLNATIONS FILTERS
///////////////////////////////////////////////////////////////////////////////////////////////////
// Allnations Filters.
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
	log.Printf("Allnations filters: %v", allnationsFilters)
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
	log.Printf("Allnations selected categories: %v", allnationsSelectedCategories)
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
// FILTERS
///////////////////////////////////////////////////////////////////////////////////////////////////
// Get filters.
func allnationsFiltersHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
		Filters     AllnationsFilters
		Validation  AllnationsFiltersValidation
	}{session, "", *allnationsFilters, AllnationsFiltersValidation{"", ""}}

	// Render filters page.
	err = tmplAllnationsFilters.ExecuteTemplate(w, "allnationsFilters.tmpl", data)
	HandleError(w, err)
}

// Save filter.
func allnationsFiltersHandlerPost(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	filters := AllnationsFilters{}
	validation := AllnationsFiltersValidation{}

	// defer req.Body.Close()
	// body, err := ioutil.ReadAll(req.Body)
	// HandleError(w, err)
	// log.Printf("receive body: %v", string(body))

	if req.PostFormValue("onlyActive") != "" {
		filters.OnlyActive = true
		log.Printf("onlyActive: true")
	}

	if req.PostFormValue("onlyAvailability") != "" {
		filters.OnlyAvailability = true
		log.Printf("onlyAvailability: true")
	}

	// Validate min price.
	filters.MinPrice = req.PostFormValue("minPrice")
	log.Printf("filters.MinPrice: %s", filters.MinPrice)
	_, err = strconv.ParseUint(filters.MinPrice, 10, 64)
	if err != nil {
		log.Printf("err: %s", err)
		validation.MinPrice = "número inválido"
	}

	// Validate max price.
	filters.MaxPrice = req.PostFormValue("maxPrice")
	_, err = strconv.ParseUint(filters.MaxPrice, 10, 64)
	if err != nil {
		validation.MaxPrice = "número inválido"
	}

	// Some invalid fields.
	if validation.MinPrice != "" && validation.MaxPrice != "" {
		data := struct {
			Session     *SessionData
			HeadMessage string
			Filters     AllnationsFilters
			Validation  AllnationsFiltersValidation
		}{session, "", filters, validation}

		log.Printf("sending data: %v", data)

		// Render page.
		err = tmplAllnationsFilters.ExecuteTemplate(w, "allnationsFilters.tmpl", data)
		HandleError(w, err)
	} else {
		// Save filters and go to products page.
		allnationsFilters.OnlyActive = filters.OnlyActive
		allnationsFilters.OnlyAvailability = filters.OnlyAvailability
		allnationsFilters.MinPrice = filters.MinPrice
		allnationsFilters.MaxPrice = filters.MaxPrice
		err = allnationsFilters.Save()
		HandleError(w, err)
		http.Redirect(w, req, "/ns/allnations/products", http.StatusSeeOther)
	}

	// http.Redirect(w, req, "/ns/allnations/filters", http.StatusSeeOther)
	// w.WriteHeader(200)
	// w.Write([]byte("500 - Something bad happened!"))
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// PRODUCT
///////////////////////////////////////////////////////////////////////////////////////////////////
// Product list.
func allnationsProductsHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
		Products    []AllnationsProduct
		ValidDate   time.Time
	}{session, "", []AllnationsProduct{}, time.Now().Add(-VALID_DATE)}

	// Get products.
	log.Println(fmt.Sprintf(
		"SELECT * FROM product WHERE category IN (%s) AND %s ORDER BY description", allnationsSelectedCategories.SqlCategories, allnationsFilters.SqlFilter))
	err = dbAllnations.Select(&data.Products, fmt.Sprintf(
		"SELECT * FROM product WHERE category IN (%s) AND %s ORDER BY description", allnationsSelectedCategories.SqlCategories, allnationsFilters.SqlFilter))
	HandleError(w, err)

	err = tmplAllnationsProducts.ExecuteTemplate(w, "allnationsProducts.tmpl", data)
	HandleError(w, err)
}

// Product item.
func allnationsProductHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params, session *SessionData) {
	data := struct {
		Session                 *SessionData
		HeadMessage             string
		Product                 *AllnationsProduct
		TechnicalDescription    template.HTML
		ProductOld              *AllnationsProduct
		TechnicalDescriptionOld template.HTML
		RMAProcedureOld         template.HTML
		Status                  string
		ShowButtonCreateProduct bool
	}{session, "", &AllnationsProduct{}, "", &AllnationsProduct{}, "", "", "", false}

	// Get product.
	err = dbAllnations.Get(data.Product, "SELECT * FROM product WHERE code=?", ps.ByName("code"))
	HandleError(w, err)

	// Not escaped.
	data.TechnicalDescription = template.HTML(data.Product.TechnicalDescription)

	// Get product history.
	productsTemp := []AllnationsProduct{}
	err = dbAllnations.Select(&productsTemp, "SELECT * FROM product_history WHERE code=? AND changed_at < ? ORDER BY changed_at DESC LIMIT 1", ps.ByName("code"), data.Product.CheckedAt)
	HandleError(w, err)
	// If some history before checked_at.
	if len(productsTemp) > 0 {
		data.ProductOld = &productsTemp[0]
		fmt.Printf("Prodcut history: %s, Price: %v, ChangedAt: %v\n", productsTemp[0].Code, productsTemp[0].PriceSale, productsTemp[0].ChangedAt)
	} else {
		// Find the fist history.
		err = dbAllnations.Select(&productsTemp, "SELECT * FROM product_history WHERE code=? ORDER BY changed_at LIMIT 1", ps.ByName("code"))
		HandleError(w, err)
		if len(productsTemp) > 0 {
			data.ProductOld = &productsTemp[0]
			fmt.Println("first history")
			fmt.Printf("Prodcut history: %s, Price: %v, ChangedAt: %v\n", productsTemp[0].Code, productsTemp[0].PriceSale, productsTemp[0].ChangedAt)
		} else {
			// No history, poduct not changed.
			data.ProductOld = data.Product
			fmt.Println("No history")
		}
	}

	data.TechnicalDescriptionOld = template.HTML(data.ProductOld.TechnicalDescription)

	// Status.
	data.Status = data.Product.Status(time.Now().Add(-VALID_DATE))
	// Show button create product on zunkasite.
	if data.Product.ZunkaProductId == "" && (data.Status == "new" || data.Status == "changed" || data.Status == "") {
		data.ShowButtonCreateProduct = true
	}

	// Render template.
	err = tmplAllnationsProduct.ExecuteTemplate(w, "allnationsProduct.tmpl", data)
	HandleError(w, err)
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// CATEGORY
///////////////////////////////////////////////////////////////////////////////////////////////////
// Get categories.
func allnationsCategoriesHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
		Categories  []AllnationsCategory
	}{session, "", []AllnationsCategory{}}

	sql := fmt.Sprintf("SELECT category as name, count(category) as products_qty, false as selected FROM product "+
		"WHERE  %s GROUP BY category ORDER BY name", allnationsFilters.SqlFilter)
	err = dbAllnations.Select(&data.Categories, sql)
	HandleError(w, err)

	m := make(map[string]bool)
	for _, selectedCategory := range allnationsSelectedCategories.Categories {
		m[selectedCategory] = true
	}
	log.Printf("m: %v", m)
	for i := range data.Categories {
		if m[data.Categories[i].Name] {
			data.Categories[i].Selected = true
		}
	}

	// Render page.
	err = tmplAllnationsCategories.ExecuteTemplate(w, "allnationsCategories.tmpl", data)
	HandleError(w, err)
}

// Save categories.
func allnationsCategoriesHandlerPost(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {

	req.ParseForm()
	allnationsSelectedCategories.Categories = []string{}
	for key := range req.PostForm {
		allnationsSelectedCategories.Categories = append(allnationsSelectedCategories.Categories, key)
	}
	allnationsSelectedCategories.UpdateSqlCategories()
	err := allnationsSelectedCategories.Save()
	HandleError(w, err)
	http.Redirect(w, req, "/ns/allnations/products", http.StatusSeeOther)
}
