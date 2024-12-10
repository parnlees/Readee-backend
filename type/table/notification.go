package table

import "time"

type Notification struct {
	NotiId         *uint64    `gorm:"primaryKey;autoIncrement"`
	NotiMessage    *string    `gorm:"type:VARCHAR(256); not null"`
	NotiSenderId   *uint64    `gorm:"not null"` // FK ชี้ไปที่ User ที่ให้ rating
	NotiSender     *User      `gorm:"foreignKey:NotiSenderId;references:UserId"`
	NotiReceiverId *uint64    `gorm:"not null"` // FK ชี้ไปที่ User ที่ให้ rating
	NotiReceiver   *User      `gorm:"foreignKey:NotiReceiverId;references:UserId"`
	SendAt         *time.Time `gorm:"autoCreateTime"`
}
