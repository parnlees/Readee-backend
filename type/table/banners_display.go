package table

import "time"

type Banners_display struct {
	id        *uint64    `gorm:"column:banner_id;primaryKey;autoIncrement"`
	userId    *uint64    `gorm:"not null"`
	user      *User      `gorm:"foreignKey:userId;references:userId"`
	bannerId  *uint64    `gorm:"not null"`
	banner    *Banners   `gorm:"foreignKey:banner_id;references:banner_id"`
	DisplayAt *time.Time `gorm:"autoCreateTime"`
}
