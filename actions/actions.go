package actions

import (
	"fmt"
	"net/http"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/toutpt/go/models"
)

var pgdb *pg.DB
var actions map[string]ActionFunc

func initActions() {
	actions = make(map[string]ActionFunc)
	actions["create"] = Create
	actions["select"] = Select
}

// ActionArgs represent arguments passed to the action
type ActionArgs struct {
	Body  interface{}
	Query map[string][]string
	model interface{}
}

// Result is an action result
type Result struct {
	data   map[string]interface{}
	status int
	err    error
}

// ActionFunc represent the go func actions
type ActionFunc func(*pg.DB, interface{}, *ActionArgs) *Result

// ActionDef is the metadata to get the Action
type ActionDef struct {
	name  string
	model string
}

// Action structure can be created from request
type Action struct {
	name  string
	model interface{}
	args  *ActionArgs
	call  ActionFunc
}

func (a *Action) String() string {
	return fmt.Sprintf("Action-%s", a.name)
}

func getActionFunc(name string) ActionFunc {
	if len(actions) == 0 {
		initActions()
	}
	return actions[name]
}

func getActionModel(name string) interface{} {
	m := models.Get()
	return m[name]
}

// NewAction create a new instance of action and resolve the func using the name
func NewAction(name string, model string) *Action {
	action := &Action{
		name:  name,
		model: models.Get()[model],
	}
	action.call = getActionFunc(name)
	return action
}

// SetArgs an action
func (a *Action) SetArgs(args *ActionArgs) {
	a.args = args
}

// Call an action
func (a *Action) Call() *Result {
	fmt.Println(a.call)
	if a.call == nil {
		return &Result{
			err: fmt.Errorf("no call found in %s", a.name),
		}
	}
	return a.call(pgdb, a.model, a.args)
}

// ActionFromReq create a new
func ActionFromReq(r *http.Request) *Action {
	query := r.URL.Query()
	action := &Action{}
	if len(query["id"]) == 0 {
		return nil
	}
	action.name = query["id"][0]
	if len(query["model"]) == 1 {
		model := query["model"][0]
		if model != "" {
			action.model = models.Get()[model]
		}
	}
	//todo: parse body to arguments
	action.args = &ActionArgs{
		Query: r.URL.Query(),
		Body:  r.Body,
	}
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

// Init initialize the action module
func Init(db *pg.DB) {
	pgdb = db
}

// HandleAction call the right action from the web
func HandleAction(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HandleAction()", r.URL.Query())

	action := ActionFromReq(r)
	if action == nil {
		fmt.Println("action not found return 404")
		http.NotFound(w, r)
	}
	result := action.Call()
	if result.err != nil {
		http.Error(w, result.err.Error(), http.StatusInternalServerError)
	}
	fmt.Println("return empty response")
	fmt.Fprintf(w, "{}")
}
