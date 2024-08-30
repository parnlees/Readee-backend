package table

import "time"

type Rating struct {
	RatingId   *uint64    `gorm:"primaryKey;autoIncrement"`
	GiverId    *uint64    `gorm:"not null"` // FK ชี้ไปที่ User ที่ให้ rating
	Giver      *User      `gorm:"foreignKey:GiverId;references:UserId"`
	ReceiverId *uint64    `gorm:"not null"` // FK ชี้ไปที่ User ที่ได้รับ rating
	Receiver   *User      `gorm:"foreignKey:ReceiverId;references:UserId"`
	Rating     *float64   `gorm:"not null"`
	Score      *uint64    `gorm:"not null"`
	NumRate    *uint64    `gorm:"not null"`
	CreatedAt  *time.Time `gorm:"precision:6"`
	UpdatedAt  *time.Time `gorm:"precision:6"`
}
