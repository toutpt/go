package models

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// Field define a field in a model to embed data
type Field struct {
	ID          string
	Title       string
	Description string
	Type        string

	DefaultJSON string
	Required    bool
	Widget      string
	Searchable  string
}

func (f *Field) String() string {
	return fmt.Sprintf("Field-%s", f.Title)
}

// Decode read the body as JSON
func (f *Field) Decode(r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

// Marshal serialize the model into json
func (f *Field) Marshal() []byte {
	data, err := json.Marshal(f)
	onError(err)
	return data
}

// Insert add the new item
func (f *Field) Insert(db *pg.DB) (orm.Result, error) {
	return db.Model(f).Insert()
}

// HandleFields call the right action from the web
func HandleFields(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		instance := &Field{}
		err := db.Model(instance).Select()
		onError(err)
		// handle error
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(instance.Marshal())
	} else if r.Method == "POST" {
		instance := &Field{}
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
