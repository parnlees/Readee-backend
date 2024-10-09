package table

import "time"

type Match struct {
	MatchId       *uint64    `gorm:"primaryKey;autoIncrement"`
	OwnerId       *uint64    `gorm:"not null"` // เจ้าของแมทช์
	Owner         *User      `gorm:"foreignKey:OwnerId;references:UserId"`
	MatchedUserId *uint64    `gorm:"not null"` // ผู้ใช้ที่ถูกแมทช์กับเจ้าของ
	MatchedUser   *User      `gorm:"foreignKey:MatchedUserId;references:UserId"`
	OwnerBookId   *uint64    `gorm:"not null"` // หนังสือของเจ้าของ
	OwnerBook     *Book      `gorm:"foreignKey:OwnerBookId;references:BookId"`
	MatchedBookId *uint64    `gorm:"not null"` // หนังสือของผู้ใช้ที่ถูกแมทช์
	MatchedBook   *Book      `gorm:"foreignKey:MatchedBookId;references:BookId"`
	MatchTime     *time.Time `gorm:"autoCreateTime"` // เวลาที่เกิดการแมทช์
	TradeTime     *time.Time `gorm:"null"`           // เวลาที่เกิดการแลกเปลี่ยน

	TradeRequestStatus string  `gorm:"type:varchar(10);default:'none'"` // none, pending, accepted, rejected
	RequestInitiatorId *uint64 `gorm:"null"`                            // id ของคนเริ่มส่ง request

}
