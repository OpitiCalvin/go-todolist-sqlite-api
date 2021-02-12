package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

var err error

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (a *App) Initialize() {
	a.Router = mux.NewRouter()

	// check if sqlite db exists and if not, create one
	_, err := os.Stat("./todolist.sqlite3")
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create("./todolist.sqlite3")
			if err != nil {
				log.Fatal(err.Error())
			}
			file.Close()
			log.Info("todolist.sqlite3 created successfully")
		}
	}

	a.DB, err = gorm.Open("sqlite3", "./todolist.sqlite3")
	if err != nil {
		panic("Failed to connect to database")
	}

	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	http.Handle("/favicon.ico", http.NotFoundHandler())
	a.Router.HandleFunc("/healthz", a.Healthz).Methods("GET")
	a.Router.HandleFunc("/todo-completed", a.GetCompletedItems).Methods("GET")
	a.Router.HandleFunc("/todo-incomplete", a.GetIncompleteItems).Methods("GET")
	a.Router.HandleFunc("/todo", a.CreateItem).Methods("POST")
	a.Router.HandleFunc("/todo/{id}", a.UpdateItem).Methods("PUT")
	a.Router.HandleFunc("/todo/{id}", a.DeleteItem).Methods("DELETE")
}

func (a *App) Run(addr string) {
	handler := cors.Default().Handler(a.Router)

	log.Info("Starting Todolist API servier")
	log.Fatal(http.ListenAndServe(addr, handler))
}
