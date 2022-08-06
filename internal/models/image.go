package models

type Image struct {
	Id   int    `json:"-" db:"id"`
	Name string `json:"name" binding:"required"`
}
