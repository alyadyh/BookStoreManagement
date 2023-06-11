package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func AddStaffHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Failed to parse form data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	staffName := r.Form.Get("staffName")
	email := r.Form.Get("email")

	staffID := generateUniqueID(staffName)

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/books_store")
	if err != nil {
		log.Println("Failed to connect to the database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Insert into staff table
	staffStmt, err := db.Prepare("INSERT INTO staff (staff_id, name, email) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("Failed to prepare SQL statement for inserting into staff table:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer staffStmt.Close()

	_, err = staffStmt.Exec(staffID, staffName, email)
	if err != nil {
		log.Println("Failed to execute SQL statement for inserting into authors table:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Staff added successfully!")
}
