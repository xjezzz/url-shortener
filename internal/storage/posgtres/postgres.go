package posgtres

import (
	"database/sql"
	"log"
)

type Storage struct {
	db
}

func New() {
	db, err := sql.Open("postgres", "user=username password=password host=localhost dbname=mydb sslmode=disable")
	if err != nil {
		log.Fatalf("Error: Unable to connect to database: %v", err)
	}
	defer db.Close()

}
