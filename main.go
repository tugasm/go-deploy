package main

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Book adalah struktur untuk representasi data buku
type Book struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var (
	db *sql.DB
)

func init() {
	//if err := godotenv.Load(); err != nil {
	//	log.Fatalf("Error loading .env file: %v", err)
	//}
	connStr := "user=postgres password=Hacktiv8123 host=db.gxdymouplidfsyylzfei.supabase.co port=5432 dbname=postgres"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer db.Close()

	e := echo.New()

	port := os.Getenv("PORT")
	// API Endpoint untuk menambahkan buku
	e.POST("/books", createBook)

	// API Endpoint untuk mendapatkan daftar semua buku
	e.GET("/books", getBooks)

	// API Endpoint untuk mendapatkan buku berdasarkan ID
	e.GET("/books/:id", getBook)

	// API Endpoint untuk mengubah buku berdasarkan ID
	e.PUT("/books/:id", updateBook)

	// API Endpoint untuk menghapus buku berdasarkan ID
	e.DELETE("/books/:id", deleteBook)

	e.Logger.Fatal(e.Start(":" + port))
}

func createBook(c echo.Context) error {
	b := new(Book)
	if err := c.Bind(b); err != nil {
		return err
	}

	query := "INSERT INTO books (name, description) VALUES ($1, $2) RETURNING id"
	var id int
	err := db.QueryRow(query, b.Name, b.Description).Scan(&id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, id)
}

func getBooks(c echo.Context) error {
	query := "SELECT * FROM books"
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	books := []Book{}
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Name, &b.Description); err != nil {
			return err
		}
		books = append(books, b)
	}

	return c.JSON(http.StatusOK, books)
}

func getBook(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	query := "SELECT * FROM books WHERE id=$1"
	var b Book
	err := db.QueryRow(query, id).Scan(&b.ID, &b.Name, &b.Description)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, b)
}

func updateBook(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	b := new(Book)
	if err := c.Bind(b); err != nil {
		return err
	}

	query := "UPDATE books SET name=$2, description=$3 WHERE id=$1"
	_, err := db.Exec(query, id, b.Name, b.Description)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func deleteBook(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	query := "DELETE FROM books WHERE id=$1"
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
