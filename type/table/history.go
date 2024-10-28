package table

import (
	"time"
)

type History struct {
	HistoryId *uint64 `gorm:"primaryKey;autoIncrement"`
	OwnerId   *uint64 `gorm:"index;not null"`
	Owner     *User   `gorm:"foreignKey:OwnerId;references:UserId"`

	MatchedUserId *uint64 `gorm:"index;not null"`
	MatchedUser   *User   `gorm:"foreignKey:MatchedUserId;references:UserId"` // เพิ่ม field ผู้ใช้ที่แมทช์ไว้

	OwnerBookId *uint64 `gorm:"not null;unique"`
	OwnerBook   *Book   `gorm:"foreignKey:OwnerBookId;references:BookId"` // เปลี่ยนชื่อ field ให้ชัดเจน

	MatchedBookId *uint64 `gorm:"not null"`
	MatchedBook   *Book   `gorm:"foreignKey:MatchedBookId;references:BookId"` // เพิ่ม field หนังสือที่เทรดด้วย

	TradeTime *time.Time `gorm:"autoCreateTime"`
}
