package table

import (
	"time"
)

type User struct {
	UserId       *uint64    `gorm:"primaryKey;autoIncrement"`
	Token        *string    `gorm:"type:varchar(256) ; unique_index"`
	Email        *string    `gorm:"type:VARCHAR(256); unique_index; not null"`
	Username     *string    `gorm:"type:VARCHAR(256); unique_index; not null"`
	Password     *string    `gorm:"type:VARCHAR(256); not null"`
	PhoneNumber  *string    `gorm:"type:VARCHAR(256)"`
	ProfileUrl   *string    `gorm:"type:VARCHAR(256)"`
	Firstname    *string    `gorm:"type:VARCHAR(256)"`
	Lastname     *string    `gorm:"type:VARCHAR(256)"`
	Gender       *string    `gorm:"type:VARCHAR(256)"`
	VerifyStatus *bool      `gorm:"default:true"` // Many-to-Many relationship
	CreatedAt    *time.Time `gorm:"autoCreateTime"`
	UpdatedAt    *time.Time `gorm:"autoCreateTime"`
	SecKey       *string    `gorm:"type:VARCHAR(256)"`

	Genres []*Genre `gorm:"many2many:user_genres;"`
}
