package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

// Init redis db.
func initRedis() {
	// Connect to Redis DB.
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := redisClient.Ping(ctx).Result()

	if err != nil || pong != "PONG" {
		log.Panicf("[panic] Couldn't connect to Redis DB. %s", err)
	}
	// log.Printf("Connected to Redis")
}

// Close redis db.
func closeRedis() {
	// log.Printf("Closing Redis connection...")
}

// Init zunka db.
func initZunkaDB() {
	// Start data base.
	dbZunka, err = sql.Open("sqlite3", dbZunkaFile)
	if err != nil {
		log.Fatal(fmt.Errorf("Error open db %v: %v", dbZunka, err))
	}
	err = dbZunka.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

// Close zunka db.
func closeZunkaDB() {
	dbZunka.Close()
}

// Init aldo db.
func initAldoDB() {
	dbAldo = sqlx.MustConnect("sqlite3", dbAldoFile)
}

// Close aldo db.
func closeAldoDB() {
	dbAldo.Close()
}

// Init allnations db.
func initAllnationsDB() {
	dbAllnations = sqlx.MustConnect("sqlite3", dbAllnationsFile)
}

// Close allnations db.
func closeAllnationsDB() {
	dbAllnations.Close()
}

// Init handytech db.
func initHandytechDB() {
	dbHandytech = sqlx.MustConnect("sqlite3", dbHandytechFile)
}

// Close handytech db.
func closeHandytechDB() {
	dbHandytech.Close()
}

// Init motospeed db.
func initMotospeedDB() {
	dbMotospeed = sqlx.MustConnect("sqlite3", dbMotospeedFile)
}

// Close motospeed db.
func closeMotospeedDB() {
	dbMotospeed.Close()
}
