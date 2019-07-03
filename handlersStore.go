package main

import (
	// "github.com/douglasmg7/bluetang"
	"github.com/julienschmidt/httprouter"
	// _ "github.com/mattn/go-sqlite3"
	// "database/sql"
	// "log"
	"net/http"
	// "time"
)

func storeProductsHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
	}{session, ""}
	// fmt.Println("session: ", data.Session)
	err = tmplStoreProducts.ExecuteTemplate(w, "storeProducts.tpl", data)
	HandleError(w, err)
}
