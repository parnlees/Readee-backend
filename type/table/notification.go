package table

import "time"

type Notification struct {
	NotiId         uint64    `gorm:"primaryKey;autoIncrement"`
	NotiMessage    string    `gorm:"type:VARCHAR(256);not null"`
	NotiType       string    `gorm:"type:VARCHAR(50);not null"`
	NotiSenderId   uint64    `gorm:"not null;index"`             // FK to User table
	NotiReceiverId uint64    `gorm:"not null;index"`             // FK to User table
	BookId         uint64    `gorm:"not null;index"`             // FK to Book table
	SendAt         time.Time `gorm:"autoCreateTime"`
}


