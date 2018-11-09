package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-pg/pg"
	"github.com/toutpt/go/actions"
)

// Note everywhere we use http://www.dublincore.org/documents/dces/
func onError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func isTableExists(db *pg.DB, model interface{}) bool {
	fmt.Print("isTableExists: ")
	err := db.Select(model)
	if err != nil {
		fmt.Println("nop", err)
		return false
	}
	fmt.Println("yep")
	return true
}

// Model to design a database
type Model struct {
	ID          string
	Title       string
	Description string
}

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

func main() {
	dbFlag := flag.String("db", "", "[*] connection string to postgresql")
	portFlag := flag.String("port", "8081", "[*] port the webapp will listen")
	flag.Parse()
	if *dbFlag == "" {
		onError(errors.New("no db provided"))
	}

	dbOptions, err := pg.ParseURL(*dbFlag)
	onError(err)
	fmt.Println("- connect to postgresql")
	db := pg.Connect(dbOptions)
	fmt.Println(db)
	defer db.Close()
	actions.Init(db)
	if !isTableExists(db, &Model{ID: "0"}) {
		onError(actions.CreateTable(db, &Model{}))
		action := actions.NewAction("create", "model")
		args := actions.ActionArgs{Body: &Model{ID: "0", Title: "App"}}
		action.SetArgs(&args)
		_, err = action.Call()
		onError(err)
	}
	if !isTableExists(db, &Field{ID: "0"}) {
		onError(actions.CreateTable(db, &Field{}))
		action := actions.NewAction("create", "field")
		args := actions.ActionArgs{Body: &Field{ID: "0", Title: "version"}}
		action.SetArgs(&args)
		_, err = action.Call()
		onError(err)
	}
	http.HandleFunc("/action", actions.HandleAction)
	onError(http.ListenAndServe(":"+*portFlag, http.FileServer(http.Dir("./elm-app/dist"))))
}
