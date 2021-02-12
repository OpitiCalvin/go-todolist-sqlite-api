package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (a *App) GetItemByID(Id int) bool {
	todo := &TodoItemModel{}
	result := a.DB.First(&todo, Id)
	if result.Error != nil {
		log.Warn("TodoItem not found in database")
		return false
	}
	return true
}

func (a *App) GetTodoItems(completed bool) interface{} {
	var todos []TodoItemModel
	TodoItems := a.DB.Where("completed = ?", completed).Find(&todos).Value
	return TodoItems
}

func (a *App) Healthz(w http.ResponseWriter, r *http.Request) {
	log.Info("API health is OK")
	w.Header().Set("Content-Type", "application-json")
	io.WriteString(w, `{"alive": true}`)
}

func (a *App) CreateItem(w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	log.WithFields(log.Fields{"description": description}).Info("Add new TodoItem. Saving to database.")
	todo := &TodoItemModel{Description: description, Completed: false}
	a.DB.Create(&todo)
	result := a.DB.Last(&todo)
	w.Header().Set("Content-Type", "application-json")
	json.NewEncoder(w).Encode(result.Value)
}

func (a *App) UpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	err := a.GetItemByID(id)
	if err == false {
		w.Header().Set("Content-Type", "application.json")
		io.WriteString(w, `{"updated": false, "error": "Record Not Found"}`)
	} else {
		completed, _ := strconv.ParseBool(r.FormValue("completed"))
		log.WithFields(log.Fields{"Id": id, "Completed": completed}).Info("Updating TodoItem")
		todo := &TodoItemModel{}
		a.DB.First(&todo, id)
		todo.Completed = completed
		a.DB.Save(&todo)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"updated": true, "error": false}`)
	}
}

func (a *App) DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	err := a.GetItemByID(id)
	if err == false {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"deleted": false, "error": "Record Not Found"}`)
	} else {
		log.WithFields(log.Fields{"Id": id}).Info("Deleting TodoItem")
		todo := &TodoItemModel{}
		a.DB.First(&todo, id)
		a.DB.Delete(&todo)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"deleted": true}`)
	}
}

func (a *App) GetCompletedItems(w http.ResponseWriter, r *http.Request) {
	log.Info("Get completed TodoItems")
	completedTodoItems := a.GetTodoItems(true)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(completedTodoItems)
}

func (a *App) GetIncompleteItems(w http.ResponseWriter, r *http.Request) {
	log.Info("Get Incomplete TodoItems")
	incompleteTodoItems := a.GetTodoItems(false)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(incompleteTodoItems)
}
