package table

import (
	"time"
)

type Report struct {
	ReportId *uint64    `gorm:"primaryKey;autoIncrement"`
	UserId   *uint64    `gorm:"not null"` // Foreign key for User
	User     *User      `gorm:"foreignKey:UserId;references:UserId" json:"-"`
	BookId   *uint64    `gorm:"not null"` // Foreign key for Book
	Book     *Book      `gorm:"foreignKey:BookId;references:BookId" json:"-"`
	ReportAt *time.Time `gorm:"autoCreateTime"`
}
