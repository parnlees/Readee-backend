package table

import "time"

type Report struct {
	ReportId *uint64    `gorm:"primaryKey;autoIncrement"`
	UserId   *uint64    `gorm:"not null"`
	User     *User      `gorm:"foreignKey:UserId;references:UserId"`
	BookId   *uint64    `gorm:"not null"`
	Book     *Book      `gorm:"foreignKey:BookId;references:BookId"`
	ReportAt *time.Time `gorm:"autoCreateTime"`
}
