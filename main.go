package main

import (
	"context"
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var ctx = context.Background()

const NAME = "zunkasrv"

/************************************************************************************************
* Templates
************************************************************************************************/
var runMode string

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
var dbHandytech *sqlx.DB
var dbZunkaFile, dbAldoFile, dbAllnationsFile, dbHandytechFile string

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

// Handytech.
var handytechFilters *HandytechFilters
var handytechSelectedCategories *HandytechSelectedCategories
var handytechSelectedMakers *HandytechSelectedMakers

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

	// Handytech db.
	dbHandytechFile = os.Getenv("HANDYTECH_DB")
	if dbHandytechFile == "" {
		panic("HANDYTECH_DB not defined.")
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

	// Load templates
	loadTemplates()
}

func main() {
	// Log start.
	runMode = "development"
	if production {
		runMode = "production"
	}
	log.Printf("Running in %v mode (version %s)\n", runMode, version)

	// Load allnations data.
	allnationsFilters = LoadAllnationsFilters(path.Join(dataPath, "filters.data"))
	allnationsSelectedCategories = LoadAllnationsSelectedCategories(path.Join(dataPath, "selected_categories.data"))
	allnationsSelectedMakers = LoadAllnationsSelectedMakers(path.Join(dataPath, "selected_makers.data"))

	// Load handytech data.
	handytechFilters = LoadHandytechFilters(path.Join(dataPath, "handytech_filters.data"))
	handytechSelectedCategories = LoadHandytechSelectedCategories(path.Join(dataPath, "handytech_selected_categories.data"))
	handytechSelectedMakers = LoadHandytechSelectedMakers(path.Join(dataPath, "handytech_selected_makers.data"))

	// Start dbs.
	initRedis()
	defer closeRedis()
	initZunkaDB()
	defer closeZunkaDB()
	initAldoDB()
	defer closeAldoDB()
	initAllnationsDB()
	defer closeAllnationsDB()
	initHandytechDB()
	defer closeHandytechDB()

	// Mercado Livre
	initMercadoLivreHandler()

	// Routers
	router := initRouter()

	// start server
	router.ServeFiles("/static/*filepath", http.Dir("./static/"))
	log.Println("listen port", port)
	// Why log.Fall work here?
	// log.Fatal(http.ListenAndServe(":"+port, router))
	log.Fatal(http.ListenAndServe(":"+port, newLogger(router)))
}
