package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	setupTest()
	code := m.Run()
	shutdownTest()

	os.Exit(code)
}

func setupTest() {
	initZunkaDB()
	initAldoDB()
	initAllnationsDB()
}

func shutdownTest() {
	closeZunkaDB()
	closeAldoDB()
	closeAllnationsDB()
}

func TestGetSimiliarProducts(t *testing.T) {
	// title := "Computador All In One Dell Inspiron 5490-M30S2"
	title := "GABINETE COOLER MASTER MASTERBOX LITE 3.1 TG LATERAL EM VIDRO TEMPERADO ATX/E-ATX/MINI-ITX/MICRO-AT"

	chanProducts := make(chan []ZunkaProductRx)
	go getProductsSimilarTitles(chanProducts, title)
	products := <-chanProducts

	for _, product := range products {
		// Debug.Print(product)
		if len(product.StoreProductTitle) == 0 {
			t.Errorf("Empty title")
		}
	}
	// if result != want {
	// t.Errorf("result = %q, want %q", result, want)
	// }
}

func TestGetProductsSameEAN(t *testing.T) {
	// title := "Computador All In One Dell Inspiron 5490-M30S2"
	ean := "7899864928406"

	chanProducts := make(chan []ZunkaProductRx)
	go getProductsSameEAN(chanProducts, ean)
	products := <-chanProducts

	for _, product := range products {
		// Debug.Print(product)
		if len(product.StoreProductTitle) == 0 {
			t.Errorf("Empty title")
		}
	}
	// if result != want {
	// t.Errorf("result = %q, want %q", result, want)
	// }
}
