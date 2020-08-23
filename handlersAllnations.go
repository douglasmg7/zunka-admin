package main

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	// "os/exec"
)

// Write struct to file.
func writeGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

// Read struct from file.
func readGob(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}

// Allnations Filters.
type AllnationsFilters struct {
	OnlyActive       bool
	OnlyAvailability bool
	MinPrice         int
	MaxPrice         int
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// FILTERS
///////////////////////////////////////////////////////////////////////////////////////////////////
// Get all filters.
func allnationsFiltersHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	// todo - maker
	// Default filters.
	filters := AllnationsFilters{
		OnlyActive:       false,
		OnlyAvailability: false,
		MinPrice:         100_00,
		MaxPrice:         1000_000_00,
	}
	// Read filter from file.
	err := readGob("filters.data", filters)
	if err != nil {
		log.Printf("Using default values for Allnations filters: %v", filters)
	} else {
		log.Printf("Using Allnations filters: %v", filters)
	}

	data := struct {
		Session     *SessionData
		HeadMessage string
		Filters     AllnationsFilters
	}{session, "", filters}

	// Render page.
	err = tmplAllnationsFilters.ExecuteTemplate(w, "allnationsFilters.tmpl", data)
	HandleError(w, err)
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

	// Get selected categories.
	categories := []AllnationsCategory{}
	err = dbAllnations.Select(&categories, "SELECT * FROM category WHERE selected = true")
	categoriesSlice := []string{}
	for _, category := range categories {
		categoriesSlice = append(categoriesSlice, fmt.Sprintf("\"%v\"", category.Name))
	}
	categoriesList := strings.Join(categoriesSlice, ", ")
	log.Printf("Categories: %v", categoriesList)

	// Get products.
	err = dbAllnations.Select(&data.Products, fmt.Sprintf(
		"SELECT * FROM product WHERE category IN (%v) AND active = true AND availability = true  ORDER BY description", categoriesList))
	HandleError(w, err)

	err = tmplAllnationsProducts.ExecuteTemplate(w, "allnationsProducts.tmpl", data)
	HandleError(w, err)
}

// // Product list.
// func allnationsProductsHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
// data := struct {
// Session     *SessionData
// HeadMessage string
// Products    []AllnationsProduct
// ValidDate   time.Time
// }{session, "", []AllnationsProduct{}, time.Now().Add(-VALID_DATE)}

// // Get selected categories.
// categories := []AllnationsCategory{}
// err = dbAllnations.Select(&categories, "SELECT * FROM category WHERE selected = true")
// categoriesSlice := []string{}
// for _, category := range categories {
// categoriesSlice = append(categoriesSlice, fmt.Sprintf("\"%v\"", category.Name))
// }
// categoriesList := strings.Join(categoriesSlice, ", ")
// log.Printf("Categories: %v", categoriesList)

// // Get products.
// err = dbAllnations.Select(&data.Products, fmt.Sprintf("SELECT * FROM product WHERE category IN (%v) ORDER BY description", categoriesList))
// HandleError(w, err)

// err = tmplAllnationsProducts.ExecuteTemplate(w, "allnationsProducts.tmpl", data)
// HandleError(w, err)
// }

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
// Get all categories.
func allnationsCategoriesHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
		Categories  []AllnationsCategory
	}{session, "", []AllnationsCategory{}}

	// Get categories from db.
	// select category, count(category) from product group by category
	// err = dbAllnations.Select(&data.Categories, "SELECT * FROM category order by name")
	sql := "SELECT category as name, count(category) as products_qty, true as selected FROM product " +
		"WHERE category IN (\"ARMAZENAMENTO\", \"SCANNER\") AND active = true AND availability = true AND price_sale > 323000 " +
		"GROUP BY category " +
		"ORDER BY name"
	err = dbAllnations.Select(&data.Categories, sql)
	HandleError(w, err)

	// Render page.
	err = tmplAllnationsCategories.ExecuteTemplate(w, "allnationsCategories.tmpl", data)
	HandleError(w, err)
}

// // Get all categories.
// func allnationsCategoriesHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
// data := struct {
// Session     *SessionData
// HeadMessage string
// Categories  []AllnationsCategory
// }{session, "", []AllnationsCategory{}}

// // Get categories from db.
// err = dbAllnations.Select(&data.Categories, "SELECT * FROM category order by name")
// HandleError(w, err)

// // Render page.
// err = tmplAllnationsCategories.ExecuteTemplate(w, "allnationsCategories.tmpl", data)
// HandleError(w, err)
// }

// Save categories.
func allnationsCategoriesHandlerPost(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	// Get categories from db.
	categories := []AllnationsCategory{}
	err = dbAllnations.Select(&categories, "SELECT * FROM category order by name")
	HandleError(w, err)
	// Prepare update.
	stmt, err := dbAllnations.Prepare(`UPDATE category SET selected = $1 WHERE name = $2;`)
	HandleError(w, err)
	defer stmt.Close()

	// For each category on db.
	for _, category := range categories {
		// Update changed categories.
		if (category.Selected && (req.PostFormValue(category.Name) == "")) || (!category.Selected && (req.PostFormValue(category.Name) != "")) {
			// fmt.Println("Updated category:", category.Name)
			_, err = stmt.Exec(!category.Selected, category.Name)
			HandleError(w, err)
		}
	}
	// Render categories page.
	http.Redirect(w, req, "/ns/allnations/categories", http.StatusSeeOther)

	// todo - Run script to process xml products.
	// cmd := exec.Command(GS + "/aldowsc/bin/process-xml-products.sh")
	// err = cmd.Start()
	// if err != nil {
	// log.Printf("Error running script to process XML Allnations products. %s", err)
	// }

	return
}
