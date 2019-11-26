package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
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

// Aldo.
var tmplAldoProducts, tmplAldoProduct, tmplAldoCategories *template.Template

// Allnations.
var tmplAllnationsProducts, tmplAllnationsConfig *template.Template

// Auth.
var tmplAuthSignup, tmplAuthSignin, tmplPasswordRecovery, tmplPasswordReset *template.Template

// Student.
var tmplStudent, tmplAllStudent, tmplNewStudent *template.Template

// Test.
var tmplTest *template.Template

var dev bool
var port string

// Db.
var dbZunka *sql.DB
var dbAldo *sqlx.DB
var dbZunkaFile, dbAldoFile string

var zunkaPath string
var GS string

var err error

// Sessions from each user.
var sessions = Sessions{
	mapUserID:      map[string]int{},
	mapSessionData: map[int]*SessionData{},
}

func init() {
	// Port.
	port = "8080"

	// Path.
	zunkaPath = os.Getenv("ZUNKAPATH")
	if zunkaPath == "" {
		panic("ZUNKAPATH env not defined.")
	}

	// Go lang source.
	GS = os.Getenv("GS")
	if GS == "" {
		panic("GS env not defined.")
	}

	logPath := path.Join(zunkaPath, "log", "zunka-srv")
	// Create log path.
	os.MkdirAll(logPath, os.ModePerm)

	// Zunka db.
	dbZunkaFile = os.Getenv("ZUNKA_SRV_DB")
	if dbZunkaFile == "" {
		panic("ZUNKASRV_DB not defined.")
	}
	dbZunkaFile = path.Join(zunkaPath, "db", dbZunkaFile)

	// Aldo db.
	dbAldoFile = os.Getenv("ZUNKA_ALDOWSC_DB")
	if dbAldoFile == "" {
		panic("ZUNKA_ALDOWSC_DB not defined.")
	}
	dbAldoFile = path.Join(zunkaPath, "db", dbAldoFile)
	// log.Println("aldoDb:", dbAldoFile)

	// Log file.
	logFile, err := os.OpenFile(path.Join(logPath, "zunkasrv.log"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	// Log configuration.
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	// log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	// log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	// log.SetFlags(log.LstdFlags | log.Ldate | log.Lshortfile)
	// log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// Run mode.
	mode := "production"
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "dev") {
		dev = true
		mode = "development"
	}
	// Log start.
	log.Printf("*** Starting zunka srv in %v mode (version %s) ***\n", mode, version)

	/************************************************************************************************
	* Load templates
	************************************************************************************************/
	// Geral.
	tmplMaster = template.Must(template.ParseGlob("templates/master/*"))
	tmplIndex = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/index.tpl"))
	tmplDeniedAccess = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/misc/deniedAccess.tpl"))
	// Misc.
	tmplMessage = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/misc/message.tpl"))
	tmplTest = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/misc/test.gohtml"))
	// User.
	tmplUserAdd = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/userAdd.tpl"))
	tmplUserAccount = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/user/userAccount.tpl"))
	tmplUserChangeName = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/user/userChangeName.tpl"))
	tmplUserChangeEmail = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/user/userChangeEmail.tpl"))
	tmplUserChangeMobile = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/user/userChangeMobile.tpl"))
	// Aldo.
	tmplAldoProduct = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/aldo/aldoProduct.tmpl"))
	tmplAldoProducts = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/aldo/aldoProducts.tmpl"))
	tmplAldoCategories = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/aldo/aldoCategories.tmpl"))
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

	// debug templates
	// for _, tplItem := range tmplAll["user_add"].Templates() {
	// 	log.Println(tplItem.Name())
	// }
}

func main() {
	// Start data base.
	dbZunka, err = sql.Open("sqlite3", dbZunkaFile)
	if err != nil {
		log.Fatal(fmt.Errorf("Error open db %v: %v", dbZunka, err))
	}
	defer dbZunka.Close()
	err = dbZunka.Ping()
	if err != nil {
		log.Fatal(err)
	}

	dbAldo = sqlx.MustConnect("sqlite3", dbAldoFile)
	defer dbAldo.Close()

	// Init router.
	router := httprouter.New()
	router.GET("/favicon.ico", faviconHandler)
	router.GET("/ns/favicon.ico", faviconHandler)
	router.GET("/", getSession(indexHandler))
	router.GET("/ns/", getSession(indexHandler))

	// Clean the session cache.
	router.GET("/ns/clean-sessions", checkPermission(cleanSessionsHandler, "admin"))

	// Aldo.
	// Products list page.
	router.GET("/ns/aldo/products", checkPermission(aldoProductsHandler, "read"))
	// Product page.
	router.GET("/ns/aldo/product/:code", checkPermission(aldoProductHandler, "read"))
	// Create product on zunka server.
	router.POST("/ns/aldo/product/:code", checkPermission(aldoProductHandlerPost, "write"))
	// Product removed from site, so remove his reference from the site system.
	router.DELETE("/ns/aldo/product/mongodb_id/:code", checkApiAuthorization(aldoProductMongodbIdHandlerDelete))
	// Categories page.
	router.GET("/ns/aldo/categories", checkPermission(aldoCategoriesHandler, "read"))
	// Save categories.
	router.POST("/ns/aldo/categories", checkPermission(aldoCategoriesHandlerPost, "write"))

	// Allnations.
	router.GET("/allnations/products", checkPermission(allnationsProductsHandler, "read"))
	router.GET("/allnations/config", checkPermission(allnationsConfigHandler, "read"))

	// Auth - signup.
	router.GET("/ns/auth/signup", confirmNoLogged(authSignupHandler))
	router.POST("/ns/auth/signup", confirmNoLogged(authSignupHandlerPost))
	router.GET("/ns/auth/signup/confirmation/:uuid", confirmNoLogged(authSignupConfirmationHandler))

	// Auth - signin/signout.
	router.GET("/ns/auth/signin", confirmNoLogged(authSigninHandler))
	router.POST("/ns/auth/signin", confirmNoLogged(authSigninHandlerPost))
	router.GET("/ns/auth/signout", authSignoutHandler)

	// Auth - password.
	router.GET("/ns/auth/password/recovery", confirmNoLogged(passwordRecoveryHandler))
	router.POST("/ns/auth/password/recovery", confirmNoLogged(passwordRecoveryHandlerPost))
	router.GET("/ns/auth/password/reset", confirmNoLogged(passwordResetHandler))

	// User.
	router.GET("/ns/user/account", checkPermission(userAccountHandler, ""))
	router.GET("/ns/user/change/name", checkPermission(userChangeNameHandler, ""))
	router.POST("/ns/user/change/name", checkPermission(userChangeNameHandlerPost, ""))
	router.GET("/ns/user/change/email", checkPermission(userChangeEmailHandler, ""))
	router.POST("/ns/user/change/email", checkPermission(userChangeEmailHandlerPost, ""))
	router.GET("/ns/user/change/email-confirmation/:uuid", checkPermission(userChangeEmailConfirmationHandler, ""))
	router.GET("/ns/user/change/mobile", checkPermission(userChangeMobileHandler, ""))
	router.POST("/ns/user/change/mobile", checkPermission(userChangeMobileHandlerPost, ""))

	// Entrance.
	router.GET("/ns/user_add", userAddHandler)

	// Student.
	router.GET("/ns/student/all", checkPermission(allStudentHandler, "editStudent"))
	router.GET("/ns/student/new", checkPermission(newStudentHandler, "editStudent"))
	router.POST("/ns/student/new", checkPermission(newStudentHandlerPost, "editStudent"))
	router.GET("/ns/student/id/:id", checkPermission(studentByIdHandler, "editStudent"))

	// Test.
	router.GET("/ns/test", checkPermission(testPageHandler, "admin"))
	router.POST("/ns/test/send-email", checkPermission(testSendMailPost, "admin"))

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
			http.Redirect(w, req, "/ns/auth/signin", http.StatusSeeOther)
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

// Api Authorization.
func checkApiAuthorization(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		user, pass, ok := req.BasicAuth()
		if ok && user == zunkaServerUser() && pass == zunkaServerPass() {
			h(w, req, p)
			return
		}
		log.Printf("Unauthorized access, %v %v, user: %v, pass: %v, ok: %v", req.Method, req.URL.Path, user, pass, ok)
		log.Printf("authorization      , %v %v, user: %v, pass: %v", req.Method, req.URL.Path, zunkaServerUser(), zunkaServerPass())
		// Unauthorised.
		w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password for this service"`)
		w.WriteHeader(401)
		w.Write([]byte("Unauthorised.\n"))
		return
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
// todo - why DELETE is logging twice?
func (l *logger) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, req)
	log.Printf("%s %s %v", req.Method, req.URL.Path, time.Since(start))
	// log.Printf("header: %v", req.Header)
}

// New logger.
func newLogger(h http.Handler) *logger {
	return &logger{handler: h}
}
