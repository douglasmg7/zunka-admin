package main

import (
	"fmt"

	"github.com/douglasmg7/bluetang"
	"github.com/julienschmidt/httprouter"

	// _ "github.com/mattn/go-sqlite3"
	"database/sql"
	"log"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type signinTplData struct {
	Session          *SessionData
	HeadMessage      string
	Email            valueMsg
	Password         valueMsg
	WarnMsgHead      string
	SuccessMsgHead   string
	WarnMsgFooter    string
	SuccessMsgFooter string
}
type signupTplData struct {
	Session         *SessionData
	HeadMessage     string
	Name            valueMsg
	Email           valueMsg
	Password        valueMsg
	PasswordConfirm valueMsg
	WarnMsg         string
	SuccessMsg      string
}
type passwordRecoveryTplData struct {
	Session          *SessionData
	HeadMessage      string
	Email            valueMsg
	WarnMsgFooter    string
	SuccessMsgFooter string
}
type passwordResetTplData struct {
	Session          *SessionData
	HeadMessage      string
	Email            valueMsg
	EmailConfirm     valueMsg
	WarnMsgFooter    string
	SuccessMsgFooter string
}

/**************************************************************************************************
* Signup
**************************************************************************************************/

// Signup page.
func authSignupHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	data := signupTplData{}
	err := tmplAuthSignup.ExecuteTemplate(w, "signup.tpl", data)
	HandleError(w, err)
}

// Signup post.
func authSignupHandlerPost(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var dataMsg messageTplData
	var data signupTplData
	data.Name.Value, data.Name.Msg = bluetang.Name(req.FormValue("name"))
	data.Email.Value, data.Email.Msg = bluetang.Email(req.FormValue("email"))
	data.Password.Value, data.Password.Msg = bluetang.Password(req.FormValue("password"))
	// Check confirm email equality.
	if data.Password.Msg == "" {
		if req.FormValue("password") != req.FormValue("passwordConfirm") {
			data.PasswordConfirm.Msg = "Confirmação da senha e senha devem ser iguais"
		}
	}
	// Return page with field erros.
	if data.Name.Msg != "" || data.Email.Msg != "" || data.Password.Msg != "" || data.PasswordConfirm.Msg != "" {
		err := tmplAuthSignup.ExecuteTemplate(w, "signup.tpl", data)
		HandleError(w, err)
		return
	}
	// Verify if email alredy registered.
	rows, err := dbApp.Query("select email from user where email = ?", data.Email.Value)
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
	// Email alredy registered.
	if data.Email.Msg != "" {
		err := tmplAuthSignup.ExecuteTemplate(w, "signup.tpl", data)
		HandleError(w, err)
		return
	}
	// Lookup for a recent email confirmation.
	var createdAt time.Time
	err = dbApp.QueryRow("SELECT createdAt FROM email_confirmation WHERE email = ?", data.Email.Value).Scan(&createdAt)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}
	// Alredy have email confirmation on db.
	if createdAt.IsZero() == false {
		// Not accept new signup in less than 5 min.
		if time.Since(createdAt).Minutes() < float64(5) {
			dataMsg.TitleMsg = "Solicitação já realizada anteriormente"
			dataMsg.WarnMsg = "A solicitação de cadastramento para o email " + data.Email.Value + " já foi realizada anteriormente, falta a confirmação do cadastro atravéz do link enviado para o respectivo email."
			err = tmplMessage.ExecuteTemplate(w, "message.tpl", dataMsg)
			HandleError(w, err)
			return
		}
		// Delete old email confirmation.
		stmt, err := dbApp.Prepare(`DELETE from email_confirmation WHERE email == ?`)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(data.Email.Value)
		if err != nil {
			log.Fatal(err)
		}
	}
	// Create uuid.
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}
	// Encrypt password.
	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password.Value), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}
	// Save email confirmation.
	stmt, err := dbApp.Prepare(`INSERT INTO email_confirmation(uuid, name, email, password, createdAt) VALUES(?, ?, ?, ?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(uuid.String(), data.Name.Value, data.Email.Value, cryptedPassword, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	// Log email confirmation on dev mode.
	if !production {
		log.Println(`http://localhost:8080/auth/signup/confirmation/` + uuid.String())
	}
	// Render page with next step to complete signup.
	dataMsg.TitleMsg = "Pŕoximo passo"
	dataMsg.SuccessMsg = "Dentro de instantes será enviado um e-mail para " + data.Email.Value + " com instruções para completar o cadastro."
	err = tmplMessage.ExecuteTemplate(w, "message.tpl", dataMsg)
	HandleError(w, err)
}

