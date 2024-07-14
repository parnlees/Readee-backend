package table

import "time"

type Message struct {
	Id       *string    `gorm:"type:CHAR(7);primaryKey"`
	RoomId   *string    `gorm:"type:VARCHAR(256);not null"`
	SenderId *uint64    `gorm:"index;not null"`
	Message  *string    `gorm:"type:VARCHAR(256);not null"`
	CreateAt *time.Time `gorm:"precision:6"`
}
