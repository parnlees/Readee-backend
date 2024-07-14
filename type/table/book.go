package table

import "time"

type Book struct {
	Id              *uint64    `gorm:"primary_key"`
	OwnerId         *uint64    `gorm:"not null"`
	Owner           *User      `gorm:"foreignKey:OwnerId"`
	BookName        *string    `gorm:"type:VARCHAR(256);not null"`
	BookPicture     *string    `gorm:"type:VARCHAR(256);not null"`
	BookDescription *string    `gorm:"type:VARCHAR(256);not null"`
	GenreId         *uint64    `gorm:"not null"`
	Genre           *Genre     `gorm:"foreignKey:GenreId"`
	Quality         *uint64    `gorm:"not null"`
	IsTraded        *bool      `gorm:"default:true"`
	CreatedAt       *time.Time `gorm:"precision:6"`
	UpdatedAt       *time.Time `gorm:"precision:6"`
}
