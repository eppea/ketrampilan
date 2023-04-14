// db.go
package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func OpenDatabase() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "root:qwerty@tcp(localhost:3306)/absen_sekolah?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}

	return db, nil
}
