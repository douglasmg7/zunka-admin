package main

import (
	// "bytes"

	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/mercadolibre/golang-sdk/sdk"
)

var mercadoLivreAPPID int64
var mercadoLivreSecretKey string
var mercadoLivreRedirectURL = "https://www.zunka.com.br/ns/ml/products"

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
	log.Printf("mercadoLivreSecretKey: %v", mercadoLivreSecretKey)
	if mercadoLivreSecretKey == "" {
		panic("MERCADO_LIVRE_SECRET_KEY env not defined.")
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// AUTHENTICATE USER
///////////////////////////////////////////////////////////////////////////////////////////////////
// Login.
func mercadoLivreAuthLoginHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	url := sdk.GetAuthURL(mercadoLivreAPPID, sdk.AuthURLMLA, "https://www.zunka.com.br/ns/ml/auth/user")
	log.Printf("url: %v", url)
	http.Redirect(w, req, url, http.StatusSeeOther)
}

// User code.
func mercadoLivreAuthUserHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	userCode := req.URL.Query().Get("code")
	Debug.Printf("url: %v", req.URL)
	Debug.Printf("user code: %v", userCode)
	w.Write([]byte(fmt.Sprintf("ok\nurl:  %v", req.URL)))
}

func mercadoLivreUserCodeUserHandlerPost(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {

	w.Write([]byte("oi"))

	// UserCode := "2345"

	// client, err := sdk.Meli(mercadoLivreAPPID, UserCode, mercadoLivreSecretKey, mercadoLivreRedirectURL)
	// resp, err := client.Get("/users/me")

	// if err != nil {
	// log.Printf("Error %s\n", err.Error())
	// }
	// userInfo, _ := ioutil.ReadAll(resp.Body)
	// fmt.Printf("response:%s\n", userInfo)

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

// Save categories.
func mercadoLivreAuthUserHandlerPost(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {

	w.Write([]byte("oi"))

	// UserCode := "2345"

	// client, err := sdk.Meli(mercadoLivreAPPID, UserCode, mercadoLivreSecretKey, mercadoLivreRedirectURL)
	// resp, err := client.Get("/users/me")

	// if err != nil {
	// log.Printf("Error %s\n", err.Error())
	// }
	// userInfo, _ := ioutil.ReadAll(resp.Body)
	// fmt.Printf("response:%s\n", userInfo)

	// req.ParseForm()
	// HandleError(w, err)
}
