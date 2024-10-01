package table

import "time"

type Message struct {
	MessageId *uint64    `gorm:"primaryKey;autoIncrement"`
	RoomId    *uint64    `gorm:"not null"` // ใช้ RoomId เป็น foreign key
	Room      *Room      `gorm:"foreignKey:RoomId;references:RoomId"`
	SenderId  *uint64    `gorm:"index;not null"`
	Message   *string    `gorm:"type:VARCHAR(256);not null"`
	CreateAt  *time.Time `gorm:"autoCreateTime"`
}
