package table

type Rating struct {
	Id      *uint64  `gorm:"primary_key"`
	OwnerId *uint64  `gorm:"not null"`
	Rating  *float64 `gorm:"not null"`
	Score   *uint64  `gorm:"not null"`
	NumRate *uint64  `gorm:"not null"`
}
