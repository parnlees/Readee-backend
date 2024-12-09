package table

import "time"

type Banners struct {
	BannerId  *uint64    `gorm:"column:banner_id;primaryKey;autoIncrement"`
	ImageUrl  *string    `gorm:"type:VARCHAR(512);not null"`
	Link      *string    `gorm:"type:VARCHAR(512);not null"`
	IsActive  *bool      `gorm:"default:true"`
	CreateAt  *time.Time `gorm:"autoCreateTime"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime"`
}
