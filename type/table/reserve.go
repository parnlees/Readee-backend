package table

import "time"

type Reserve struct {
	Id        *uint64    `gorm:"primary_key"`
	User      *User      `gorm:"foreignKey:UserId"`
	Book      *Book      `gorm:"foreignKey:BookId"`
	CreatedAt *time.Time `gorm:"precision:6"`
	UpdatedAt *time.Time `gorm:"precision:6"`
}
