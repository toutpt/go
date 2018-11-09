package actions

import (
	"reflect"

	"github.com/go-pg/pg"
)

// Create create an entry in the db
func Create(db *pg.DB, model reflect.Type, args *ActionArgs) *Result {
	// return db.Create
	_, err := db.Model(args.Body).Insert()
	return &Result{
		err: err,
	}
}
