package actions

import (
	"fmt"
	"net/http"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

var pgdb *pg.DB

// ActionArgs represent arguments passed to the action
type ActionArgs struct {
	Body interface{}
}

// Action structure can be created from request
type Action struct {
	name  string
	model string
	args  *ActionArgs
	call  func(*pg.DB, *ActionArgs) error
}

func (a *Action) String() string {
	return fmt.Sprintf("Action-%s", a.name)
}

// Call an action
func (a *Action) Call() error {
	return a.call(pgdb, a.args)
}

// ActionFromReq create a new
func ActionFromReq(r *http.Request) *Action {
	action := &Action{}
	action.name = r.URL.Path[1:]
	//todo: parse body to arguments
	action.args = &ActionArgs{}
	if action.name == "create" {
		action.call = Create
	}
	return action
}

// CreateTable add a table in the database
func CreateTable(db *pg.DB, model interface{}) error {
	err := db.CreateTable(model, &orm.CreateTableOptions{
		FKConstraints: true,
	})
	return err
}

// Create create an entry in the db
func Create(db *pg.DB, args *ActionArgs) error {
	// return db.Create
	return nil
}

// Init initialize the action module
func Init(db *pg.DB) {
	pgdb = db
}

// HandleAction call the right action from the web
func HandleAction(w http.ResponseWriter, r *http.Request) {
	action := ActionFromReq(r)
	err := action.Call()
	if err != nil {
		// write in the response status 500
	}
	// response the result
	fmt.Println("return response")
	fmt.Fprintf(w, "{}")
}
