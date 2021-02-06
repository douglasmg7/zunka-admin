package main

import (
	// "bytes"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/mercadolibre/golang-sdk/sdk"
)

var mercadoLivreAPPID int64
var mercadoLivreSecretKey string
var mercadoLivreRedirectURL string
var mercadoLivreUserCode string

func loadMercadoLivreEnv() {
	// MERCADO_LIVRE_APP_ID
	mercadoLivreAPPIDString := os.Getenv("MERCADO_LIVRE_APP_ID")
	if mercadoLivreAPPIDString == "" {
		panic("MERCADO_LIVRE_APP_ID env not defined.")
	}
	mercadoLivreAPPID, err = strconv.ParseInt(mercadoLivreAPPIDString, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("parsing MERCADO_LIVRE_APP_ID env: %v", err))
	}

	// MERCADO_LIVRE_SECRET_KEY
	mercadoLivreSecretKey = os.Getenv("MERCADO_LIVRE_SECRET_KEY")
	if mercadoLivreSecretKey == "" {
		panic("MERCADO_LIVRE_SECRET_KEY env not defined.")
	}

	// MERCADO_LIVRE_REDIRECT_URL
	mercadoLivreRedirectURL = os.Getenv("MERCADO_LIVRE_REDIRECT_URL")
	if mercadoLivreRedirectURL == "" {
		panic("MERCADO_LIVRE_REDIRECT_URL env not defined.")
	}

	// MERCADO_LIVRE_USER_CODE
	mercadoLivreUserCode = os.Getenv("MERCADO_LIVRE_USER_CODE")
	if mercadoLivreUserCode == "" {
		panic("MERCADO_LIVRE_USER_CODE env not defined.")
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// AUTHENTICATE USER
///////////////////////////////////////////////////////////////////////////////////////////////////
// Login.
func mercadoLivreAuthLoginHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	url := sdk.GetAuthURL(mercadoLivreAPPID, sdk.AuthURLMLA, mercadoLivreRedirectURL)
	log.Printf("url: %v", url)
	http.Redirect(w, req, url, http.StatusSeeOther)
}

// User code.
func mercadoLivreAuthUserHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	Debug.Printf("url: %v", req.URL)
	// Save user code received.
	mercadoLivreUserCode = req.URL.Query().Get("code")
	Debug.Printf("mercado livre user code: %v", mercadoLivreUserCode)

	w.Write([]byte(fmt.Sprintf("ok\nurl:  %v", req.URL)))
}

// Notification.
func mercadoLivreNotificationHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	Debug.Printf("method: %v", req.Method)
	Debug.Printf("url: %v", req.URL)
	w.Write([]byte(fmt.Sprintf("ok\nurl:  %v", req.URL)))
}

// func mercadoLivreUserCodeUserHandlerPost(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
func mercadoLivreUserCodeHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	Debug.Printf("mercado livre user code: %v", mercadoLivreUserCode)
	w.Write([]byte(mercadoLivreUserCode))
}

// func mercadoLivreUserCodeUserHandlerPost(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
func mercadoLivreUsersMeHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	client, err := sdk.Meli(mercadoLivreAPPID, mercadoLivreUserCode, mercadoLivreSecretKey, mercadoLivreRedirectURL)
	if err != nil {
		HandleError(w, err)
		return
	}

	resp, err := client.Get("/users/me")
	if err != nil {
		HandleError(w, err)
		return
	}

	// if err != nil {
	// log.Printf("Error %s\n", err.Error())
	// }

	userInfo, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		HandleError(w, err)
		return
	}

	fmt.Printf("response:%s\n", userInfo)
	w.Write([]byte(userInfo))

	// req.ParseForm()
	// HandleError(w, err)
}

// Get categories.
func mercadoLivreAuthUserHandler2(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {

	url := sdk.GetAuthURL(mercadoLivreAPPID, sdk.AuthURLMLA, "https://www.zunka.com.br/ns/ml/auth")
	log.Printf("url: %v", url)

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
