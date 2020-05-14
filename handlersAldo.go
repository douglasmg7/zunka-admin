package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/douglasmg7/aldoutil"
	"github.com/julienschmidt/httprouter"
)

// func isProductNew(product aldoutil.Product) bool {
// if product.ChangedAt.Equal(product.CreatedAt) {
// limitDate := time.Now().Add(time.Hour * 24 * 10)
// if product.ChangedAt.Before(limitDate) {
// return true
// }
// }
// return false
// }

// func isProductModified(product aldoutil.Product) bool {
// if !product.ChangedAt.Equal(product.CreatedAt) {
// limitDate := time.Now().Add(time.Hour * 24 * 10)
// if product.ChangedAt.Before(limitDate) {
// return true
// }
// }
// return false
// }

// Product list.
func aldoProductsHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
		Products    []aldoutil.Product
		ValidDate   time.Time
	}{session, "", []aldoutil.Product{}, time.Now().Add(-time.Hour * 240)}

	err = dbAldo.Select(&data.Products, "SELECT * FROM product order by description")
	// err = dbAldo.Select(&data.Products, "SELECT * FROM product order by description LIMIT 100 ")
	HandleError(w, err)

	// // Test - keep it commented.
	// i := 0
	// // Force new.
	// almostNow := time.Now().Add(-time.Hour * 48)
	// data.Products[i].CreatedAt = almostNow
	// data.Products[i].ChangedAt = almostNow
	// // Force changed.
	// // data.Products[i].ChangedAt = almostNow.Add(time.Hour * 24)
	// // Logs.
	// // log.Printf("Code: %v", data.Products[i].Code)
	// // log.Printf("CreatedAt: %v", data.Products[i].CreatedAt)
	// // log.Printf("ChangedAt: %v", data.Products[i].ChangedAt)
	// // log.Printf("ValidateDate: %v", data.ValidDate)
	// // Force unavailable.
	// // data.Products[i].Availability = false
	// // Force removed.
	// // data.Products[i].Removed = true
	// log.Printf("status: %v", data.Products[i].Status(data.ValidDate))
	// // End test.

	err = tmplAldoProducts.ExecuteTemplate(w, "aldoProducts.tmpl", data)
	HandleError(w, err)
}

// Product item.
func aldoProductHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params, session *SessionData) {
	data := struct {
		Session                 *SessionData
		HeadMessage             string
		Product                 *aldoutil.Product
		TechnicalDescription    template.HTML
		RMAProcedure            template.HTML
		Status                  string
		ShowButtonCreateProduct bool
	}{session, "", &aldoutil.Product{}, "", "", "", false}

	err = dbAldo.Get(data.Product, "SELECT * FROM product WHERE code=?", ps.ByName("code"))
	HandleError(w, err)

	// Not escaped.
	data.TechnicalDescription = template.HTML(data.Product.TechnicalDescription)
	data.RMAProcedure = template.HTML(data.Product.RMAProcedure)

	// // Test - keep commented.
	// // Force new.
	// almostNow := time.Now().Add(-time.Hour * 48)
	// data.Product.CreatedAt = almostNow
	// data.Product.ChangedAt = almostNow
	// // Force changed.
	// // data.Product.ChangedAt = almostNow.Add(time.Hour * 24)
	// // Logs.
	// // Force unavailable.
	// // data.Product.Availability = false
	// // Force removed.
	// // data.Product.Removed = true
	// // End test.

	// Status.
	data.Status = data.Product.Status(time.Now().Add(-time.Hour * 240))
	// Show button create product on zunkasite.
	if data.Product.MongodbId == "" && (data.Status == "new" || data.Status == "changed" || data.Status == "") {
		data.ShowButtonCreateProduct = true
	}
	// Render template.
	err = tmplAldoProduct.ExecuteTemplate(w, "aldoProduct.tmpl", data)
	HandleError(w, err)
}

