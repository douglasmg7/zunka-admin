package main

import (
	"time"
)

// Store product to create a new product on store.
type ZunkaProduct struct {
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
}
