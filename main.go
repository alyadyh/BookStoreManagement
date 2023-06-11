package main

import (
	"log"
	"net/http"

	"github.com/book-crud/handlers"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/addbook", handlers.AddBookHandler)
	http.HandleFunc("/deletebook", handlers.DeleteBookHandler)
	http.HandleFunc("/addstaff", handlers.AddStaffHandler)
	http.HandleFunc("/addtransaction", handlers.AddTransactionHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
