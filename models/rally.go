package models

import "gorm.io/gorm"

type Rally struct {
	gorm.Model
	RallyName string  `gorm:"not null;unique"`
	Password  string  `gorm:"not null"`
	Stamp     []Stamp `gorm:"constraint:Ondelete:CASCADE"`
}
