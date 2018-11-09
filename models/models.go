package models

import (
	"fmt"

	"github.com/go-pg/pg"
)

func onError(err error) {
	if err != nil {
		// log.Fatal(err)
		fmt.Println(err)
	}
}

var db *pg.DB

//InjectDB expose the DB in the scope of the module
func InjectDB(pgdb *pg.DB) {
	db = pgdb
}
