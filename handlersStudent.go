package main

import (
	"github.com/douglasmg7/bluetang"
	"github.com/julienschmidt/httprouter"

	// _ "github.com/mattn/go-sqlite3"
	"database/sql"
	"log"
	"net/http"
	"time"
)

// Student by email.
func studentByIdHandler(w http.ResponseWriter, req *http.Request, p httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
		Name        string
		Email       string
		Mobile      string
	}{
		Session: session,
	}
	// Get the student.
	err := dbZunka.QueryRow("select name, email, mobile from student where id = ?", p.ByName("id")).Scan(&data.Name, &data.Email, &data.Mobile)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}
	err = tmplStudent.ExecuteTemplate(w, "student.tpl", data)
	HandleError(w, err)
}

// List all stundents.
func allStudentHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	type student struct {
		Id   string
		Name string
	}
	data := struct {
		Session     *SessionData
		HeadMessage string
		Students    []student
	}{session, "", []student{}}
	// names := make([]string, 0)
	// Get all students.
	rows, err := dbZunka.Query("select id, name from student")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var student student
		err := rows.Scan(&student.Id, &student.Name)
		if err != nil {
			log.Fatal(err)
		}
		data.Students = append(data.Students, student)
		// log.Println(id, name)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	err = tmplAllStudent.ExecuteTemplate(w, "allStudent.tpl", data)
	HandleError(w, err)
}

// New student page.
func newStudentHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
		Name        valueMsg
		Email       valueMsg
		Mobile      valueMsg
	}{
		Session: session,
	}
	err = tmplNewStudent.ExecuteTemplate(w, "newStudent.tpl", data)
	HandleError(w, err)
}

// Save new student.
func newStudentHandlerPost(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
		Name        valueMsg
		Email       valueMsg
		Mobile      valueMsg
	}{
		Session: session,
	}
	data.Name.Value, data.Name.Msg = bluetang.Name(req.FormValue("name"))
	data.Email.Value, data.Email.Msg = bluetang.Email(req.FormValue("email"))
	data.Mobile.Value, data.Mobile.Msg = bluetang.Mobile(req.FormValue("mobile"))
	// return page with field erros
	if data.Name.Msg != "" || data.Email.Msg != "" || data.Mobile.Msg != "" {
		err := tmplNewStudent.ExecuteTemplate(w, "newStudent.tpl", data)
		HandleError(w, err)
		// save student
	} else {
		// verify if student name alredy exist
		rows, err := dbZunka.Query("select email from student where email = ?", data.Email.Value)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			data.Email.Msg = "Email já cadastrado"
		}
		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}
		// student alredy registered
		if data.Email.Msg != "" {
			err := tmplNewStudent.ExecuteTemplate(w, "newStudent.tpl", data)
			HandleError(w, err)
			// insert student into db
		} else {
			stmt, err := dbZunka.Prepare(`INSERT INTO student(name, email, mobile, createdAt) VALUES(?, ?, ?, ?)`)
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()
			_, err = stmt.Exec(data.Name.Value, data.Email.Value, data.Mobile.Value, time.Now())
			if err != nil {
				log.Fatal(err)
			}
			http.Redirect(w, req, "/ns/", http.StatusSeeOther)
		}
	}
}
