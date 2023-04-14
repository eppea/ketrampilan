package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Student struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	// Setting up database connection
	db, err := sql.Open("mysql", "root:qwerty@tcp(127.0.0.1:3306)/absen_sekolah")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/students", func(c echo.Context) error {
		// Query database for all students
		rows, err := db.Query("SELECT id, name FROM students")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Create slice of students
		students := []Student{}

		// Loop through rows and append to slice
		for rows.Next() {
			var s Student
			err := rows.Scan(&s.ID, &s.Name)
			if err != nil {
				log.Fatal(err)
			}
			students = append(students, s)
		}

		// Return slice of students as JSON
		return c.JSON(http.StatusOK, students)
	})

	e.GET("/students/:id", func(c echo.Context) error {
		id := c.Param("id")

		// Query database for student with given ID
		row := db.QueryRow("SELECT id, name FROM students WHERE id=?", id)

		// Create new student struct
		var s Student

		// Scan row into student struct
		err := row.Scan(&s.ID, &s.Name)
		if err != nil {
			log.Fatal(err)
		}

		// Return student struct as JSON
		return c.JSON(http.StatusOK, s)
	})

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
