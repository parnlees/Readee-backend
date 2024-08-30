package table

import "time"

type Review struct {
	ReviewId   *uint64    `gorm:"primaryKey;autoIncrement"`
	GiverId    *uint64    `gorm:"not null"` // FK ชี้ไปที่ User ที่ให้ rating
	Giver      *User      `gorm:"foreignKey:GiverId;references:UserId"`
	ReceiverId *uint64    `gorm:"not null"` // FK ชี้ไปที่ User ที่ได้รับ rating
	Receiver   *User      `gorm:"foreignKey:ReceiverId;references:UserId"`
	TextReview string     `gorm:"type:VARCHAR(256); not null"`
	CreatedAt  *time.Time `gorm:"precision:6"`
	UpdatedAt  *time.Time `gorm:"precision:6"`
}
