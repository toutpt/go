package actions

import (
	"fmt"

	"github.com/go-pg/pg"
)

// Select create an entry in the db
func Select(db *pg.DB, model interface{}, args *ActionArgs) *Result {
	fmt.Println("call actions.Select !!", model, args)
	err := db.Select(args.Body)
	return &Result{
		err: err,
	}
}
