package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

type blog struct {
	Id    string `json:"id" form:"id" query:"id"`
	Title string `json:"title" form:"title" query:"title"`
	Body  string `json:"body" form:"body" query:"body"`
}

type json_error struct {
	Err     error
	Message string
}

var all_blogs []blog
var database *sql.DB

var error_response = json_error{Err: nil, Message: "There was an error"}

func main() {
	fmt.Println("Hello, World!")
	var err error
	database, err = sql.Open("sqlite3", "./blog_db.db")
	if err != nil {
		log.Fatal(err)
	}
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS blog_entry (id TEXT, title TEXT, body TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()

	e := echo.New()
	e.GET("/", get_request)
	e.POST("/post", post_request)
	e.DELETE("/post/:id", delete_request)
	e.PUT("/post/:id", put_request)

	e.Start(":8080")
}

func post_request(c echo.Context) error {
	b := new(blog)
	err := c.Bind(b)
	if err != nil {
		log.Fatal(err)
		error_response.Err = err
		return c.JSON(http.StatusBadRequest, error_response)
	}
	add_entry(b.Id, b.Title, b.Body)
	return c.JSON(http.StatusOK, b)
}

func get_request(c echo.Context) error {
	all_blogs = nil
	rows, err := database.Query("SELECT id, title, body FROM blog_entry")
	error_response.Err = err
	if err != nil {
		log.Fatal(err)
		return c.JSON(http.StatusNotFound, error_response)
	}
	var blog_id, blog_title, blog_body string

	defer rows.Close()
	for rows.Next() {
		rows.Scan(&blog_id, &blog_title, &blog_body)
		temp_blog := blog{Id: blog_id, Title: blog_title, Body: blog_body}
		all_blogs = append(all_blogs, temp_blog)
	}
	return c.JSON(http.StatusOK, all_blogs)
}

func delete_request(c echo.Context) error {
	id := c.Param("id")
	if check_id(id) == false {
		return c.JSON(http.StatusBadRequest, error_response)
	}
	b := new(blog)
	err := database.QueryRow("SELECT id, title, body FROM blog_entry WHERE id=?", id).Scan(&b.Id, &b.Title, &b.Body)
	if err != nil {
		log.Fatal(err)
		error_response.Err = err
		return c.JSON(http.StatusInternalServerError, error_response)
	}
	delete_entry(id)
	return c.JSON(http.StatusOK, b)
}

func put_request(c echo.Context) error {
	id := c.Param("id")
	if check_id(id) {
		return c.JSON(http.StatusBadRequest, error_response)
	}
	b := new(blog)
	err := c.Bind(b)
	if err != nil {
		log.Fatal(err)
		error_response.Err = err
		return c.JSON(http.StatusInternalServerError, error_response)
	}
	statement, err := database.Prepare("UPDATE blog_entry SET title=?, body=? WHERE id=?")
	if err != nil {
		log.Fatal(err)
		error_response.Err = err
		return c.JSON(http.StatusNotFound, error_response)
	}
	statement.Exec(b.Title, b.Body, id)
	return c.JSON(http.StatusOK, b)
}

func add_entry(id string, title string, body string) {
	statement, err := database.Prepare("INSERT INTO blog_entry (id, title, body) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec(id, title, body)
	log.Println("added")
}

func delete_entry(blog_id string) {
	statement, err := database.Prepare("DELETE FROM blog_entry WHERE id=?")
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec(blog_id)
	log.Println("deleted")
}

func check_id(id string) bool {
	if id == " " {
		error_response.Err = nil
		return false
	}
	return true
}
