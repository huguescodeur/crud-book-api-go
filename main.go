package main

import (
	"crud-book/internal/service"
	"crud-book/internal/store"
	"crud-book/internal/transport"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var schema = `
	CREATE TABLE IF NOT EXISTS books(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT NOT NULL
	);
`

func main() {
	//? Config database
	dbName := "book_db.db"
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal("error open db: ", err)
	}
	defer func() {
		fmt.Println("Closing database...")

		if err := db.Close(); err != nil {
			log.Fatal("error closing: ", err)
		}
	}()

	if err := db.Ping(); err != nil {
		log.Fatal("ping error: ", err)
	}

	fmt.Println("Database connection establish")

	// if _, err := db.Exec(schema); err != nil {
	// 	log.Fatal("error create table: ", err)
	// }

	// fmt.Println("Table was created")

	//? injection dependance
	bookStore := store.NewStore(db)
	bookService := service.NewService(bookStore)
	bookHandler := transport.NewHandler(bookService)

	//? Routing
	http.HandleFunc("/book", bookHandler.HandlerBook)
	http.HandleFunc("/book/{id}", bookHandler.HandlerBookByID)

	//? Server
	addr := ":8080"

	fmt.Printf("Server start on http://localhost%s\n", addr)

	http.ListenAndServe(addr, nil)
}
