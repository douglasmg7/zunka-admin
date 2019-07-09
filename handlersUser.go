package main

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/douglasmg7/bluetang"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"

	// _ "github.com/mattn/go-sqlite3"

	"net/http"
)

// Change name template data.
type changeNameTplData struct {
	Session     *SessionData
	HeadMessage string
	NewName     valueMsg
	Password    valueMsg
}

// Change email template data.
type changeEmailTplData struct {
	Session     *SessionData
	HeadMessage string
	NewEmail    valueMsg
	Password    valueMsg
}

// Change mobile number template data.
type changeMobileTplData struct {
	Session     *SessionData
	HeadMessage string
	NewMobile   valueMsg
	Password    valueMsg
}

/**************************************************************************************************
* Account
**************************************************************************************************/

// Account page.
func userAccountHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	data := struct {
		Session     *SessionData
		HeadMessage string
		Name        string
		Email       string
		Mobile      string
		RG          string
		CPF         string
	}{Session: session}

	// Get user data.
	err := dbApp.QueryRow("select name, email, mobile, rg, cpf from user where id = ?", session.UserID).Scan(&data.Name, &data.Email, &data.Mobile, &data.RG, &data.CPF)
	if err != nil {
		log.Fatal(err)
	}
	// Render page.
	err = tmplUserAccount.ExecuteTemplate(w, "userAccount.tpl", data)
	HandleError(w, err)
}

// Change name page.
func userChangeNameHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	var data changeNameTplData
	data.Session = session
	// Get current name.
	var name string
	err = dbApp.QueryRow("SELECT name FROM user WHERE id = ?", session.UserID).Scan(&name)
	data.NewName.Value = name
	// Render page.
	err = tmplUserChangeName.ExecuteTemplate(w, "userChangeName.tpl", data)
	HandleError(w, err)
}

// Change name post.
func userChangeNameHandlerPost(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	var data changeNameTplData
	data.Session = session
	// Check fields.
	data.NewName.Value, data.NewName.Msg = bluetang.Name(req.FormValue("new-name"))
	// Return page with field erros.
	if data.NewName.Msg != "" {
		err := tmplUserChangeName.ExecuteTemplate(w, "userChangeName.tpl", data)
		HandleError(w, err)
		return
	}
	// Wrong password.
	if !session.PasswordIsCorrect(req.FormValue("password")) {
		data.Password.Value = ""
		data.Password.Msg = "Senha incorreta"
		err := tmplUserChangeName.ExecuteTemplate(w, "userChangeName.tpl", data)
		HandleError(w, err)
		return
	}
	// Update name.
	stmt, err := dbApp.Prepare(`UPDATE user SET name = ? WHERE id = ?`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(data.NewName.Value, session.UserID)
	if err != nil {
		log.Fatal(err)
	}
	// Sinalize to session be refreshed from db.
	session.Outdated = true
	// Redirect to account page.
	http.Redirect(w, req, "/user/account", http.StatusSeeOther)
	return
}

// Change email page.
func userChangeEmailHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	var data changeEmailTplData
	data.Session = session
	// Get current email.
	var email string
	err = dbApp.QueryRow("SELECT email FROM user WHERE id = ?", session.UserID).Scan(&email)
	data.NewEmail.Value = email
	// Render page.
	err = tmplUserChangeEmail.ExecuteTemplate(w, "userChangeEmail.tpl", data)
	HandleError(w, err)
}

// Change Email post.
func userChangeEmailHandlerPost(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	var data changeEmailTplData
	data.Session = session
	// Check fields.
	data.NewEmail.Value, data.NewEmail.Msg = bluetang.Email(req.FormValue("new-email"))
	// Return page with field erros.
	if data.NewEmail.Msg != "" {
		err := tmplUserChangeEmail.ExecuteTemplate(w, "userChangeEmail.tpl", data)
		HandleError(w, err)
		return
	}
	// Wrong password.
	if !session.PasswordIsCorrect(req.FormValue("password")) {
		data.Password.Value = ""
		data.Password.Msg = "Senha incorreta"
		err := tmplUserChangeEmail.ExecuteTemplate(w, "userChangeEmail.tpl", data)
		HandleError(w, err)
		return
	}
	// Verify if email alredy in use.
	var userName string
	err := dbApp.QueryRow("select name from user where email = ?", data.NewEmail.Value).Scan(&userName)
	// Email alredy in use.
	if err == nil {
		data.NewEmail.Msg = "Email já cadastrado"
		err := tmplUserChangeEmail.ExecuteTemplate(w, "userChangeEmail.tpl", data)
		HandleError(w, err)
		return
	}
	if (err != nil) && (err != sql.ErrNoRows) {
		log.Fatal(err)
	}
	// Delete email confirmation. Can messed up signup from same e-mail.
	stmt, err := dbApp.Prepare(`DELETE from email_confirmation WHERE email == ?`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(data.NewEmail.Value)
	if err != nil {
		log.Fatal(err)
	}
	// Create uuid.
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}
	// Save email confirmation.
	stmt, err = dbApp.Prepare(`INSERT INTO email_confirmation(uuid, name, email, password, createdAt) VALUES(?, ?, ?, ?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(uuid.String(), "", data.NewEmail.Value, "", time.Now())
	if err != nil {
		log.Fatal(err)
	}
	// Log change email confirmation on dev mode.
	if !production {
		log.Println(`http://localhost:8080/user/change/email-confirmation/` + uuid.String())
	}
	// Render page with next step to complete the email change.
	var dataMsg messageTplData
	dataMsg.Session = session
	dataMsg.TitleMsg = "Pŕoximo passo"
	dataMsg.SuccessMsg = "Dentro de instantes será enviado um e-mail para " + data.NewEmail.Value + " com instruções para completar a alteração do email."
	err = tmplMessage.ExecuteTemplate(w, "message.tpl", dataMsg)
	HandleError(w, err)
}

