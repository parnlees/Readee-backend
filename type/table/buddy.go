package table

type Buddy struct {
	Id          *uint64 `gorm:"primaryKey"`
	OwnerId     *uint64 `gorm:"not null"`
	Owner       *User   `gorm:"foreignKey:OwnerId"`
	BookOwnerId *uint64 `gorm:"not null"`
	BookBuddyId *uint64 `gorm:"not null"`
}
