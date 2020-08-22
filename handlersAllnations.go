package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	// "os/exec"
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

///////////////////////////////////////////////////////////////////////////////////////////////////
// CATEGORY
///////////////////////////////////////////////////////////////////////////////////////////////////
// Get all categories.
func allnationsCategoriesHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
		Categories  []AllnationsCategory
	}{session, "", []AllnationsCategory{}}

	// Get categories from db.
	err = dbAllnations.Select(&data.Categories, "SELECT * FROM category order by name")
	HandleError(w, err)

	log.Printf("Categories: %v", data.Categories)
	// Render page.
	err = tmplAllnationsCategories.ExecuteTemplate(w, "allnationsCategories.tmpl", data)
	HandleError(w, err)
}

// Save categories.
func allnationsCategoriesHandlerPost(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	// Get categories from db.
	categories := []AllnationsCategory{}
	err = dbAllnations.Select(&categories, "SELECT * FROM category order by name")
	HandleError(w, err)
	// Prepare update.
	stmt, err := dbAllnations.Prepare(`UPDATE category SET selected = $1 WHERE name = $2;`)
	HandleError(w, err)
	defer stmt.Close()

	// For each category on db.
	for _, category := range categories {
		// Update changed categories.
		if (category.Selected && (req.PostFormValue(category.Name) == "")) || (!category.Selected && (req.PostFormValue(category.Name) != "")) {
			// fmt.Println("Updated category:", category.Name)
			_, err = stmt.Exec(!category.Selected, category.Name)
			HandleError(w, err)
		}
	}
	// Render categories page.
	http.Redirect(w, req, "/ns/allnations/categories", http.StatusSeeOther)

	// todo - Run script to process xml products.
	// cmd := exec.Command(GS + "/aldowsc/bin/process-xml-products.sh")
	// err = cmd.Start()
	// if err != nil {
	// log.Printf("Error running script to process XML Allnations products. %s", err)
	// }

	return
}