// Change Email confirmation.
func userChangeEmailConfirmationHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params, session *SessionData) {
	var msgData messageTplData
	msgData.Session = session
	// Find email certify.
	uuid := ps.ByName("uuid")
	var name, newEmail string
	err = dbApp.QueryRow("SELECT name, email FROM email_confirmation WHERE uuid = ?", uuid).Scan(&name, &newEmail)
	// No email confirmation.
	if err != nil {
		msgData.TitleMsg = "Link inválido"
		msgData.WarnMsg = "A alteração do email já foi confirmada anteriormente, ou a tentativa de uma nova alteração de email invalidou este link."
		err := tmplMessage.ExecuteTemplate(w, "message.tpl", msgData)
		HandleError(w, err)
		return
	}
	// Someone trying to signup using the same email.
	if name != "" {
		// Delete email confirmation from signup, so user can try to change to this email again.
		stmt, err := dbApp.Prepare(`DELETE from email_confirmation WHERE uuid == ?`)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(uuid)
		if err != nil {
			log.Fatal(err)
		}
		// Render page message.
		msgData.TitleMsg = "Link inválido"
		msgData.WarnMsg = "Já existe uma tentative de criação de conta utilizando este email."
		err = tmplMessage.ExecuteTemplate(w, "message.tpl", msgData)
		HandleError(w, err)
		return
	}
	// Update user email.
	stmt, err := dbApp.Prepare(`UPDATE user SET email = ? WHERE id = ?`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(newEmail, session.UserID)
	if err != nil {
		log.Fatal(err)
	}
	// Delete email confirmation.
	stmt, err = dbApp.Prepare(`DELETE from email_confirmation WHERE uuid == ?`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(uuid)
	if err != nil {
		log.Fatal(err)
	}
	// Render success message.
	msgData.TitleMsg = "Alteração de email"
	msgData.SuccessMsg = "Email alterado para " + newEmail + "."
	err = tmplMessage.ExecuteTemplate(w, "message.tpl", msgData)
	HandleError(w, err)
	return
}

// Change mobile number page.
func userChangeMobileHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	var data changeMobileTplData
	data.Session = session
	// Get current mobile number.
	var mobile string
	err = dbApp.QueryRow("SELECT mobile FROM user WHERE id = ?", session.UserID).Scan(&mobile)
	data.NewMobile.Value = mobile
	// Render page.
	err = tmplUserChangeMobile.ExecuteTemplate(w, "userChangeMobile.tpl", data)
	HandleError(w, err)
}

// Change mobile number post.
func userChangeMobileHandlerPost(w http.ResponseWriter, req *http.Request, _ httprouter.Params, session *SessionData) {
	var data changeMobileTplData
	data.Session = session
	// Remove spaces.
	data.NewMobile.Value = strings.TrimSpace(req.FormValue("new-mobile"))
	// If some number validate.
	if data.NewMobile.Value != "" {
		// Check fields.
		data.NewMobile.Value, data.NewMobile.Msg = bluetang.Mobile(data.NewMobile.Value)
		// Return page with field erros.
		if data.NewMobile.Msg != "" {
			err := tmplUserChangeMobile.ExecuteTemplate(w, "userChangeMobile.tpl", data)
			HandleError(w, err)
			return
		}
	}
	// Wrong password.
	if !session.PasswordIsCorrect(req.FormValue("password")) {
		data.Password.Value = ""
		data.Password.Msg = "Senha incorreta"
		err := tmplUserChangeMobile.ExecuteTemplate(w, "userChangeMobile.tpl", data)
		HandleError(w, err)
		return
	}
	// Update mobile number.
	stmt, err := dbApp.Prepare(`UPDATE user SET mobile = ? WHERE id = ?`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(data.NewMobile.Value, session.UserID)
	if err != nil {
		log.Fatal(err)
	}
	// Sinalize to session be refreshed from db.
	session.Outdated = true
	// Redirect to account page.
	http.Redirect(w, req, "/user/account", http.StatusSeeOther)
	return
}
