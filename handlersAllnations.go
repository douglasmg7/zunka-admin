package main

import (
	// "bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

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
	// log.Println(fmt.Sprintf(
	// "SELECT * FROM product WHERE category IN (%s) AND %s ORDER BY description", allnationsSelectedCategories.SqlCategories, allnationsFilters.SqlFilter))
	err = dbAllnations.Select(&data.Products, fmt.Sprintf(
		"SELECT * FROM product WHERE category IN (%s) AND maker IN (%s) AND %s ORDER BY description",
		allnationsSelectedCategories.SqlCategories, allnationsSelectedMakers.SqlMakers, allnationsFilters.SqlFilter))
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
	// log.Printf("m: %v", m)
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

///////////////////////////////////////////////////////////////////////////////////////////////////
// MAKERS
///////////////////////////////////////////////////////////////////////////////////////////////////
// Get categories.
func allnationsMakersHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
		Makers      []AllnationsMaker
	}{session, "", []AllnationsMaker{}}

	sql := fmt.Sprintf("SELECT maker as name, count(maker) as products_qty, false as selected FROM product "+
		"WHERE  %s GROUP BY maker ORDER BY name", allnationsFilters.SqlFilter)
	err = dbAllnations.Select(&data.Makers, sql)
	HandleError(w, err)

	m := make(map[string]bool)
	for _, selectedMaker := range allnationsSelectedMakers.Makers {
		m[selectedMaker] = true
	}
	// log.Printf("m: %v", m)
	for i := range data.Makers {
		if m[data.Makers[i].Name] {
			data.Makers[i].Selected = true
		}
	}

	// Render page.
	err = tmplAllnationsMakers.ExecuteTemplate(w, "allnationsMakers.gohtml", data)
	HandleError(w, err)
}

// Save categories.
func allnationsMakersHandlerPost(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {

	req.ParseForm()
	allnationsSelectedMakers.Makers = []string{}
	for key := range req.PostForm {
		allnationsSelectedMakers.Makers = append(allnationsSelectedMakers.Makers, key)
	}
	allnationsSelectedMakers.UpdateSqlMakers()
	err := allnationsSelectedMakers.Save()
	HandleError(w, err)
	http.Redirect(w, req, "/ns/allnations/products", http.StatusSeeOther)
}
