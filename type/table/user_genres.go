package table

type UserGenres struct {
	Genre_genre_id *uint64 `gorm:"index"`
	User_user_id   *uint64 `gorm:"index"`
}
