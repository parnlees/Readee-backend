package table

import "time"

type History struct {
	Id           *uint64    `gorm:"primaryKey;autoIncrement"`
	OwnerId      *uint64    `gorm:"index;not null"`
	Owner        *User      `gorm:"foreignKey:OwnerId;references:UserId"`
	OwnerMatchId *uint64    `gorm:"index"`
	BookMatchId  *uint64    `gorm:"not null;unique"`
	Book         *Book      `gorm:"foreignKey:BookMatchId;references:BookId"`
	MatchTime    *time.Time `gorm:"autoCreateTime"`
}
