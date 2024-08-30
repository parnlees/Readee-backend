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
	CreatedAt    *time.Time `gorm:"precision:6"`
	UpdatedAt    *time.Time `gorm:"precision:6"`

	Books  []*Book  `gorm:"foreignKey:OwnerId"`
	Genres []*Genre `gorm:"many2many:user_genres;"`
	//Ratings      []*Rating  `gorm:"foreignKey:GiverId"`
	//Reviews      []*Review  `gorm:"foreignKey:GiverId"`
	//Histories    []*History `gorm:"foreignKey:OwnerId"`
	Logs []*Log `gorm:"foreignKey:UserId"`
	//Matches      []*Match   `gorm:"foreignKey:OwnerId"`
	//SentRooms    []*Room    `gorm:"foreignKey:SenderId"`
	//ReceiveRooms []*Room    `gorm:"foreignKey:ReceiverId"`
}
