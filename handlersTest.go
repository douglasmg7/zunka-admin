package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

// Test page.
func testPageHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session *SessionData
	}{session}

	err = tmplTest.ExecuteTemplate(w, "test.gohtml", data)
	HandleError(w, err)
}

// Test send mail.
func testSendMailPost(w http.ResponseWriter, req *http.Request, ps httprouter.Params, _ *SessionData) {
	msg := time.Now().String()
	err := sendMail([]string{"douglasmg7@gmail.com"}, "Teste (zunkasrv).", msg)
	if err == nil {
		w.WriteHeader(200)
		return
	}
	HandleError(w, err)
}
