package models

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// Model to design a database
type Model struct {
	ID          string
	Title       string
	Description string
}

func (m *Model) String() string {
	return fmt.Sprintf("Model-%s", m.Title)
}
func (m *Model) Decode(r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

// Marshal serialize the model into json
func (m *Model) Marshal() []byte {
	data, err := json.Marshal(m)
	onError(err)
	return data
}

// Insert add the new item
func (m *Model) Insert(db *pg.DB) (orm.Result, error) {
	return db.Model(m).Insert()
}

// HandleModels call the right action from the web
func HandleModels(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		instance := &Model{}
		err := db.Model(instance).Select()
		onError(err)
		// handle error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(instance.Marshal())
	} else if r.Method == "POST" {
		instance := &Model{}
		instance.Decode(r)
		if instance.ID != "" {
			_, err := db.Model(instance).Insert()
			onError(err)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "{ \"error\": \"%s\" }", err.Error())
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write(instance.Marshal())
			}
			w.Header().Set("Content-Type", "application/json")
		}
	} else {
		fmt.Println("Method not supported:", r.Method)
		http.NotFound(w, r)
	}
}
