package table

import "time"

type Match struct {
	Id           *uint64    `gorm:"primary_key"`
	OwnerId      *uint64    `gorm:"not null"`
	Owner        *User      `gorm:"foreignKey:OwnerId"`
	OwnerMatchId *uint64    `gorm:"not null"`
	MatchTime    *time.Time `gorm:"precision:6"`
	TradeTime    *time.Time `gorm:"precision:6"`
	BookUserId   uint64     `gorm:"not null"`
	BookMatchId  uint64     `gorm:"not null"`
}
