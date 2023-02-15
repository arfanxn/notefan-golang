package config

import (
	"database/sql"

	"github.com/gorilla/mux"
)

type App struct {
	DB     *sql.DB
	Router *mux.Router
}

func NewApp(DB *sql.DB, router *mux.Router) *App {
	return &App{
		DB:     DB,
		Router: router,
	}
}

func InitializeApp() *App {
	// Initialize the environment file
	err := LoadENV()
	if err != nil {
		panic(err)
	}

	// Initialize the Database connection
	db, err := InitializeDB()
	if err != nil {
		panic(err)
	}

	// Instantiate the router
	router := mux.NewRouter()

	// Instantiate the Application
	return NewApp(db, router)
}

func InitializeTestApp() *App {
	// Initialize the environment file
	err := LoadTestENV()
	if err != nil {
		panic(err)
	}

	// Initialize the Database connection
	db, err := InitializeDB()
	if err != nil {
		panic(err)
	}

	// Instantiate the router
	router := mux.NewRouter()

	// Instantiate the Application
	return NewApp(db, router)
}
