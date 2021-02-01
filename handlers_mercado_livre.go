package main

import (
	// "bytes"

	"net/http"

	"github.com/julienschmidt/httprouter"
)

///////////////////////////////////////////////////////////////////////////////////////////////////
// AUTHENTICATE USER
///////////////////////////////////////////////////////////////////////////////////////////////////
// Get categories.
func mercadoLivreAuthUserHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session          *SessionData
		HeadMessage      string
		User             valueMsg
		Password         valueMsg
		WarnMsgHead      string
		SuccessMsgHead   string
		WarnMsgFooter    string
		SuccessMsgFooter string
	}{session, "", valueMsg{"", ""}, valueMsg{"", ""}, "", "", "", ""}

	// Render page.
	err = tmplMercadoLivreAuthUser.ExecuteTemplate(w, "mercadoLivreAuthUser.gohtml", data)
	HandleError(w, err)
}

// Save categories.
func mercadoLivreAuthUserHandlerPost(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {

	req.ParseForm()
	allnationsSelectedMakers.Makers = []string{}
	for key := range req.PostForm {
		allnationsSelectedMakers.Makers = append(allnationsSelectedMakers.Makers, key)
	}
	allnationsSelectedMakers.UpdateSqlMakers()
	err := allnationsSelectedMakers.Save()
	HandleError(w, err)
	http.Redirect(w, req, "/ns/allnations/products", http.StatusSeeOther)
}
