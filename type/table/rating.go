package table

import "time"

type Rating struct {
	RatingId   *uint64    `gorm:"primaryKey;autoIncrement"`
	ReviewId   *uint64    `gorm:"not null"` // FK ชี้ไปที่ Review ที่เกี่ยวข้อง
	GiverId    *uint64    `gorm:"not null"` // FK ชี้ไปที่ User ที่ให้ rating
	Giver      *User      `gorm:"foreignKey:GiverId;references:UserId"`
	ReceiverId *uint64    `gorm:"not null"` // FK ชี้ไปที่ User ที่ได้รับ rating
	Receiver   *User      `gorm:"foreignKey:ReceiverId;references:UserId"`
	Rating     *float64   `gorm:"not null"`
	Score      *uint64    `gorm:"not null"`
	NumRate    *uint64    `gorm:"not null"`
	CreatedAt  *time.Time `gorm:"autoCreateTime"`
	UpdatedAt  *time.Time `gorm:"autoCreateTime"`
}
