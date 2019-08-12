package main

import (
	"net/http"
	"path"
	"strings"

	"github.com/douglasmg7/aldoutil"
	"github.com/julienschmidt/httprouter"
)

// Products.
func aldoProductsHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
		Products    []aldoutil.Product
	}{session, "", []aldoutil.Product{}}

	err = dbAldo.Select(&data.Products, "SELECT * FROM product LIMIT 10")
	HandleError(w, err)

	err = tmplAldoProducts.ExecuteTemplate(w, "aldoProducts.tmpl", data)
	HandleError(w, err)
}

// Product.
func aldoProductHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
		Product     aldoutil.Product
	}{session, "", aldoutil.Product{}}

	err = dbAldo.Get(&data.Product, "SELECT * FROM product WHERE code=?", ps.ByName("code"))
	HandleError(w, err)

	// resp, err := http.Post("http://localhost:3080/product-config/product", "asdf", &bug)
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// log.Println(resp.Body)

	err = tmplAldoProduct.ExecuteTemplate(w, "aldoProduct.tmpl", data)
	HandleError(w, err)
}

// Product.
func aldoProductHandlerPost(w http.ResponseWriter, req *http.Request, ps httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
		Product     aldoutil.Product
	}{session, "", aldoutil.Product{}}

	err = dbAldo.Get(&data.Product, "SELECT * FROM product WHERE code=?", ps.ByName("code"))
	HandleError(w, err)

	err = tmplAldoProduct.ExecuteTemplate(w, "aldoProduct.tmpl", data)
	HandleError(w, err)
}

// All categories.
func aldoCategAllHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
		Categories  []string
	}{session, "", []string{}}
	data.Categories = aldoutil.ReadCategoryList(path.Join(listPath, "categAll.list"))
	err = tmplAldoCategoryAll.ExecuteTemplate(w, "aldoCategoryAll.tmpl", data)
	HandleError(w, err)
}

// Selected categories.
func aldoCategSelHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
		Categories  string
	}{session, "", ""}
	data.Categories = strings.Join(aldoutil.ReadCategoryList(path.Join(listPath, "categSel.list")), "\n")
	err = tmplAldoCategorySel.ExecuteTemplate(w, "aldoCategorySel.tmpl", data)
	HandleError(w, err)
}

// Save selected categories.
func aldoCategSelHandlerPost(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	err := aldoutil.WriteCategoryListFromString(req.FormValue("categories"), path.Join(listPath, "categSel.list"))
	HandleError(w, err)

	// categories := strings.Split(strings.ReplaceAll(req.FormValue("categories"), " ", ""), "\n")
	// fmt.Println("Categories:", categories)

	http.Redirect(w, req, "/aldo/category/sel", http.StatusSeeOther)
	return
}

// Categories in use.
func aldoCategUseHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
		Categories  []string
	}{session, "", []string{}}
	data.Categories = aldoutil.ReadCategoryList(path.Join(listPath, "categUse.list"))
	err = tmplAldoCategoryUse.ExecuteTemplate(w, "aldoCategoryUse.tmpl", data)
	HandleError(w, err)
}
