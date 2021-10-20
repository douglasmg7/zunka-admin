package main

import (
	// "bytes"

	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	// "math"
	"net/http"
	"strconv"

	// "strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

///////////////////////////////////////////////////////////////////////////////////////////////////
// FILTERS
///////////////////////////////////////////////////////////////////////////////////////////////////
// Get filters.
func motospeedFiltersHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
		Filters     MotospeedFilters
		Validation  MotospeedFiltersValidation
	}{session, "", *motospeedFilters, MotospeedFiltersValidation{"", ""}}

	// Render filters page.
	err = tmplMotospeedFilters.ExecuteTemplate(w, "motospeed_filters.gohtml", data)
	HandleError(w, err)
}

// Save filter.
func motospeedFiltersHandlerPost(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	filters := MotospeedFilters{}
	validation := MotospeedFiltersValidation{}

	// defer req.Body.Close()
	// body, err := ioutil.ReadAll(req.Body)
	// HandleError(w, err)
	// log.Printf("receive body: %v", string(body))

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
			Filters     MotospeedFilters
			Validation  MotospeedFiltersValidation
		}{session, "", filters, validation}

		log.Printf("sending data: %v", data)

		// Render page.
		err = tmplMotospeedFilters.ExecuteTemplate(w, "motospeed_filters.gohtml", data)
		HandleError(w, err)
	} else {
		// Save filters and go to products page.
		motospeedFilters.MinPrice = filters.MinPrice
		motospeedFilters.MaxPrice = filters.MaxPrice
		err = motospeedFilters.Save()
		HandleError(w, err)
		http.Redirect(w, req, "/ns/motospeed/products", http.StatusSeeOther)
	}

	// http.Redirect(w, req, "/ns/motospeed/filters", http.StatusSeeOther)
	// w.WriteHeader(200)
	// w.Write([]byte("500 - Something bad happened!"))
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// PRODUCT
///////////////////////////////////////////////////////////////////////////////////////////////////
// Product list.
func motospeedProductsHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
		Products    []MotospeedProduct
		ValidDate   time.Time
	}{session, "", []MotospeedProduct{}, time.Now().Add(-VALID_DATE)}

	// Get products.
	// log.Println(fmt.Sprintf(
	// "SELECT * FROM product WHERE category IN (%s) AND %s ORDER BY description", motospeedSelectedCategories.SqlCategories, MotospeedFilters.SqlFilter))

	// With sql filter.
	// err = dbMotospeed.Select(&data.Products, fmt.Sprintf(
	// "SELECT * FROM product WHERE categoria IN (%s) AND fabricante IN (%s) AND %s ORDER BY desc",
	// motospeedSelectedCategories.SqlCategories, motospeedSelectedMakers.SqlMakers, MotospeedFilters.SqlFilter))

	// err = dbMotospeed.Select(&data.Products, fmt.Sprintf(
	// "SELECT * FROM product WHERE categoria IN (%s) ORDER BY desc",
	// motospeedSelectedCategories.SqlCategories))

	err = dbMotospeed.Select(&data.Products, fmt.Sprintf("SELECT * FROM product ORDER BY desc"))

	HandleError(w, err)

	log.Println("Products count:", len(data.Products))

	err = tmplMotospeedProducts.ExecuteTemplate(w, "motospeed_products.gohtml", data)
	HandleError(w, err)
}

// Product item.
func motospeedProductHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params, session *SessionData) {
	data := struct {
		Session                 *SessionData
		HeadMessage             string
		Product                 *MotospeedProduct
		Status                  string
		ShowButtonCreateProduct bool
		SimiliarZunkaProducts   []ZunkaSiteProductRx
		SameEANZunkaProducts    []ZunkaSiteProductRx
	}{session, "", &MotospeedProduct{}, "", false, []ZunkaSiteProductRx{}, []ZunkaSiteProductRx{}}

	// Get product.
	err = dbMotospeed.Get(data.Product, "SELECT * FROM product WHERE code=?", ps.ByName("code"))
	HandleError(w, err)

	// Status.
	data.Status = data.Product.Status()

	// Show option to create product on zunkasite.
	if !data.Product.ZunkaProductId.Valid || data.Product.ZunkaProductId.String == "" {
		data.ShowButtonCreateProduct = true

		// Similar titles.
		chanProductsSimilarTitles := make(chan []ZunkaSiteProductRx)
		go getProductsSimilarTitles(chanProductsSimilarTitles, data.Product.Desc.String)

		// Same EAN.
		// todo - search product data for ean.
		ean := ""
		if len(ean) > 0 {
			chanProductsSameEAN := make(chan []ZunkaSiteProductRx)
			go getProductsSameEAN(chanProductsSameEAN, ean)
			data.SimiliarZunkaProducts, data.SameEANZunkaProducts = <-chanProductsSimilarTitles, <-chanProductsSameEAN
		} else {
			data.SimiliarZunkaProducts = <-chanProductsSimilarTitles
		}
		// Debug.Printf("SimiliarZunkaPorudcts: %v", data.SimiliarZunkaProducts)
	}

	// Render template.
	err = tmplMotospeedProduct.ExecuteTemplate(w, "motospeed_product.gohtml", data)
	HandleError(w, err)
}

