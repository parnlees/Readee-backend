package table

import (
	"time"
)

type Book struct {
	BookId          *uint64    `gorm:"primaryKey;autoIncrement"`
	OwnerId         *uint64    `gorm:"not null"`
	Owner           *User      `gorm:"foreignKey:OwnerId;references:UserId"`
	BookName        *string    `gorm:"type:VARCHAR(256);not null"`
	BookPicture     *string    `gorm:"type:VARCHAR(256);not null"`
	BookDescription *string    `gorm:"type:VARCHAR(256);not null"`
	GenreId         *uint64    `gorm:"not null"`
	Quality         *uint64    `gorm:"not null"`
	IsTraded        *bool      `gorm:"default:true"`
	CreatedAt       *time.Time `gorm:"precision:6"`
	UpdatedAt       *time.Time `gorm:"precision:6"`

	Logs []*Log `gorm:"foreignKey:BookLikeId"`
}
