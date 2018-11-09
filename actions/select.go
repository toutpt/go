package actions

import (
	"fmt"
	"reflect"

	"github.com/go-pg/pg"
)

// Select create an entry in the db
func Select(db *pg.DB, model reflect.Type, args *ActionArgs) *Result {
	fmt.Println("call actions.Select !!", model, args)
	data := reflect.ValueOf(reflect.New(model))
	err := db.Model(model).Select()
	fmt.Println("db.Select(data) ", data)
	return &Result{
		err: err,
	}
}
