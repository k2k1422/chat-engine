package Model

import "time"

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
}

type Chat struct {
	ID            uint      `gorm:"primaryKey"`
	FromUser      string    `gorm:"not null"`
	Username      string    `gorm:"not null"`
	Message       string    `gorm:"size:1024"`
	Time          time.Time `gorm:"not null"`
	Read          bool
	ReadTime      time.Time
	Delivered     bool
	DeliveredTime time.Time
}
