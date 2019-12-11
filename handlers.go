package main

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type valueMsg struct {
	Value string
	Msg   string
}

// Template message data.
type messageTplData struct {
	Session    *SessionData
	TitleMsg   string
	WarnMsg    string
	SuccessMsg string
}

// Handler error.
func HandleError(w http.ResponseWriter, err error) {
	if err != nil {
		// http.Error(w, "Some thing wrong", 404)
		if dev {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "Alguma coisa deu errado", http.StatusInternalServerError)
		}
		log.Println(err.Error())
		return
	}
}

// Favicon handler.
func faviconHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	http.ServeFile(w, req, "./static/img/favicon.ico")
}

// Index handler.
func indexHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
	}{session, "Aviso de regatta na Lagoa dos Ingleses, dia 18/03/2019"}
	// fmt.Println("session: ", data.Session)
	err = tmplIndex.ExecuteTemplate(w, "index.tpl", data)
	HandleError(w, err)
}

// Clean sessions cache, needed when some db update occurs.
func cleanSessionsHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	sessions.CleanSessions()
	http.Redirect(w, req, "/ns/", http.StatusSeeOther)
}

// Changelog page.
func changelogHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	// Read change log markdown file.
	md, err := ioutil.ReadFile("changelog.md")
	HandleError(w, err)
	// md := []byte("## markdown document")
	output := markdown.ToHTML(md, nil, nil)

	data := struct {
		Session     *SessionData
		HeadMessage string
		Version     string
		Changelog   template.HTML
	}{session, "", version, template.HTML(output)}

	err = tmplChangelog.ExecuteTemplate(w, "changelog.gohtml", data)
	HandleError(w, err)
}

/**************************************************************************************************
* To organizer
**************************************************************************************************/

func userHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "USER, %s!\n", ps.ByName("name"))
}

func userAddHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	err := tmplUserAdd.ExecuteTemplate(w, "user_add.tpl", nil)
	HandleError(w, err)
}
