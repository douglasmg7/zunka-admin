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

	"github.com/douglasmg7/aldoutil"
	"github.com/julienschmidt/httprouter"
)

// Product list.
func aldoProductsHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
		Products    []aldoutil.Product
	}{session, "", []aldoutil.Product{}}

	err = dbAldo.Select(&data.Products, "SELECT * FROM product order by description LIMIT 100 ")
	HandleError(w, err)

	err = tmplAldoProducts.ExecuteTemplate(w, "aldoProducts.tmpl", data)
	HandleError(w, err)
}

// Product item.
func aldoProductHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params, session *SessionData) {
	data := struct {
		Session              *SessionData
		HeadMessage          string
		Product              aldoutil.Product
		TechnicalDescription template.HTML
		RMAProcedure         template.HTML
	}{session, "", aldoutil.Product{}, "", ""}

	err = dbAldo.Get(&data.Product, "SELECT * FROM product WHERE code=?", ps.ByName("code"))
	HandleError(w, err)

	// Not escaped.
	data.TechnicalDescription = template.HTML(data.Product.TechnicalDescription)
	data.RMAProcedure = template.HTML(data.Product.RMAProcedure)

	err = tmplAldoProduct.ExecuteTemplate(w, "aldoProduct.tmpl", data)
	HandleError(w, err)
}

// Add product to zunka site.
func aldoProductHandlerPost(w http.ResponseWriter, req *http.Request, ps httprouter.Params, session *SessionData) {
	data := struct {
		Session              *SessionData
		HeadMessage          string
		Product              aldoutil.Product
		TechnicalDescription template.HTML
		RMAProcedure         template.HTML
	}{session, "", aldoutil.Product{}, "", ""}

	// Get product.
	err = dbAldo.Get(&data.Product, "SELECT * FROM product WHERE code=?", ps.ByName("code"))
	HandleError(w, err)

	// Set store product.
	storeProduct := aldoutil.StoreProduct{}
	storeProduct.DealerName = "Aldo"
	storeProduct.DealerProductId = data.Product.Code

	// Title.
	// storeProduct.DealerProductTitle = strings.Title(strings.ToLower(data.Product.Description))
	storeProduct.DealerProductTitle = data.Product.Description

	// Category.
	storeProduct.DealerProductCategory = strings.Title(strings.ToLower(data.Product.Category))

	// Maker.
	storeProduct.DealerProductMaker = strings.Title(strings.ToLower(data.Product.Brand))

	// Description.
	storeProduct.DealerProductDesc = strings.TrimSpace(data.Product.TechnicalDescription)
	// log.Println("data.Product.TechnicalDescription:", data.Product.TechnicalDescription)
	// log.Println("storeProduct.DealerProductDesc:", storeProduct.DealerProductDesc)

	// // Description.
	// data.Product.TechnicalDescription = strings.TrimSpace(data.Product.TechnicalDescription)

	// // Description.
	// // Split by <br/>
	// techDescs := regexp.MustCompile(`\<\s*br\s*\/*\>`).Split(data.Product.TechnicalDescription, -1)
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
	storeProduct.DealerProductWarrantyDays = data.Product.WarrantyPeriod * 30
	// Length in cm.
	storeProduct.DealerProductDeep = int(math.Ceil(float64(data.Product.Length) / 10))
	// Width in cm.
	storeProduct.DealerProductWidth = int(math.Ceil(float64(data.Product.Width) / 10))
	// Height in cm.
	storeProduct.DealerProductHeight = int(math.Ceil(float64(data.Product.Height) / 10))
	// Weight in grams.
	storeProduct.DealerProductWeight = data.Product.Weight
	// Price.
	storeProduct.DealerProductPrice = int(data.Product.DealerPrice)
	// Suggestion price.
	storeProduct.DealerProductFinalPriceSuggestion = int(data.Product.SuggestionPrice)
	// Last update.
	storeProduct.DealerProductLastUpdate = data.Product.ChangedAt
	// Active.
	storeProduct.DealerProductActive = data.Product.Availability

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
	data.Product.MongodbId = string(resBody)
	// Remove suround double quotes.
	data.Product.MongodbId = data.Product.MongodbId[1 : len(data.Product.MongodbId)-1]

	// Update product with _id from mongodb store.
	stmt, err := dbAldo.Prepare(`UPDATE product SET mongodb_id = $1 WHERE id = $2;`)
	HandleError(w, err)
	defer stmt.Close()
	_, err = stmt.Exec(data.Product.MongodbId, data.Product.Id)
	HandleError(w, err)

	// Not escaped.
	data.TechnicalDescription = template.HTML(data.Product.TechnicalDescription)
	data.RMAProcedure = template.HTML(data.Product.RMAProcedure)

	// Render template.
	err = tmplAldoProduct.ExecuteTemplate(w, "aldoProduct.tmpl", data)
	HandleError(w, err)
}

// Remove mongodb id from Product.
func aldoProductMongodbIdHandlerDelete(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// Update mongodb store id.
	stmt, err := dbAldo.Prepare(`UPDATE product SET mongodb_id = $1 WHERE code = $2;`)
	HandleError(w, err)
	defer stmt.Close()
	fmt.Println("code:", ps.ByName("code"))
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
