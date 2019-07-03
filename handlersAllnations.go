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

func allnationsProductsHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
	}{session, ""}
	// fmt.Println("session: ", data.Session)
	err = tmplAllnationsProducts.ExecuteTemplate(w, "allnationsProducts.tpl", data)
	HandleError(w, err)
}

func allnationsConfigHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
	}{session, ""}
	// fmt.Println("session: ", data.Session)
	err = tmplAllnationsConfig.ExecuteTemplate(w, "allnationsConfig.tpl", data)
	HandleError(w, err)
}