// Add product to zunka site.
func aldoProductHandlerPost(w http.ResponseWriter, req *http.Request, ps httprouter.Params, session *SessionData) {
	// Get product.
	product := aldoutil.Product{}
	err = dbAldo.Get(&product, "SELECT * FROM product WHERE code=?", ps.ByName("code"))
	HandleError(w, err)

	// Set store product.
	storeProduct := aldoutil.StoreProduct{}
	storeProduct.DealerName = "Aldo"
	storeProduct.DealerProductId = product.Code

	// Title.
	// storeProduct.DealerProductTitle = strings.Title(strings.ToLower(product.Description))
	storeProduct.DealerProductTitle = product.Description

	// Category.
	storeProduct.DealerProductCategory = strings.Title(strings.ToLower(product.Category))

	// Maker.
	storeProduct.DealerProductMaker = strings.Title(strings.ToLower(product.Brand))

	// Description.
	storeProduct.DealerProductDesc = strings.TrimSpace(product.TechnicalDescription)
	// log.Println("product.TechnicalDescription:", product.TechnicalDescription)
	// log.Println("storeProduct.DealerProductDesc:", storeProduct.DealerProductDesc)

	// Image.
	storeProduct.DealerProductImagesLink = product.PictureLink

	// // Description.
	// product.TechnicalDescription = strings.TrimSpace(product.TechnicalDescription)

	// // Description.
	// // Split by <br/>
	// techDescs := regexp.MustCompile(`\<\s*br\s*\/*\>`).Split(product.TechnicalDescription, -1)
	// // Remove html tags.
	// // Remove blank intens.
	// reg := regexp.MustCompile(`\<[^>]*\>`)
	// buffer := bytes.Buffer{}
	// for _, val := range techDescs {
	// val = strings.TrimSpace(reg.ReplaceAllString(val, ""))
	// if val != "" {
	// buffer.WriteString(strings.ReplaceAll(val, ":", ";"))
	// buffer.WriteString("\n")
	// }
	// }
	// storeProduct.DealerProductDesc = strings.TrimSpace(buffer.String())

	// Months in days.
	storeProduct.DealerProductWarrantyDays = product.WarrantyPeriod * 30
	// Length in cm.
	storeProduct.DealerProductDeep = int(math.Ceil(float64(product.Length) / 10))
	// Width in cm.
	storeProduct.DealerProductWidth = int(math.Ceil(float64(product.Width) / 10))
	// Height in cm.
	storeProduct.DealerProductHeight = int(math.Ceil(float64(product.Height) / 10))
	// Weight in grams.
	storeProduct.DealerProductWeight = product.Weight
	// Price.
	storeProduct.DealerProductPrice = int(product.DealerPrice)
	// Suggestion price.
	storeProduct.DealerProductFinalPriceSuggestion = int(product.SuggestionPrice)
	// Last update.
	storeProduct.DealerProductLastUpdate = product.ChangedAt
	// Active.
	storeProduct.DealerProductActive = product.Availability

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
		HandleError(w, errors.New(fmt.Sprintf("Error ao solicitar a criação do produto no servidor zunka.\n\nstatus: %v\n\nbody: %v", res.StatusCode, string(resBody))))
		// HandleError(w, errors.New(fmt.Fprintf("Error ao solicitar a criação do produto no servidor zunka.\n\nStatus: %v\nError: %v", res.StatusCode, resBody)))
		return
	}
	// Mongodb id from created product.
	product.MongodbId = string(resBody)
	// Remove suround double quotes.
	product.MongodbId = product.MongodbId[1 : len(product.MongodbId)-1]

	// Update product with _id from mongodb store.
	// stmt, err := dbAldo.Prepare(`UPDATE product SET mongodb_id = $1 WHERE id = $2;`)
	stmt, err := dbAldo.Prepare(`UPDATE product SET mongodb_id = $1 WHERE code = $2;`)
	HandleError(w, err)
	defer stmt.Close()
	_, err = stmt.Exec(product.MongodbId, product.Code)
	HandleError(w, err)

	// Render categories page.
	http.Redirect(w, req, "/ns/aldo/product/"+product.Code, http.StatusSeeOther)
}

// Remove mongodb id from Product.
func aldoProductMongodbIdHandlerDelete(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// Update mongodb store id.
	stmt, err := dbAldo.Prepare(`UPDATE product SET mongodb_id = $1 WHERE code = $2;`)
	HandleError(w, err)
	defer stmt.Close()
	// fmt.Println("code:", ps.ByName("code"))
	_, err = stmt.Exec("", ps.ByName("code"))
	// _, err = stmt.Exec(ps.ByName("code"), "")
	HandleError(w, err)
	// w.Write([]byte("OK"))
	w.WriteHeader(200)
}

// Categories.
func aldoCategoriesHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
		Categories  []aldoutil.Category
	}{session, "", []aldoutil.Category{}}

	// Get categories from db.
	err = dbAldo.Select(&data.Categories, "SELECT * FROM category order by name")
	HandleError(w, err)

	// Render page.
	err = tmplAldoCategories.ExecuteTemplate(w, "aldoCategories.tmpl", data)
	HandleError(w, err)
}

// Save categories.
func aldoCategoriesHandlerPost(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	// Get categories from db.
	categories := []aldoutil.Category{}
	err = dbAldo.Select(&categories, "SELECT * FROM category order by name")
	HandleError(w, err)
	// Prepare update.
	stmt, err := dbAldo.Prepare(`UPDATE category SET selected = $1 WHERE name = $2;`)
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
	http.Redirect(w, req, "/ns/aldo/categories", http.StatusSeeOther)

	// Run script to process xml products.
	// log.Println(GS + "/aldowsc/bin/process-xml-products.sh")
	cmd := exec.Command(GS + "/aldowsc/bin/process-xml-products.sh")
	err = cmd.Start()
	if err != nil {
		log.Printf("Error running script to process XML Aldo products. %s", err)
	}
	// log.Printf("Waiting for command to finish...")
	// err = cmd.Wait()
	// log.Printf("Command finished with error: %v", err)

	return
}
