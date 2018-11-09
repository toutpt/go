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

// ActionFunc represent the go func actions
type ActionFunc func(*pg.DB, *ActionArgs) (orm.Result, error)

// ActionDef is the metadata to get the Action
type ActionDef struct {
	name  string
	model string
}

// Action structure can be created from request
type Action struct {
	name  string
	model string
	args  *ActionArgs
	call  ActionFunc
}

func (a *Action) String() string {
	return fmt.Sprintf("Action-%s", a.name)
}

func getActionFunc(name string) func(*pg.DB, *ActionArgs) (orm.Result, error) {
	if name == "create" {
		return Create
	}
	return nil
}

// NewAction create a new instance of action and resolve the func using the name
func NewAction(name string, model string) *Action {
	action := &Action{
		name:  name,
		model: model,
	}
	action.call = getActionFunc(name)
	return action
}

// SetArgs an action
func (a *Action) SetArgs(args *ActionArgs) {
	a.args = args
}

// Call an action
func (a *Action) Call() (orm.Result, error) {
	fmt.Println(a.call)
	if a.call == nil {
		return nil, fmt.Errorf("no call found in %s", a.name)
	}
	return a.call(pgdb, a.args)
}

// ActionFromReq create a new
func ActionFromReq(r *http.Request) *Action {
	action := &Action{}
	action.name = r.URL.Path[1:]
	//todo: parse body to arguments
	action.args = &ActionArgs{}
	action.call = getActionFunc(action.name)
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
func Create(db *pg.DB, args *ActionArgs) (orm.Result, error) {
	// return db.Create
	return db.Model(args.Body).Insert()
}

// Init initialize the action module
func Init(db *pg.DB) {
	pgdb = db
}

// HandleAction call the right action from the web
func HandleAction(w http.ResponseWriter, r *http.Request) {
	action := ActionFromReq(r)
	_, err := action.Call()
	if err != nil {
		// write in the response status 500
	}
	// response the result
	fmt.Println("return response")
	fmt.Fprintf(w, "{}")
}
