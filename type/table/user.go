package table

import (
	"time"
)

type User struct {
	UserId       *uint64    `gorm:"primaryKey;autoIncrement"`
	Token        *string    `gorm:"type:varchar(256) ; unique_index; not null"`
	Email        *string    `gorm:"type:VARCHAR(256); unique_index; not null"`
	Username     *string    `gorm:"type:VARCHAR(256); unique_index; not null"`
	Password     *string    `gorm:"type:VARCHAR(256); not null"`
	PhoneNumber  *string    `gorm:"type:VARCHAR(256); not null"`
	ProfileUrl   *string    `gorm:"type:VARCHAR(256); not null"`
	Firstname    *string    `gorm:"type:VARCHAR(256); not null"`
	Lastname     *string    `gorm:"type:VARCHAR(256); not null"`
	Gender       *string    `gorm:"type:VARCHAR(256); not null"`
	VerifyStatus *bool      `gorm:"default:true"` // Many-to-Many relationship
	CreatedAt    *time.Time `gorm:"autoCreateTime"`
	UpdatedAt    *time.Time `gorm:"autoCreateTime"`

	Genres []*Genre `gorm:"many2many:user_genres;"`
}
