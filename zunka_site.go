package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Store product to create a new product on store.
type ZunkaSiteProductTx struct {
	// MongodbId                 string    `json:"_id"` // Identify product into store site using mongo _id.
	DealerName                        string    `json:"dealerName"`
	DealerProductId                   string    `json:"dealerProductId"`
	DealerProductTitle                string    `json:"dealerProductTitle"`
	DealerProductDesc                 string    `json:"dealerProductDesc"`
	DealerProductMaker                string    `json:"dealerProductMaker"`
	DealerProductCategory             string    `json:"dealerProductCategory"`
	DealerProductWarrantyDays         int       `json:"dealerProductWarrantyDays"`
	DealerProductDeep                 int       `json:"dealerProductDeep"`   // Deep (comprimento) in cm.
	DealerProductHeight               int       `json:"dealerProductHeight"` // Height in cm.
	DealerProductWidth                int       `json:"dealerProductWidth"`  // Width in cm.
	DealerProductWeight               int       `json:"dealerProductWeight"` // Weight in grams.
	DealerProductActive               bool      `json:"dealerProductActive"`
	DealerProductFinalPriceSuggestion int       `json:"dealerProductFinalPriceSuggestion"`
	DealerProductPrice                int       `json:"dealerProductPrice"`
	DealerProductLastUpdate           time.Time `json:"dealerProductLastUpdate"`
	DealerProductImagesLink           string    `json:"dealerProductImagesLink"` // Images link separated by "__,__".
	DealerProductLocation             string    `json:"dealerProductLocation"`
	StoreProductQtd                   int       `json:"storeProductQtd"`
	Ean                               string    `json:"ean"`
	ProductIdTemplate                 string    `json:"productIdTemplate"` // Identify product to be used as template.
}

// Store product to create a new product on store.
type ZunkaSiteProductRx struct {
	MongodbId         string `json:"_id"` // Identify product into store site using mongo _id.
	StoreProductId    string `json:"storeProductId"`
	StoreProductTitle string `json:"storeProductTitle"`
}

// Get products similar titles.
func getProductsSimilarTitles(c chan []ZunkaSiteProductRx, title string) {
	products := []ZunkaSiteProductRx{}
	// Request product add.
	client := &http.Client{}
	// title = "GABINETE COOLER MASTER MASTERBOX LITE 3.1 TG LATERAL EM VIDRO TEMPERADO ATX/E-ATX/MINI-ITX/MICRO-AT"
	req, err := http.NewRequest("GET", zunkaSiteHost()+"/setup/products-similar-title", nil)
	if err != nil {
		Error.Print(err)
		c <- products
		return
	}
	// Query params
	q := req.URL.Query()
	q.Add("title", title)
	req.URL.RawQuery = q.Encode()
	// Head.
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(zunkaSiteUser(), zunkaSitePass())
	res, err := client.Do(req)
	if err != nil {
		Error.Print(err)
		c <- products
		return
	}
	// res, err := http.Post("http://localhost:3080/setup/product/add", "application/json", bytes.NewBuffer(reqBody))
	defer res.Body.Close()

	// Result.
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		Error.Print(err)
		c <- products
		return
	}
	// No 200 status.
	if res.StatusCode != 200 {
		Error.Print(errors.New(fmt.Sprintf("Getting products with similiar title from zunkasite.\nstatus: %v\nbody: %v", res.StatusCode, string(resBody))))
		c <- products
		return
	}
	err = json.Unmarshal(resBody, &products)
	if err != nil {
		Error.Print(err)
	}
	// Debug.Printf("Product[0]: %v", products[0])
	c <- products
	return
}

// Get products same EAN.
func getProductsSameEAN(c chan []ZunkaSiteProductRx, ean string) {
	products := []ZunkaSiteProductRx{}

	// Request product add.
	client := &http.Client{}
	// title = "GABINETE COOLER MASTER MASTERBOX LITE 3.1 TG LATERAL EM VIDRO TEMPERADO ATX/E-ATX/MINI-ITX/MICRO-AT"
	req, err := http.NewRequest("GET", zunkaSiteHost()+"/setup/products-same-ean", nil)
	if err != nil {
		Error.Print(err)
		c <- products
		return
	}
	// Query params
	q := req.URL.Query()
	q.Add("ean", ean)
	req.URL.RawQuery = q.Encode()
	// Head.
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(zunkaSiteUser(), zunkaSitePass())
	res, err := client.Do(req)
	if err != nil {
		Error.Print(err)
		c <- products
		return
	}
	// res, err := http.Post("http://localhost:3080/setup/product/add", "application/json", bytes.NewBuffer(reqBody))
	defer res.Body.Close()

	// Result.
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		Error.Print(err)
		c <- products
		return
	}
	// No 200 status.
	if res.StatusCode != 200 {
		Error.Print(errors.New(fmt.Sprintf("Getting products same Ean from zunkasite.\nstatus: %v\nbody: %v", res.StatusCode, string(resBody))))
		c <- products
		return
	}
	err = json.Unmarshal(resBody, &products)
	if err != nil {
		Error.Print(err)
	}
	// Debug.Printf("Product[0]: %v", products[0])
	c <- products
	return
}
