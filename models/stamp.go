package models

import "gorm.io/gorm"

type Stamp struct {
	gorm.Model
	RallyID            uint   `gorm:"not null"`
	StampNumber        uint   `gorm:"not null"`
	ImageURL           string `gorm:"not null"`
	InstallationPotion string `gorm:"not null"`
}
