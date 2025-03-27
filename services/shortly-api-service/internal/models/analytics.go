package models

import (
	"gorm.io/gorm"
	"time"
)

type Analytics struct {
	gorm.Model

	UrlID     string    `gorm:"index;not null"`
	Url       Url       `gorm:"foreignKey:UrlID"`
	ClickedAt time.Time `gorm:"autoCreateTime"`
	IPAddress string    `gorm:"not null"`
	UserAgent string    `gorm:"not null"`
	Referrer  string    `gorm:"size:255"`
	Country   string    `gorm:"size:100"`
	Device    string    `gorm:"size:50"`
	Browser   string    `gorm:"size:50"`
	OS        string    `gorm:"size:50"`
}
