package main

import (
	// "bytes"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/douglasmg7/aldoutil"
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

// Add product to zunka site.
func allnationsProductHandlerPost(w http.ResponseWriter, req *http.Request, ps httprouter.Params, session *SessionData) {
	// Get product.
	product := AllnationsProduct{}
	// Get product.
	err = dbAllnations.Get(&product, "SELECT * FROM product WHERE code=?", ps.ByName("code"))
	HandleError(w, err)

	// Set store product.
	storeProduct := aldoutil.StoreProduct{}
	storeProduct.DealerName = "Allnations"
	storeProduct.DealerProductId = product.Code

	// Title.
	// storeProduct.DealerProductTitle = strings.Title(strings.ToLower(product.Description))
	storeProduct.DealerProductTitle = product.Description
	// Category.
	storeProduct.DealerProductCategory = strings.Title(strings.ToLower(product.Category))
	// Maker.
	storeProduct.DealerProductMaker = strings.Title(strings.ToLower(product.Maker))
	// Description.
	storeProduct.DealerProductDesc = strings.TrimSpace(product.TechnicalDescription)
	// log.Println("product.TechnicalDescription:", product.TechnicalDescription)
	// log.Println("storeProduct.DealerProductDesc:", storeProduct.DealerProductDesc)
	// Image.
	storeProduct.DealerProductImagesLink = product.UrlImage

	// Months in days.
	storeProduct.DealerProductWarrantyDays = product.WarrantyMonth * 30
	// Length in cm.
	storeProduct.DealerProductDeep = int(math.Ceil(float64(product.LengthMm) / 10))
	// Width in cm.
	storeProduct.DealerProductWidth = int(math.Ceil(float64(product.WidthMm) / 10))
	// Height in cm.
	storeProduct.DealerProductHeight = int(math.Ceil(float64(product.HeightMm) / 10))
	// Weight in grams.
	storeProduct.DealerProductWeight = product.WeightG
	// Price.
	storeProduct.DealerProductPrice = int(product.PriceSale)
	// Suggestion price.
	storeProduct.DealerProductFinalPriceSuggestion = int(product.PriceSale)
	// Last update.
	storeProduct.DealerProductLastUpdate = product.ChangedAt
	// Active.
	storeProduct.DealerProductActive = product.Availability && product.Active
	// Stock.
	storeProduct.StoreProductQtd = product.StockQty

	// Convert to json.
	reqBody, err := json.Marshal(storeProduct)
	HandleError(w, err)

	// Log request.
	// log.Println("request body: " + string(reqBody))

	// Request product add.
	client := &http.Client{}
	req, err = http.NewRequest("POST", zunkaSiteHost()+"/setup/product/add", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	HandleError(w, err)
	req.SetBasicAuth(zunkaSiteUser(), zunkaSitePass())
	res, err := client.Do(req)
	HandleError(w, err)

	// res, err := http.Post("http://localhost:3080/setup/product/add", "application/json", bytes.NewBuffer(reqBody))
	defer res.Body.Close()
	HandleError(w, err)

	// Result.
	resBody, err := ioutil.ReadAll(res.Body)
	HandleError(w, err)
	// No 200 status.
	if res.StatusCode != 200 {
		HandleError(w, errors.New(fmt.Sprintf("Error ao solicitar a criação do produto allnations no servidor zunka.\n\nstatus: %v\n\nbody: %v", res.StatusCode, string(resBody))))
		return
	}
	// Mongodb id from created product.
	product.ZunkaProductId = string(resBody)
	// Remove suround double quotes.
	product.ZunkaProductId = product.ZunkaProductId[1 : len(product.ZunkaProductId)-1]

	// Update product with _id from mongodb store and set checked_at.
	stmt, err := dbAllnations.Prepare(`UPDATE product SET zunka_product_id = $1, checked_at=$2 WHERE code = $3;`)
	HandleError(w, err)
	defer stmt.Close()
	_, err = stmt.Exec(product.ZunkaProductId, time.Now(), product.Code)
	HandleError(w, err)

	// Render product page.
	http.Redirect(w, req, "/ns/allnations/product/"+product.Code, http.StatusSeeOther)
}

// Aldo product checked.
func allnationsProductCheckedHandlerPost(w http.ResponseWriter, req *http.Request, ps httprouter.Params, session *SessionData) {
	productCode := ps.ByName("code")
	// Update product status_cleaned_at field.
	stmt, err := dbAllnations.Prepare(`UPDATE product SET checked_at=$1 WHERE code = $2;`)
	HandleError(w, err)
	defer stmt.Close()
	_, err = stmt.Exec(time.Now(), productCode)
	HandleError(w, err)

	// Render categories page.
	http.Redirect(w, req, "/ns/allnations/product/"+ps.ByName("code"), http.StatusSeeOther)
}

// Remove mongodb id from Product.
func allnationsProductZunkaProductIdHandlerDelete(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// Update mongodb store id.
	stmt, err := dbAllnations.Prepare(`UPDATE product SET zunka_product_id = $1 WHERE code = $2;`)
	HandleError(w, err)
	defer stmt.Close()
	// fmt.Println("code:", ps.ByName("code"))
	_, err = stmt.Exec("", ps.ByName("code"))
	// _, err = stmt.Exec(ps.ByName("code"), "")
	HandleError(w, err)
	// w.Write([]byte("OK"))
	w.WriteHeader(200)
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
