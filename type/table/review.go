package table

import "time"

type Review struct {
	Id         *uint64    `gorm:"primary_key"`
	OwnerId    *uint64    `gorm:"not null"`
	BuddyId    *uint64    `gorm:"not null"`
	TextReview string     `gorm:"type:VARCHAR(256); not null"`
	CreatedAt  *time.Time `gorm:"precision:6"`
	UpdatedAt  *time.Time `gorm:"precision:6"`
}
