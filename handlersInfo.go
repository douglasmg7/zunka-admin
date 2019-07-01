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

// Institutional.
func institutionalHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
	}{session, ""}
	// fmt.Println("session: ", data.Session)
	err = tmplInstitutional.ExecuteTemplate(w, "institutional.tpl", data)
	HandleError(w, err)
}

// Children sailing lessons.
func childrensSailingLessons(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
	}{session, ""}
	// fmt.Println("session: ", data.Session)
	err = tmplChildrenSailingLessons.ExecuteTemplate(w, "childrensSailingLessons.tpl", data)
	HandleError(w, err)
}

// Adults sailing lessons.
func adultsSailingLessons(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
	}{session, ""}
	// fmt.Println("session: ", data.Session)
	// err = tmplAdultsSailingLessons.ExecuteTemplate(w, "adultsSailingLessons.tpl", data)
	err = tmplAdultsSailingLessons.ExecuteTemplate(w, "adultsSailingLessons.tpl", data)
	HandleError(w, err)
}

// Rowing lessons.
func rowingLessons(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
	}{session, ""}
	// fmt.Println("session: ", data.Session)
	err = tmplRowingLessons.ExecuteTemplate(w, "rowingLessons.tpl", data)
	HandleError(w, err)
}

// Sailboat rental.
func sailboatRental(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
	}{session, ""}
	// fmt.Println("session: ", data.Session)
	err = tmplSailboatRental.ExecuteTemplate(w, "sailboatRental.tpl", data)
	HandleError(w, err)
}

// Kayaks and aquatic bikes rental.
func kayaksAndAquaticBikesRental(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
	}{session, ""}
	// fmt.Println("session: ", data.Session)
	err = tmplKayaksAndAquaticBikesRental.ExecuteTemplate(w, "kayaksAndAquaticBikesRental.tpl", data)
	HandleError(w, err)
}

// Sailboat ride.
func sailboatRide(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
	}{session, ""}
	// fmt.Println("session: ", data.Session)
	err = tmplSailboatRide.ExecuteTemplate(w, "sailboatRide.tpl", data)
	HandleError(w, err)
}

// Projects and initiatives.
func projectsAndInitiatives(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
	}{session, ""}
	// fmt.Println("session: ", data.Session)
	err = tmplProjectsAndInitiatives.ExecuteTemplate(w, "projectsAndInitiatives.tpl", data)
	HandleError(w, err)
}

// Contato.
func contact(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
	}{session, ""}
	// fmt.Println("session: ", data.Session)
	err = tmplContact.ExecuteTemplate(w, "contact.tpl", data)
	HandleError(w, err)
}

// Sutents area.
func studentsArea(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
	}{session, ""}
	// fmt.Println("session: ", data.Session)
	err = tmplStudentsArea.ExecuteTemplate(w, "studentsArea.tpl", data)
	HandleError(w, err)
}
