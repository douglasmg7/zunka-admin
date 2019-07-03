package main

import (
	"net/http"

	"github.com/douglasmg7/aldoutil"
	"github.com/julienschmidt/httprouter"
)

func aldoProductsHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
		Products    []aldoutil.Product
	}{session, "", nil}
	data.Products, err = aldoutil.FindAllProducts(db)
	HandleError(w, err)
	// fmt.Println("session: ", data.Session)
	err = tmplAldoProducts.ExecuteTemplate(w, "aldoProducts.tpl", data)
	HandleError(w, err)
}

func aldoConfigHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
	}{session, ""}
	// fmt.Println("session: ", data.Session)
	err = tmplAldoConfig.ExecuteTemplate(w, "aldoConfig.tpl", data)
	HandleError(w, err)
}
