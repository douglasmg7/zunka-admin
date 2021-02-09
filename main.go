package main

import (
	"context"
	"database/sql"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

var ctx = context.Background()

const NAME = "zunkasrv"

/************************************************************************************************
* Templates
************************************************************************************************/
// Geral.
var tmplMaster, tmplIndex, tmplDeniedAccess *template.Template

// Misc.
var tmplMessage *template.Template
var tmplChangelog *template.Template
var tmplTest *template.Template

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
var tmplAllnationsProducts, tmplAllnationsProduct, tmplAllnationsFilters, tmplAllnationsCategories, tmplAllnationsMakers *template.Template

// Mercado Livre.
var tmplMercadoLivreAuthUser, tmplMercadoLivreText *template.Template

// Auth.
var tmplAuthSignup, tmplAuthSignin, tmplPasswordRecovery, tmplPasswordReset *template.Template

// Student.
var tmplStudent, tmplAllStudent, tmplNewStudent *template.Template

var production bool
var port string

// Log.
var Trace *log.Logger
var Debug *log.Logger
var Info *log.Logger
var Warn *log.Logger
var Error *log.Logger

// Db.
var redisClient *redis.Client
var dbZunka *sql.DB
var dbAldo *sqlx.DB
var dbAllnations *sqlx.DB
var dbZunkaFile, dbAldoFile, dbAllnationsFile string

var zunkaPath string
var GS string

var err error

// Data path.
var dataPath string

// Sessions from each user.
var sessions = Sessions{
	mapUserID:      map[string]int{},
	mapSessionData: map[int]*SessionData{},
}

// Allnations.
var allnationsFilters *AllnationsFilters
var allnationsSelectedCategories *AllnationsSelectedCategories
var allnationsSelectedMakers *AllnationsSelectedMakers

func init() {
	// Check if production mode.
	if os.Getenv("RUN_MODE") == "production" {
		production = true
	}

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

	// Log path.
	logPath := path.Join(zunkaPath, "log", "zunkasrv")
	os.MkdirAll(logPath, os.ModePerm)

	// Data path.
	dataPath = path.Join(zunkaPath, "data", "zunkasrv")
	os.MkdirAll(dataPath, os.ModePerm)

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

	// Allnations db.
	dbAllnationsFile = os.Getenv("ALLNATIONS_DB")
	if dbAllnationsFile == "" {
		panic("ALLNATIONS_DB not defined.")
	}

	// Dev mode.
	if !production {
		dbAllnationsFile += "-dev"
	}

	// Log file.
	logFile, err := os.OpenFile(path.Join(logPath, "zunkasrv.log"), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	// Log configuration.
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.SetPrefix("[zunkasrv] ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lmsgprefix)
	// log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile | log.Lmsgprefix)
	// log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lmsgprefix)
	// log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	// log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	// log.SetFlags(log.LstdFlags | log.Ldate | log.Lshortfile)
	// log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	Trace = log.New(mw, "["+NAME+"] [trace] ", log.Ldate|log.Lmicroseconds|log.Lmsgprefix)
	Debug = log.New(mw, "["+NAME+"] [debug] ", log.Ldate|log.Lmicroseconds|log.Lmsgprefix)
	Info = log.New(mw, "["+NAME+"] [info ] ", log.Ldate|log.Lmicroseconds|log.Lmsgprefix)
	Warn = log.New(mw, "["+NAME+"] [warn ] ", log.Ldate|log.Lmicroseconds|log.Lmsgprefix)
	Error = log.New(mw, "["+NAME+"] [error] ", log.Ldate|log.Lmicroseconds|log.Lmsgprefix|log.Lshortfile)

	/************************************************************************************************
	* Load templates
	************************************************************************************************/
	// Geral.
	tmplMaster = template.Must(template.ParseGlob("templates/master/*"))
	tmplIndex = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/index.tpl"))
	tmplDeniedAccess = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/misc/deniedAccess.tpl"))
	// Misc.
	tmplMessage = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/misc/message.tpl"))
	tmplChangelog = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/misc/changelog.gohtml"))
	tmplTest = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/misc/test.gohtml"))
	// User.
	tmplUserAdd = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/userAdd.tpl"))
	tmplUserAccount = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/user/userAccount.tpl"))
	tmplUserChangeName = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/user/userChangeName.tpl"))
	tmplUserChangeEmail = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/user/userChangeEmail.tpl"))
	tmplUserChangeMobile = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/user/userChangeMobile.tpl"))
	// Aldo.
	tmplAldoProducts = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/aldo/aldoProducts.tmpl"))
	tmplAldoProduct = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/aldo/aldoProduct.tmpl"))
	tmplAldoCategories = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/aldo/aldoCategories.tmpl"))
	// Allnations.
	tmplAllnationsProducts = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/allnations/allnationsProducts.tmpl"))
	tmplAllnationsProduct = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/allnations/allnationsProduct.tmpl"))
	tmplAllnationsFilters = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/allnations/allnationsFilters.tmpl"))
	tmplAllnationsCategories = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/allnations/allnationsCategories.tmpl"))
	tmplAllnationsMakers = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/allnations/allnationsMakers.gohtml"))
	// Mercado Livre.
	tmplMercadoLivreAuthUser = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/mercado_livre/mercadoLivreAuthUser.gohtml"))
	tmplMercadoLivreText = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/mercado_livre/mercadoLivreText.gohtml"))

	// Auth.
	tmplAuthSignup = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/auth/signup.tpl"))
	tmplAuthSignin = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/auth/signin.tpl"))
	tmplPasswordRecovery = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/auth/passwordRecovery.tpl"))
	tmplPasswordReset = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/auth/passwordReset.tpl"))
	// Student.
	tmplStudent = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/student.tpl"))
	tmplAllStudent = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/allStudent.tpl"))
	tmplNewStudent = template.Must(template.Must(tmplMaster.Clone()).ParseFiles("templates/newStudent.tpl"))
}

func main() {
	// Log start.
	runMode := "development"
	if production {
		runMode = "production"
	}
	log.Printf("Running in %v mode (version %s)\n", runMode, version)

	// Load allnations data.
	allnationsFilters = LoadAllnationsFilters(path.Join(dataPath, "filters.data"))
	allnationsSelectedCategories = LoadAllnationsSelectedCategories(path.Join(dataPath, "selected_categories.data"))
	allnationsSelectedMakers = LoadAllnationsSelectedMakers(path.Join(dataPath, "selected_makers.data"))

	// Start dbs.
	initRedis()
	defer closeRedis()
	initZunkaDB()
	defer closeZunkaDB()
	initAldoDB()
	defer closeAldoDB()
	initAllnationsDB()
	defer closeAllnationsDB()

	// Mercado Livre.
	initMercadoLivreHandler()

	// Init router.
	router := httprouter.New()
	router.GET("/favicon.ico", faviconHandler)
	router.GET("/ns/favicon.ico", faviconHandler)
	router.GET("/", getSession(indexHandler))
	router.GET("/ns/", getSession(indexHandler))
	router.GET("/ping", getSession(indexPing))
	router.GET("/ns/ping", getSession(indexPing))

	// Clean the session cache.
	router.GET("/ns/clean-sessions", checkPermission(cleanSessionsHandler, "admin"))
	// Changelog page.
	router.GET("/ns/changelog", checkPermission(changelogHandler, "admin"))
	// Test.
	router.GET("/ns/test", checkPermission(testPageHandler, "admin"))
	router.POST("/ns/test/send-email", checkPermission(testSendMailPost, "admin"))

	// Aldo.
	// Products list page.
	router.GET("/ns/aldo/products", checkPermission(aldoProductsHandler, "read"))
	// Product page.
	router.GET("/ns/aldo/product/:code", checkPermission(aldoProductHandler, "read"))
	// Create product on zunka server.
	router.POST("/ns/aldo/product/:code", checkPermission(aldoProductHandlerPost, "write"))
	// Check product change.
	router.POST("/ns/aldo/product/:code/checked", checkPermission(aldoProductCheckedHandlerPost, "write"))
	// Product removed from site, so remove his reference from the site system.
	router.DELETE("/ns/aldo/product/mongodb_id/:code", checkApiAuthorization(aldoProductMongodbIdHandlerDelete))
	// Categories page.
	router.GET("/ns/aldo/categories", checkPermission(aldoCategoriesHandler, "read"))
	// Save categories.
	router.POST("/ns/aldo/categories", checkPermission(aldoCategoriesHandlerPost, "write"))

	// Allnations.
	// Products list page.
	router.GET("/ns/allnations/products", checkPermission(allnationsProductsHandler, "read"))
	// Product page.
	router.GET("/ns/allnations/product/:code", checkPermission(allnationsProductHandler, "read"))
	// Create product on zunka server.
	router.POST("/ns/allnations/product/:code", checkPermission(allnationsProductHandlerPost, "write"))
	// Check product change.
	router.POST("/ns/allnations/product/:code/checked", checkPermission(allnationsProductCheckedHandlerPost, "write"))
	// Product removed from site, so remove his reference from zunkasrv.
	router.DELETE("/ns/allnations/product/zunka_product_id/:code", checkApiAuthorization(allnationsProductZunkaProductIdHandlerDelete))
	// Filter page.
	router.GET("/ns/allnations/filters", checkPermission(allnationsFiltersHandler, "read"))
	// Save filter.
	router.POST("/ns/allnations/filters", checkPermission(allnationsFiltersHandlerPost, "write"))
	// Categories page.
	router.GET("/ns/allnations/categories", checkPermission(allnationsCategoriesHandler, "read"))
	// Save categories.
	router.POST("/ns/allnations/categories", checkPermission(allnationsCategoriesHandlerPost, "write"))
	// Makers page.
	router.GET("/ns/allnations/makers", checkPermission(allnationsMakersHandler, "read"))
	// Save categories.
	router.POST("/ns/allnations/makers", checkPermission(allnationsMakersHandlerPost, "write"))

	// Mercado Livre
	// Login
	router.GET("/ns/ml/auth/login", checkPermission(mercadoLivreAuthLoginHandler, "write"))
	// User code
	router.GET("/ns/ml/auth/user", mercadoLivreAuthUserHandler)
	// Notification
	router.GET("/ns/ml/notification", mercadoLivreNotificationHandler)
	// Load user code from zunka server. Used by zunka server in development.
	router.GET("/ns/ml/user/load-code", checkPermission(mercadoLivreLoadUserCode, "read"))
	// Show user code
	router.GET("/ns/ml/user/code", checkPermission(mercadoLivreUserCodeHandler, "read"))
	router.GET("/ns/ml/api/user-code", checkApiAuthorization(mercadoLivreAPIUserCodeHandler))
	// Products
	router.GET("/ns/ml/users/me", checkPermission(mercadoLivreUsersMeHandler, "read"))

	// Autheticate user.
	// router.GET("/ns/ml/auth/login", checkPermission(mercadoLivreAuthUserHandler, "read"))
	// router.POST("/ns/ml/auth/login", checkPermission(mercadoLivreAuthUserHandlerPost, "write"))

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
	// log.Printf("%s %s - begin", req.Method, req.URL.Path)
	start := time.Now()
	l.handler.ServeHTTP(w, req)
	log.Printf("%s %s %v", req.Method, req.URL.Path, time.Since(start))
	// log.Printf("header: %v", req.Header)
}

// New logger.
func newLogger(h http.Handler) *logger {
	return &logger{handler: h}
}

/**************************************************************************************************
* Error
**************************************************************************************************/
func checkError(err error) bool {
	if err != nil {
		// notice that we're using 1, so it will actually log where
		// the error happened, 0 = this function, we don't want that.
		_, file, line, _ := runtime.Caller(1)
		log.Printf("[error] %s:%d: %v", filepath.Base(file), line, err)
		// function, file, line, _ := runtime.Caller(1)
		// log.Printf("[error] [%s] [%s:%d] %v", filepath.Base(file), runtime.FuncForPC(function).Name(), line, err)
		return true
	}
	return false
}
