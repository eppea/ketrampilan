// model.go

package main

type Murid struct {
	ID    int    `json:"id" gorm:"primary_key"`
	Nama  string `json:"nama" gorm:"not null"`
	Hadir bool   `json:"hadir" gorm:"not null;default:false"`
}
