package main

import (
	// "bytes"

	"bytes"
	"encoding/json"
	"errors"
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

// type mlUserMe struct {
// id     int    `json:"id"`
// nickname    string `json:"nickname"`
// registration_date date    `json:"registration_date"`
// permalink    string    `json:"permalink"`
// MsgError string `xml:"MsgErro"`
// }

// Initialize Mercado Livre handler
func initMercadoLivreHandler() {
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
	setMLUserCode("TG-60228432dbe8c8000639e79c-360790045")
	mercadoLivreUserCode = getMLUserCode()
	// Debug.Printf("mercadoLivreUserCode: %v", mercadoLivreUserCode)
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// AUTHENTICATE USER
///////////////////////////////////////////////////////////////////////////////////////////////////
// Login
// Redirect user to Mercado Livre login page.
func mercadoLivreAuthLoginHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	// url := sdk.GetAuthURL(mercadoLivreAPPID, sdk.AuthURLMLA, mercadoLivreRedirectURL)
	url := sdk.GetAuthURL(mercadoLivreAPPID, sdk.AuthURLMLB, mercadoLivreRedirectURL)
	log.Printf("url: %v", url)
	http.Redirect(w, req, url, http.StatusSeeOther)
}

// User code.
// After user has logged into ML, ML call this handler to pass the user code.
func mercadoLivreAuthUserHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	Debug.Printf("url: %v", req.URL)
	// Save user code received.
	mercadoLivreUserCode = req.URL.Query().Get("code")
	setMLUserCode(mercadoLivreUserCode)
	Debug.Printf("mercado livre user code: %v", mercadoLivreUserCode)

	// w.Write([]byte(fmt.Sprintf("ok\nurl:  %v", req.URL)))
	w.Write([]byte(fmt.Sprintf("%v", req.URL)))
}

// Notification.
func mercadoLivreNotificationHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	Debug.Printf("method: %v", req.Method)
	Debug.Printf("url: %v", req.URL)
	w.Write([]byte(fmt.Sprintf("ok\nurl:  %v", req.URL)))
}

// Show ML user code.
func mercadoLivreUserCodeHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session *SessionData
		Text    string
	}{&SessionData{}, ""}

	data.Text = mercadoLivreUserCode
	err = tmplMercadoLivreText.ExecuteTemplate(w, "mercadoLivreText.gohtml", data)
	HandleError(w, err)
}

// Send ML user code.
func mercadoLivreAPIUserCodeHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if mercadoLivreUserCode != "" {
		w.Write([]byte(mercadoLivreUserCode))
	} else {
		w.Write([]byte("No user code"))
	}
}

// Show user info.
func mercadoLivreUsersMeHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session *SessionData
		Text    string
	}{&SessionData{}, ""}

	// No user code.
	if mercadoLivreUserCode == "" {
		http.Redirect(w, req, "../auth/login", http.StatusSeeOther)
	}

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

	userInfo, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		HandleError(w, err)
		return
	}

	var out bytes.Buffer
	json.Indent(&out, userInfo, "", "\t")
	// out.WriteTo(os.Stdout)

	fmt.Printf("response:%s\n", out.String())
	// w.Write([]byte(out.String()))

	data.Text = out.String()
	err = tmplMercadoLivreText.ExecuteTemplate(w, "mercadoLivreText.gohtml", data)
	HandleError(w, err)
}

// Load user code from zunka server in production.
func mercadoLivreLoadUserCode(w http.ResponseWriter, req *http.Request, ps httprouter.Params, session *SessionData) {
	data := struct {
		Session *SessionData
		Text    string
	}{session, ""}

	// Request user code from production server.
	client := &http.Client{}
	// req, err = http.NewRequest("GET", "http://localhost:8080/ns/ml/api/user-code", nil)
	req, err = http.NewRequest("GET", "https://www.zunka.com.br/ns/ml/api/user-code", nil)
	req.Header.Set("Content-Type", "application/json")
	HandleError(w, err)
	// req.SetBasicAuth(zunkaServerUser(), zunkaServerPass())
	req.SetBasicAuth(zunkaServerUserProduction, zunkaServerPassProduction)
	res, err := client.Do(req)
	HandleError(w, err)

	// res, err := http.Post("http://localhost:3080/setup/product/add", "application/json", bytes.NewBuffer(reqBody))
	defer res.Body.Close()
	HandleError(w, err)

	// Result.
	resBody, err := ioutil.ReadAll(res.Body)
	HandleError(w, err)
	// No 200 status.
	if res.StatusCode != 200 {
		HandleError(w, errors.New(fmt.Sprintf("Error ao solicitar ml user code no servidor zunka.\n\nstatus: %v\n\nbody: %v", res.StatusCode, string(resBody))))
		return
	}
	mercadoLivreUserCode = string(resBody)
	Debug.Printf("Mercado Livre user code loaded: %v", mercadoLivreUserCode)

	data.Text = "Código do usuário carregado com sucesso"
	err = tmplMercadoLivreText.ExecuteTemplate(w, "mercadoLivreText.gohtml", data)
	HandleError(w, err)
}
