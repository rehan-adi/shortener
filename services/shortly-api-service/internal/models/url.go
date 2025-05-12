package models

import (
	"gorm.io/gorm"
)

type Url struct {
	gorm.Model

	OriginalURL string      `gorm:"not null"`
	ShortKey    string      `gorm:"size:50;uniqueIndex;not null"`
	Title       string      `gorm:"size:255"`
	UserID      *string     `gorm:"index"`
	User        *User       `gorm:"foreignKey:UserID"`
	Clicks      int         `gorm:"default:0"`
	Analytics   []Analytics `gorm:"foreignKey:UrlID"`
}
