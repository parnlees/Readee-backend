package table

import "time"

type Match struct {
	MatchId       *uint64    `gorm:"primaryKey;autoIncrement"`
	OwnerId       *uint64    `gorm:"not null"` // เจ้าของแมทช์
	Owner         *User      `gorm:"foreignKey:OwnerId;references:UserId"`
	MatchedUserId *uint64    `gorm:"not null"` // ผู้ใช้ที่ถูกแมทช์กับเจ้าของ
	MatchedUser   *User      `gorm:"foreignKey:MatchedUserId;references:UserId"`
	OwnerBookId   *uint64    `gorm:"not null; uniqueIndex"` // หนังสือของเจ้าของ
	OwnerBook     *Book      `gorm:"foreignKey:OwnerBookId;references:BookId"`
	MatchedBookId *uint64    `gorm:"not null; uniqueIndex"` // หนังสือของผู้ใช้ที่ถูกแมทช์
	MatchedBook   *Book      `gorm:"foreignKey:MatchedBookId;references:BookId"`
	MatchTime     *time.Time `gorm:"precision:6"` // เวลาที่เกิดการแมทช์
	TradeTime     *time.Time `gorm:"precision:6"` // เวลาที่เกิดการแลกเปลี่ยน
}
