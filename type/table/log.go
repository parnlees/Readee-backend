package table

import "time"

type Log struct {
	//id
	LogId      *uint64    `gorm:"primaryKey;autoIncrement"`
	BookLikeId *uint64    `gorm:"not null"` // ใช้ BookId เป็น foreign key
	Book       *Book      `gorm:"foreignKey:BookLikeId; references:BookId"`
	UserLikeId *uint64    `gorm:"not null"`                                // User ที่ถูกใจหนังสือ
	User       *User      `gorm:"foreignKey:UserLikeId;references:UserId"` // สอดคล้องกับ UserId
	Liked      *bool      `gorm:"default:false"`                           // swipe to left = false, swipe to right = true
	CreatedAt  *time.Time `gorm:"autoCreateTime"`
}