// Add product to zunka site.
func motospeedProductHandlerPost(w http.ResponseWriter, req *http.Request, ps httprouter.Params, session *SessionData) {
	req.ParseForm()
	HandleError(w, err)

	// Get product.
	product := MotospeedProduct{}
	// Get product.
	err = dbMotospeed.Get(&product, "SELECT * FROM product WHERE code=?", ps.ByName("code"))
	HandleError(w, err)

	// Set store product.
	// storeProduct := aldoutil.StoreProduct{}
	storeProduct := ZunkaSiteProductTx{}

	storeProduct.ProductIdTemplate = req.FormValue("similar-product")
	storeProduct.DealerName = "Motospeed"
	// storeProduct.DealerProductId = product.Code

	// // Title.
	// // storeProduct.DealerProductTitle = strings.Title(strings.ToLower(product.Description))
	// storeProduct.DealerProductTitle = product.Description
	// // Category.
	// storeProduct.DealerProductCategory = strings.Title(strings.ToLower(product.Category))
	// // Maker.
	// storeProduct.DealerProductMaker = strings.Title(strings.ToLower(product.Maker))
	// // Description.
	// storeProduct.DealerProductDesc = strings.TrimSpace(product.TechnicalDescription)
	// // log.Println("product.TechnicalDescription:", product.TechnicalDescription)
	// // log.Println("storeProduct.DealerProductDesc:", storeProduct.DealerProductDesc)
	// // Image.
	// storeProduct.DealerProductImagesLink = product.UrlImage

	// // Months in days.
	// storeProduct.DealerProductWarrantyDays = product.WarrantyMonth * 30
	// // Length in cm.
	// storeProduct.DealerProductDeep = int(math.Ceil(float64(product.LengthMm) / 10))
	// // Width in cm.
	// storeProduct.DealerProductWidth = int(math.Ceil(float64(product.WidthMm) / 10))
	// // Height in cm.
	// storeProduct.DealerProductHeight = int(math.Ceil(float64(product.HeightMm) / 10))
	// // Weight in grams.
	// storeProduct.DealerProductWeight = product.WeightG
	// // Price.
	// storeProduct.DealerProductPrice = int(product.PriceSale)
	// // Suggestion price.
	// storeProduct.DealerProductFinalPriceSuggestion = int(product.PriceSale + product.PriceSale/3)
	// // Last update.
	// storeProduct.DealerProductLastUpdate = product.ChangedAt
	// // Active.
	// storeProduct.DealerProductActive = product.Availability && product.Active
	// // Origin.
	// storeProduct.DealerProductLocation = product.StockOrigin
	// // Stock.
	// storeProduct.StoreProductQtd = product.StockQty
	// // Ean.
	// storeProduct.Ean = product.Ean

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
		HandleError(w, errors.New(fmt.Sprintf("Error ao solicitar a criação do produto motospeed no servidor zunka.\n\nstatus: %v\n\nbody: %v", res.StatusCode, string(resBody))))
		return
	}
	// Mongodb id from created product.
	product.ZunkaProductId.String = string(resBody)
	// Remove suround double quotes.
	product.ZunkaProductId.String = product.ZunkaProductId.String[1 : len(product.ZunkaProductId.String)-1]

	// Update product with _id from mongodb store and set checked_at.
	stmt, err := dbMotospeed.Prepare(`UPDATE product SET zunka_product_id = $1, checked_at=$2 WHERE code = $3;`)
	HandleError(w, err)
	defer stmt.Close()
	_, err = stmt.Exec(product.ZunkaProductId, time.Now(), product.Code)
	HandleError(w, err)

	// Render product page.
	// http.Redirect(w, req, "/ns/motospeed/product/"+product.Code, http.StatusSeeOther)

	// Back to product list.
	http.Redirect(w, req, "/ns/motospeed/products", http.StatusSeeOther)
}

