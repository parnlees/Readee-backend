package table

import (
	"time"
)

type Book struct {
	//Id              *uint64    `gorm:"primaryKey;autoIncrement"`
	BookId          *uint64    `gorm:"column:book_id;primaryKey;autoIncrement"`
	OwnerId         *uint64    `gorm:"not null"`
	Owner           *User      `gorm:"foreignKey:OwnerId;references:UserId"`
	BookName        *string    `gorm:"type:VARCHAR(256);not null"`
	Author          *string    `gorm:"type:VARCHAR(256);not null"`
	BookPicture     *string    `gorm:"type:text;not null"`
	BookDescription *string    `gorm:"type:VARCHAR(256);not null"`
	GenreId         *uint64    `gorm:"not null"`
	Quality         *uint64    `gorm:"not null"`
	IsTraded        *bool      `gorm:"default:false"`
	CreatedAt       *time.Time `gorm:"autoCreateTime"`
	UpdatedAt       *time.Time `gorm:"autoUpdateTime"`
	IsReported      *bool      `gorm:"default:false"`
	ExpiredDate     *time.Time `gorm:"autoCreateTime"`
}
