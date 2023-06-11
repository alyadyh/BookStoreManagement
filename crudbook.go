package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/addbook", addBookHandler)
	http.HandleFunc("/deletebook", deleteBookHandler)
	http.HandleFunc("/addstaff", addStaffHandler)
	http.HandleFunc("/addtransaction", addTransactionHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func addBookHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Failed to parse form data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	bookTitle := r.Form.Get("bookTitle")
	authorName := r.Form.Get("bookAuthor")
	publisherName := r.Form.Get("bookPublisher")
	publicationDate := r.Form.Get("publicationDate")
	ISBN := r.Form.Get("isbn")
	price := r.Form.Get("price")
	stockQty := r.Form.Get("stockQuantity")

	bookID := generateUniqueID(bookTitle)
	authorID := generateUniqueID(authorName)
	publisherID := generateUniqueID(publisherName)

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/books_store")
	if err != nil {
		log.Println("Failed to connect to the database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Insert into authors table
	authorStmt, err := db.Prepare("INSERT INTO authors (author_id, name) VALUES (?, ?)")
	if err != nil {
		log.Println("Failed to prepare SQL statement for inserting into authors table:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer authorStmt.Close()

	_, err = authorStmt.Exec(authorID, authorName)
	if err != nil {
		log.Println("Failed to execute SQL statement for inserting into authors table:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Insert into publishers table
	publisherStmt, err := db.Prepare("INSERT INTO publishers (publisher_id, name) VALUES (?, ?)")
	if err != nil {
		log.Println("Failed to prepare SQL statement for inserting into publishers table:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer publisherStmt.Close()

	_, err = publisherStmt.Exec(publisherID, publisherName)
	if err != nil {
		log.Println("Failed to execute SQL statement for inserting into publishers table:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Insert into books table
	bookStmt, err := db.Prepare("INSERT INTO books (book_id, title, author_id, publisher_id, publication_date, ISBN, price, stock_qty) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Failed to prepare SQL statement for inserting into books table:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer bookStmt.Close()

	_, err = bookStmt.Exec(bookID, bookTitle, authorID, publisherID, publicationDate, ISBN, price, stockQty)
	if err != nil {
		log.Println("Failed to execute SQL statement for inserting into books table:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Book added successfully!")
}

func deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Failed to parse form data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	bookTitle := r.Form.Get("bookTitle")

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/books_store")
	if err != nil {
		log.Println("Failed to connect to the database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM books WHERE title = ?")
	if err != nil {
		log.Println("Failed to prepare SQL statement:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(bookTitle)
	if err != nil {
		log.Println("Failed to execute SQL statement:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Book deleted successfully!")
}

// ----- STAFF ----- //
func addStaffHandler(w http.ResponseWriter, r *http.Request) {
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

// ----- CART TRANSACTION ----- //

func addTransactionHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Failed to parse form data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	bookTitle := r.Form.Get("bookTitle")
	quantityStr := r.Form.Get("quantity")
	staffName := r.Form.Get("staffName")

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		log.Println("Invalid quantity:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/books_store")
	if err != nil {
		log.Println("Failed to connect to the database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Retrieve staff details from the staff table
	staffQuery := "SELECT staff_id FROM staff WHERE name = ?"
	staffStmt, err := db.Prepare(staffQuery)
	if err != nil {
		log.Println("Failed to prepare SQL statement for retrieving staff details:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer staffStmt.Close()
	var staffID int

	err = staffStmt.QueryRow(staffName).Scan(&staffID)
	if err != nil {
		log.Println("Failed to execute SQL statement for retrieving book details:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Retrieve book details from the books table
	bookQuery := "SELECT book_id, price, stock_qty FROM books WHERE title = ?"
	bookStmt, err := db.Prepare(bookQuery)
	if err != nil {
		log.Println("Failed to prepare SQL statement for retrieving book details:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer bookStmt.Close()

	var bookID int
	var itemPrice float64
	var stockQty int
	err = bookStmt.QueryRow(bookTitle).Scan(&bookID, &itemPrice, &stockQty)
	if err != nil {
		log.Println("Failed to execute SQL statement for retrieving book details:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if quantity > stockQty {
		log.Println("Insufficient stock quantity")
		http.Error(w, "Insufficient stock quantity", http.StatusBadRequest)
		return
	}

	// Update stock quantity in books table
	updateQtyStmt, err := db.Prepare("UPDATE books SET stock_qty = stock_qty - ? WHERE book_id = ?")
	if err != nil {
		log.Println("Failed to prepare SQL statement for updating stock quantity:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer updateQtyStmt.Close()

	_, err = updateQtyStmt.Exec(quantity, bookID)
	if err != nil {
		log.Println("Failed to execute SQL statement for updating stock quantity:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Insert into cart table
	cartInsertStmt, err := db.Prepare("INSERT INTO cart (cart_id, book_id, qty, item_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Println("Failed to prepare SQL statement for inserting into cart table:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer cartInsertStmt.Close()

	// Generate a unique cart ID
	cartID := generateUniqueID(bookTitle + staffName)

	_, err = cartInsertStmt.Exec(cartID, bookID, quantity, itemPrice)
	if err != nil {
		log.Println("Failed to execute SQL statement for inserting into cart table:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Calculate the total price
	totalPrice := float64(quantity) * itemPrice

	// Insert into transaction table
	transInsertStmt, err := db.Prepare("INSERT INTO transactions (trans_id, cart_id, trans_date, total_price, staff_id) VALUES (?, ?, now(), ?, ?)")
	if err != nil {
		log.Println("Failed to prepare SQL statement for inserting into transaction table:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer transInsertStmt.Close()

	// Generate a unique transaction ID
	transID := generateUniqueID(staffName + bookTitle + "ABC")

	_, err = transInsertStmt.Exec(transID, cartID, totalPrice, staffID)
	if err != nil {
		log.Println("Failed to execute SQL statement for inserting into transactions table:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Transaction added successfully!")
}

func generateUniqueID(input string) string {
	id := ""
	for _, c := range input {
		id += strconv.Itoa(int(c))
	}
	return id
}
