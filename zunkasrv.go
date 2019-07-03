package main

import (
	"database/sql"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

/************************************************************************************************
* Templates
************************************************************************************************/
// Geral.
var tmplMaster, tmplIndex, tmplDeniedAccess *template.Template

// Misc.
var tmplMessage *template.Template

// Store.
var tmplStoreProducts *template.Template

// Aldo.
var tmplAldoProducts, tmplAldoConfig *template.Template

// Allnations.
var tmplAllnationsProducts, tmplAllnationsConfig *template.Template

// Auth.
var tmplAuthSignup, tmplAuthSignin, tmplPasswordRecovery, tmplPasswordReset *template.Template

// Student.
var tmplStudent, tmplAllStudent, tmplNewStudent *template.Template

// Db.
var db *sql.DB
var err error

// User.
var tmplUserAdd *template.Template
var tmplUserAccount *template.Template
var tmplUserChangeName *template.Template
var tmplUserChangeEmail *template.Template
var tmplUserChangeMobile *template.Template
var tmplUserChangePassword *template.Template
var tmplUserChangeRG *template.Template
var tmplUserChangeCPF *template.Template
var tmplUserDeleteAccount *template.Template

// Development mode.
var devMode bool

const port = "8080"

// Sessions from each user.
var sessions = Sessions{
	mapUserID:      map[string]int{},
	mapSessionData: map[int]*SessionData{},
}

func init() {
	// Log file.
	logFile, err := os.OpenFile("./log/main.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	// Log cnfiguration.
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	// log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	// log.SetFlags(log.LstdFlags | log.Ldate | log.Lshortfile)
	// log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	// production or development mode
	setMode()

	/************************************************************************************************
	* Load templates
	************************************************************************************************/
	// Geral.
	tmplMaster = template.Must(template.ParseGlob("templates/master/*"))
	tmplIndex = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/index.tpl"))
	tmplDeniedAccess = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/misc/deniedAccess.tpl"))
	// Misc.
	tmplMessage = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/misc/message.tpl"))
	// Store.
	tmplStoreProducts = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/store/storeProducts.tpl"))
	// Aldo.
	tmplAldoProducts = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/aldo/aldoProducts.tpl"))
	tmplAldoConfig = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/aldo/aldoConfig.tpl"))
	// Allnations.
	tmplAllnationsProducts = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/allnations/allnationsProducts.tpl"))
	tmplAllnationsConfig = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/allnations/allnationsConfig.tpl"))
	// Auth.
	tmplAuthSignup = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/auth/signup.tpl"))
	tmplAuthSignin = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/auth/signin.tpl"))
	tmplPasswordRecovery = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/auth/passwordRecovery.tpl"))
	tmplPasswordReset = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/auth/passwordReset.tpl"))
	// Student.
	tmplStudent = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/student.tpl"))
	tmplAllStudent = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/allStudent.tpl"))
	tmplNewStudent = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/newStudent.tpl"))
	// User.
	tmplUserAdd = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/userAdd.tpl"))
	tmplUserAccount = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/user/userAccount.tpl"))
	tmplUserChangeName = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/user/userChangeName.tpl"))
	tmplUserChangeEmail = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/user/userChangeEmail.tpl"))
	tmplUserChangeMobile = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/user/userChangeMobile.tpl"))

	// debug templates
	// for _, tplItem := range tmplAll["user_add"].Templates() {
	// 	log.Println(tplItem.Name())
	// }
}

func main() {
	// Start data base.
	db, err = sql.Open("sqlite3", "./db/zunka.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Init router.
	router := httprouter.New()
	router.GET("/favicon.ico", faviconHandler)
	router.GET("/", getSession(indexHandler))

	// Clean the session cache.
	router.GET("/clean-sessions", checkPermission(cleanSessionsHandler, "admin"))

	// Store.
	router.GET("/store/products", getSession(storeProductsHandler))

	// Aldo.
	router.GET("/aldo/products", getSession(aldoProductsHandler))
	router.GET("/aldo/config", getSession(aldoConfigHandler))

	// Allnations.
	router.GET("/allnations/products", getSession(allnationsProductsHandler))
	router.GET("/allnations/config", getSession(allnationsConfigHandler))

	// Auth - signup.
	router.GET("/auth/signup", confirmNoLogged(authSignupHandler))
	router.POST("/auth/signup", confirmNoLogged(authSignupHandlerPost))
	router.GET("/auth/signup/confirmation/:uuid", confirmNoLogged(authSignupConfirmationHandler))

	// Auth - signin/signout.
	router.GET("/auth/signin", confirmNoLogged(authSigninHandler))
	router.POST("/auth/signin", confirmNoLogged(authSigninHandlerPost))
	router.GET("/auth/signout", authSignoutHandler)

	// Auth - password.
	router.GET("/auth/password/recovery", confirmNoLogged(passwordRecoveryHandler))
	router.POST("/auth/password/recovery", confirmNoLogged(passwordRecoveryHandlerPost))
	router.GET("/auth/password/reset", confirmNoLogged(passwordResetHandler))

	// User.
	router.GET("/user/account", checkPermission(userAccountHandler, ""))
	router.GET("/user/change/name", checkPermission(userChangeNameHandler, ""))
	router.POST("/user/change/name", checkPermission(userChangeNameHandlerPost, ""))
	router.GET("/user/change/email", checkPermission(userChangeEmailHandler, ""))
	router.POST("/user/change/email", checkPermission(userChangeEmailHandlerPost, ""))
	router.GET("/user/change/email-confirmation/:uuid", checkPermission(userChangeEmailConfirmationHandler, ""))
	router.GET("/user/change/mobile", checkPermission(userChangeMobileHandler, ""))
	router.POST("/user/change/mobile", checkPermission(userChangeMobileHandlerPost, ""))

	// Entrance.
	router.GET("/user_add", userAddHandler)

	// Student.
	router.GET("/student/all", checkPermission(allStudentHandler, "editStudent"))
	router.GET("/student/new", checkPermission(newStudentHandler, "editStudent"))
	router.POST("/student/new", checkPermission(newStudentHandlerPost, "editStudent"))
	router.GET("/student/id/:id", checkPermission(studentByIdHandler, "editStudent"))

	// // Example.
	// router.GET("/user/:name", userHandler)

	// start server
	router.ServeFiles("/static/*filepath", http.Dir("./static/"))
	log.Println("listen port", port)
	// Why log.Fall work here?
	// log.Fatal(http.ListenAndServe(":"+port, router))
	log.Fatal(http.ListenAndServe(":"+port, newLogger(router)))
}

/**************************************************************************************************
* Middleware
**************************************************************************************************/

// Handle with session.
type handleS func(w http.ResponseWriter, req *http.Request, p httprouter.Params, session *SessionData)

// Get session middleware.
func getSession(h handleS) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		// Get session.
		session, err := sessions.GetSession(req)
		if err != nil {
			log.Fatal(err)
		}
		h(w, req, p, session)
	}
}

// Check permission middleware.
func checkPermission(h handleS, permission string) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		// Get session.
		session, err := sessions.GetSession(req)
		if err != nil {
			log.Fatal(err)
		}
		// Not signed.
		if session == nil {
			http.Redirect(w, req, "/auth/signin", http.StatusSeeOther)
			return
		}
		// Have the permission.
		if permission == "" || session.CheckPermission(permission) {
			h(w, req, p, session)
			return
		}
		// No Permission.
		// fmt.Fprintln(w, "Not allowed")
		data := struct {
			Session     *SessionData
			HeadMessage string
		}{Session: session}
		err = tmplDeniedAccess.ExecuteTemplate(w, "deniedAccess.tpl", data)
		HandleError(w, err)
	}
}

// Check if not logged.
func confirmNoLogged(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		// Get session.
		session, err := sessions.GetSession(req)
		if err != nil {
			log.Fatal(err)
		}
		// Not signed.
		if session == nil {
			h(w, req, p)
			return
		}
		// fmt.Fprintln(w, "Not allowed")
		data := struct{ Session *SessionData }{session}
		err = tmplDeniedAccess.ExecuteTemplate(w, "deniedAccess.tpl", data)
		HandleError(w, err)

	}
}

/**************************************************************************************************
* Logger middleware
**************************************************************************************************/

// Logger struct.
type logger struct {
	handler http.Handler
}

// Handle interface.
func (l *logger) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, req)
	log.Printf("%s %s %v", req.Method, req.URL.Path, time.Since(start))
}

// New logger.
func newLogger(h http.Handler) *logger {
	return &logger{handler: h}
}

/**************************************************************************************************
* Run mode.
**************************************************************************************************/

// Define production or development mode.
func setMode() {
	for _, arg := range os.Args[1:] {
		if arg == "dev" {
			devMode = true
			log.Println("development mode")
			return
		}
	}
	log.Println("production mode")
}
