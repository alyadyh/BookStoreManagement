package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/book-crud/generates"
	_ "github.com/go-sql-driver/mysql"
)

func AddBook(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Failed to parse form data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	bookTitle := r.Form.Get("bookTitle")
	authorName := r.Form.Get("authorName")
	publisherName := r.Form.Get("publisherName")
	publicationDate := r.Form.Get("publicationDate")
	ISBN := r.Form.Get("ISBN")
	priceStr := r.Form.Get("price")
	stockQtyStr := r.Form.Get("stockQty")
	genreName := r.Form.Get("genreName")

	bookID := generateUniqueID(bookTitle)
	authorID := generateUniqueID(authorName)
	publisherID := generateUniqueID(publisherName)
	genreID := generateUniqueID(genreName)

	log.Println("bookTitle:", bookTitle)
	log.Println("authorName:", authorName)
	log.Println("publisherName:", publisherName)
	log.Println("publicationDate:", publicationDate)
	log.Println("ISBN:", ISBN)
	log.Println("price:", priceStr)
	log.Println("stockQty:", stockQtyStr)
	log.Println("genreName:", genreName)

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

	// Insert into genres table
	genreStmt, err := db.Prepare("INSERT INTO genres (genre_id, genre_name) VALUES (?, ?)")
	if err != nil {
		log.Println("Failed to prepare SQL statement for inserting into genres table:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer genreStmt.Close()

	_, err = genreStmt.Exec(genreID, genreName)
	if err != nil {
		log.Println("Failed to execute SQL statement for inserting into genres table:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Convert price and stockQty to float64 and int types respectively
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		log.Println("Failed to parse price value:", err)
		http.Error(w, "Invalid price value", http.StatusBadRequest)
		return
	}

	stockQty, err := strconv.Atoi(stockQtyStr)
	if err != nil {
		log.Println("Failed to parse stock quantity value:", err)
		http.Error(w, "Invalid stock quantity value", http.StatusBadRequest)
		return
	}

	// Insert into books table
	bookStmt, err := db.Prepare("INSERT INTO books (book_id, title, author_id, publisher_id, publication_date, ISBN, price, stock_qty, genre_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Failed to prepare SQL statement for inserting into books table:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer bookStmt.Close()

	_, err = bookStmt.Exec(bookID, bookTitle, authorID, publisherID, publicationDate, ISBN, price, stockQty, genreID)
	if err != nil {
		log.Println("Failed to execute SQL statement for inserting into books table:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Book added successfully!")
}

// Helper function to generate a unique ID based on the ASCII values of the provided string
func generateUniqueID(str string) string {
	// Convert the book title to lowercase and remove spaces
	title := strings.ToLower(strings.ReplaceAll(str, " ", ""))

	// Calculate the sum of ASCII values of each character in the title
	sum := 0
	for _, ch := range title {
		sum += int(ch)
	}

	// Generate the unique ID using the sum of ASCII values
	return fmt.Sprintf("ID_%d", sum)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

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
