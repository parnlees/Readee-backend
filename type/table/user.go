package table

import "time"

type User struct {
	Id           *uint64 `gorm:"primarykey"`
	Token        *string `gorm:"type:varchar(256) ; unique_index; not null"`
	Email        *string `gorm:"type:VARCHAR(256); unique_index; not null"`
	Username     *string `gorm:"type:VARCHAR(256); unique_index; not null"`
	Password     *string `gorm:"type:VARCHAR(256); not null"`
	PhoneNumber  *string `gorm:"type:VARCHAR(256); not null"`
	ProfileUrl   *string `gorm:"type:VARCHAR(256); not null"`
	Firstname    *string `gorm:"type:VARCHAR(256); not null"`
	Lastname     *string `gorm:"type:VARCHAR(256); not null"`
	Gender       *string `gorm:"type:VARCHAR(256); not null"`
	VerifyStatus *bool   `gorm:"default:true"`
	//FavoriteGenre *Genre  `gorm:"foreignKey:FavoriteGenreID"`
	//Books         []Book     `gorm:"foreignKey:OwnerId"`
	//Histories     []History  `gorm:"foreignKey:OwnerId"`
	CreatedAt *time.Time `gorm:"precision:6"`
	UpdatedAt *time.Time `gorm:"precision:6"`
}
