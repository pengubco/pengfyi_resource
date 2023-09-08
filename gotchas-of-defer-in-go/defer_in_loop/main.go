package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgres://postgres:password@localhost:15432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Error opening database connection:", err)
		return
	}
	defer db.Close()

	foo(db, 5)

	fooFixed(db, 5)
}

func foo(db *sql.DB, n int) {
	for i := 0; i < n; i++ {
		fmt.Printf("query %d, ", i)
		rows, err := db.Query("SELECT 1;")
		if err != nil {
			fmt.Println("Error querying the database:", err)
			return
		} else {
			fmt.Println("success")
		}
		defer rows.Close()
	}
}

func fooFixed(db *sql.DB, n int) {
	for i := 0; i < n; i++ {
		func() {
			fmt.Printf("query %d, ", i)
			rows, err := db.Query("SELECT 1;")
			if err != nil {
				fmt.Println("Error querying the database:", err)
				return
			} else {
				fmt.Println("success")
			}
			defer rows.Close()
		}()
	}
}
