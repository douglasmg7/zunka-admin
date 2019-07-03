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

func aldoProductsHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
	}{session, ""}
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
