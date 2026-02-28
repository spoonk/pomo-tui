package storage

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	ID uint `gorm:"primarykey;autoIncrement"`

	Name        string
	Description string
}
