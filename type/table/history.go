package table

import "time"

type History struct {
	Id           *uint64    `gorm:"primaryKey"`
	OwnerId      *uint64    `gorm:"index;not null"`
	Owner        *User      `gorm:"foreignKey:OwnerId"`
	OwnerMatchId *uint64    `gorm:"index"`
	BookUserId   *uint64    `gorm:"index"`
	BookMatchId  *uint64    `gorm:"index"`
	MatchTime    *time.Time `gorm:"precision:6"`
}
