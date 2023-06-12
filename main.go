package main

import (
	"log"
	"net/http"

	"github.com/book-crud/handlers"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/addbook", handlers.AddBook)
	http.HandleFunc("/deletebook", handlers.DeleteBook)
	http.HandleFunc("/addstaff", handlers.AddStaff)
	http.HandleFunc("/deletestaff", handlers.DeleteStaff)
	http.HandleFunc("/addtransaction", handlers.AddTransaction)
	http.HandleFunc("/deletetransaction", handlers.DeleteTransaction)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
