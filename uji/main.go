package main

import (
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type Murid struct {
	ID    int    `json:"id" gorm:"primary_key"`
	Nama  string `json:"nama"`
	Hadir bool   `json:"hadir"`
}

func OpenDatabase() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "root:qwerty@tcp(localhost:3306)/absen_sekolah?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	// Membuat koneksi database
	db, err := OpenDatabase()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Membuat instance Echo
	e := echo.New()

	// Middleware untuk menghubungkan database ke konteks Echo
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	})

	// Route untuk daftar murid
	e.GET("/murid", func(c echo.Context) error {
		db := c.Get("db").(*gorm.DB)

		var murids []Murid
		db.Find(&murids)

		return c.JSON(http.StatusOK, murids)
	})

	// Route untuk menandai murid hadir
	e.PUT("/murid/:id/hadir", func(c echo.Context) error {
		db := c.Get("db").(*gorm.DB)

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}

		var murid Murid
		db.First(&murid, id)
		if murid.ID == 0 {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Murid tidak ditemukan"})
		}

		murid.Hadir = true
		db.Save(&murid)

		return c.JSON(http.StatusOK, murid)
	})

	// Menjalankan server di port 8080
	e.Logger.Fatal(e.Start(":8080"))
}
