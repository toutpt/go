package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Note everywhere we use http://www.dublincore.org/documents/dces/
func onError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func listTables(db *sql.DB) *sql.Rows {
	// SELECT table_schema,table_name
	// FROM information_schema.tables
	// WHERE table_schema='public' ORDER BY table_schema,table_name;
	rows, err := db.Query(`SELECT table_name FROM information_schema.tables
	WHERE table_schema='public' ORDER BY table_name;`)
	onError(err)
	return rows
}

func main() {
	dbname := flag.String("dbname", "", "The name of the database to connect to")
	user := flag.String("user", "", "The user to sign in as")
	// password := flag.String("password", "", "The user's password")
	// host := flag.String("host", "localhost", "The host to connect to. Values that start with / are for unix domain sockets. (default is localhost)")
	// port := flag.String("port", "5432", "The port to bind to. (default is 5432)")
	sslmode := flag.String("sslmode", "disable", "Whether or not to use SSL (default is require, this is not the default for libpq)")
	// fallback := flag.String("fallback_application_name", "", "An application_name to fall back to if one isn't provided.")
	// timeout := flag.String("connect_timeout", "", "Maximum wait for connection, in seconds. Zero or not specified means wait indefinitely.")
	// sslcert := flag.String("sslcert", "", "Cert file location. The file must contain PEM encoded data.")
	// sslkey := flag.String("sslkey", "", "Key file location. The file must contain PEM encoded data.")
	// sslrootcert := flag.String("sslrootcert", "", "The location of the root certificate file. The file must contain PEM encoded data.")
	flag.Parse()

	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=%s", *user, *dbname, *sslmode)
	db, err := sql.Open("postgres", connStr)
	onError(err)
	tables := listTables(db)
	fmt.Println(tables)
}
