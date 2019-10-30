package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
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
	msg := "Email padrão\r\n" +
		"Favor não responder\r\n"
	err := sendMail([]string{"douglasmg7@gmail.com"}, "Criação de conta", msg)
	HandleError(w, err)
	w.WriteHeader(200)
}
