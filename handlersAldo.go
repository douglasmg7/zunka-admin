package main

import (
	"net/http"
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
	data.Categories = aldoutil.ReadCategoryList("../aldowsc/list/categAll.list")
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
	data.Categories = strings.Join(aldoutil.ReadCategoryList("../aldowsc/list/categSel.list"), "\n")
	err = tmplAldoCategorySel.ExecuteTemplate(w, "aldoCategorySel.tmpl", data)
	HandleError(w, err)
}

// Categories in use.
func aldoCategUseHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
		Categories  []string
	}{session, "", []string{}}
	data.Categories = aldoutil.ReadCategoryList("../aldowsc/list/categUse.list")
	err = tmplAldoCategoryUse.ExecuteTemplate(w, "aldoCategoryUse.tmpl", data)
	HandleError(w, err)
}
