package main

import (
	"go-twt/internal/db"
	"go-twt/internal/env"
	"go-twt/internal/store"
	"log"

	"github.com/gorilla/sessions"
)

var sessionStore *sessions.CookieStore

type application struct {
	config   config
	store    store.Storage
	sessions *sessions.CookieStore
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func main() {
	// Sessions stuff

	// We need to generate a random key for prod
	sessionStore = sessions.NewCookieStore([]byte("super-secret-key"))
	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 1,
		HttpOnly: true,
		Secure:   false,
	}

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	dbConn, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}
	defer dbConn.Close()
	log.Println("DB Connected - connection pool established!")

	store := store.NewStorage(dbConn)

	app := &application{
		config:   cfg,
		store:    store,
		sessions: sessionStore,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