// Product checked.
func motospeedProductCheckedHandlerPost(w http.ResponseWriter, req *http.Request, ps httprouter.Params, session *SessionData) {
	productCode := ps.ByName("code")
	// Update product status_cleaned_at field.
	stmt, err := dbMotospeed.Prepare(`UPDATE product SET checked_at=$1 WHERE code = $2;`)
	HandleError(w, err)
	defer stmt.Close()
	_, err = stmt.Exec(time.Now(), productCode)
	HandleError(w, err)

	// Render categories page.
	http.Redirect(w, req, "/ns/motospeed/product/"+ps.ByName("code"), http.StatusSeeOther)
}

// Remove mongodb id from Product.
func motospeedProductZunkaProductIdHandlerDelete(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// Update mongodb store id.
	stmt, err := dbMotospeed.Prepare(`UPDATE product SET zunka_product_id = $1 WHERE code = $2;`)
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
func motospeedCategoriesHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
		Categories  []MotospeedCategory
	}{session, "", []MotospeedCategory{}}

	// sql := fmt.Sprintf("SELECT categoria as name, count(categoria) as products_qty, false as selected FROM product "+
	// "WHERE  %s GROUP BY categoria ORDER BY name", MotospeedFilters.SqlFilter)
	sql := fmt.Sprintf("SELECT categoria as name, count(categoria) as products_qty, false as selected FROM product " +
		"GROUP BY categoria ORDER BY name")
	// log.Printf("sql: %v", sql)
	err = dbMotospeed.Select(&data.Categories, sql)
	HandleError(w, err)

	m := make(map[string]bool)
	for _, selectedCategory := range motospeedSelectedCategories.Categories {
		m[selectedCategory] = true
	}
	// log.Printf("m: %v", m)
	for i := range data.Categories {
		if m[data.Categories[i].Name] {
			data.Categories[i].Selected = true
		}
	}

	// log.Printf("data: %v", len(data.Categories))
	// Render page.
	err = tmplMotospeedCategories.ExecuteTemplate(w, "motospeed_categories.gohtml", data)
	HandleError(w, err)
}

// Save categories.
func motospeedCategoriesHandlerPost(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {

	req.ParseForm()
	motospeedSelectedCategories.Categories = []string{}
	for key := range req.PostForm {
		motospeedSelectedCategories.Categories = append(motospeedSelectedCategories.Categories, key)
	}
	motospeedSelectedCategories.UpdateSqlCategories()
	err := motospeedSelectedCategories.Save()
	HandleError(w, err)
	http.Redirect(w, req, "/ns/motospeed/products", http.StatusSeeOther)
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// CSV
///////////////////////////////////////////////////////////////////////////////////////////////////
// Load csv page.
func motospeedLoadCSVHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
	}{session, ""}

	// Render page.
	err = tmplMotospeedLoadCSV.ExecuteTemplate(w, "motospeed_load_csv.gohtml", data)
	HandleError(w, err)
}

// Load csv.
func motospeedLoadCSVHandlerPost(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	// log.Println("File Upload Endpoint Hit")
	// log.Println(req.Header)

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	req.ParseMultipartForm(10 << 20)

	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, header, err := req.FormFile("csv-file")
	if err != nil {
		HandleError(w, err)
		return
	}
	defer file.Close()
	log.Printf("Motorspeed csv uploaded, file: %+v, size: %+v", header.Filename, header.Size)
	// fmt.Printf("Uploaded File: %+v\n", header.Filename)
	// fmt.Printf("File Size: %+v\n", header.Size)
	// fmt.Printf("MIME Header: %+v\n", header.Header)

	tmpfile, _ := os.Create(MOTOSPEED_CSV)
	// tmpfile, _ := os.Create(path.Join(motospeedDataPath, "/motospeed_products.csv"))
	io.Copy(tmpfile, file)
	tmpfile.Close()
	file.Close()

	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Arquivo CSV carregado com sucesso.")
}