// Signup confirmation.
func authSignupConfirmationHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// Find email confirmation.
	uuid := ps.ByName("uuid")
	var name, email string
	var password []byte
	var data signinTplData
	err = dbApp.QueryRow("SELECT name, email, password FROM email_confirmation WHERE uuid = ?", uuid).Scan(&name, &email, &password)
	if err == sql.ErrNoRows {
		var msgData messageTplData
		msgData.TitleMsg = "Link inválido"
		msgData.WarnMsg = "O cadastro já foi confirmado anteriormente, ou a tentativa de gerar o cadastro novamente invalidou este link."
		err := tmplMessage.ExecuteTemplate(w, "message.tpl", msgData)
		HandleError(w, err)
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	// Someone trying to change email to this same email.
	if name == "" {
		// Delete email confirmation from change email, so user can try to signup with this email again.
		stmt, err := dbApp.Prepare(`DELETE from email_confirmation WHERE uuid == ?`)
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(uuid)
		if err != nil {
			log.Fatal(err)
		}
		var msgData messageTplData
		msgData.TitleMsg = "Link inválido"
		msgData.WarnMsg = "Já existe uma tentativa de alteração de email para este mesmo email."
		err = tmplMessage.ExecuteTemplate(w, "message.tpl", msgData)
		HandleError(w, err)
		return
	}
	// Create a user from email confirmation.
	stmt, err := dbApp.Prepare(`INSERT INTO user(name, email, password, createdAt, updatedAt) VALUES(?, ?, ?, ?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	now := time.Now()
	_, err = stmt.Exec(name, email, password, now, now)
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
	data.SuccessMsgHead = "Cadastro concluído, você já pode entrar."
	data.WarnMsgHead = ""
	err = tmplAuthSignin.ExecuteTemplate(w, "signin.tpl", data)
	HandleError(w, err)
}

/**************************************************************************************************
* Signin
**************************************************************************************************/

// Signin page.
func authSigninHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	data := signinTplData{}
	err := tmplAuthSignin.ExecuteTemplate(w, "signin.tpl", data)
	HandleError(w, err)
}

// Signin post.
func authSigninHandlerPost(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var data signinTplData
	// Test email format.
	data.Email.Value, data.Email.Msg = bluetang.Email(req.FormValue("email"))
	if data.Email.Msg != "" {
		err := tmplAuthSignin.ExecuteTemplate(w, "signin.tpl", data)
		HandleError(w, err)
		return
	}
	// Get user by email.
	var userID int
	var cryptedPassword []byte
	err = dbApp.QueryRow("SELECT id, password FROM user WHERE email = ?", data.Email.Value).Scan(&userID, &cryptedPassword)
	// no registred user
	if err == sql.ErrNoRows {
		data.Email.Msg = "Email não cadastrado"
		err := tmplAuthSignin.ExecuteTemplate(w, "signin.tpl", data)
		HandleError(w, err)
		return
	}
	// Internal error.
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}
	// Test password format.
	data.Password.Value, data.Password.Msg = bluetang.Password(req.FormValue("password"))
	if data.Password.Msg != "" {
		err := tmplAuthSignin.ExecuteTemplate(w, "signin.tpl", data)
		HandleError(w, err)
		return
	}
	// Test password.
	err = bcrypt.CompareHashAndPassword(cryptedPassword, []byte(data.Password.Value))
	// Incorrect password.
	if err != nil {
		data.Password.Msg = "Senha incorreta"
		err := tmplAuthSignin.ExecuteTemplate(w, "signin.tpl", data)
		HandleError(w, err)
		return
	}
	// Create session.
	err = sessions.CreateSession(w, userID)
	if err != nil {
		log.Fatal(err)
	}
	// Logged, redirect to main page.
	http.Redirect(w, req, "/", http.StatusSeeOther)
	return
}

// Signout.
func authSignoutHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	sessions.RemoveSession(w, req)
}

/**************************************************************************************************
* Reset password
**************************************************************************************************/
// Password recovery page.
func passwordRecoveryHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	data := passwordRecoveryTplData{}
	err := tmplPasswordRecovery.ExecuteTemplate(w, "passwordRecovery.tpl", data)
	HandleError(w, err)
}

// Password recovery post.
func passwordRecoveryHandlerPost(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var data passwordRecoveryTplData
	// Verify email format.
	data.Email.Value, data.Email.Msg = bluetang.Email(req.FormValue("email"))
	if data.Email.Msg != "" {
		err := tmplPasswordRecovery.ExecuteTemplate(w, "passwordRecovery.tpl", data)
		HandleError(w, err)
		return
	}
	// Get user by email.
	var userID int
	err = dbApp.QueryRow("SELECT id FROM user WHERE email = ?", data.Email.Value).Scan(&userID)
	// No user.
	if err == sql.ErrNoRows {
		data.WarnMsgFooter = "Email não cadastrado."
		err := tmplPasswordRecovery.ExecuteTemplate(w, "passwordRecovery.tpl", data)
		HandleError(w, err)
		return
	}
	// Internal error.
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}
	// Create a token to change password.
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}
	// Save token.
	stmt, err := dbApp.Prepare(`INSERT INTO password_reset(uuid, user_email, createdAt) VALUES(?, ?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(uuid.String(), data.Email.Value, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	// Log email confirmation on dev mode.
	if !production {
		log.Println(`http://localhost:8080/auth/password/reset/` + uuid.String())
	}
	// Render page with next step to reset password.
	var dataMsg messageTplData
	dataMsg.TitleMsg = "Pŕoximo passo"
	dataMsg.SuccessMsg = fmt.Sprintf("Foi enviado um e-mail para %s com as instruções para a recuperação da senha.", data.Email.Value)
	err = tmplMessage.ExecuteTemplate(w, "message.tpl", dataMsg)
	HandleError(w, err)
}

// Password reset page.
func passwordResetHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	data := passwordResetTplData{}
	err := tmplPasswordReset.ExecuteTemplate(w, "passwordReset.tpl", data)
	HandleError(w, err)
}

// let emailOptions = {
//   from: '',
//   to: req.body.email,
//   subject: 'Solicitação de criação de conta no site da Zunka.',
//   text: 'Você recebeu este e-mail porquê você (ou alguem) requisitou a criação de uma conta no site da Zunka (https://www.zunka.com.br) usando este e-mail.\n\n' +
//   'Por favor clique no link, ou cole-o no seu navegador de internet para concluir a criação da conta.\n\n' +
//   'https://' + req.app.get('hostname') + '/user/signin/' + token + '\n\n' +
//   'Se não foi você que requisitou esta criação de conta, por favor, ignore este e-mail e nenhuma conta será criada.',
// };
