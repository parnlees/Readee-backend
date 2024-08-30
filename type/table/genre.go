package table

type Genre struct {
	GenreId *uint64 `gorm:"primaryKey;autoIncrement"`
	Name    *string `gorm:"type:VARCHAR(256); unique_index; not null"`
	Users   []*User `gorm:"many2many:user_genres;"` // Many-to-Many relationship
}
