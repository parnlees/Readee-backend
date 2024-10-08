package table

type BookPicture struct {
	PictureId  *uint64 `gorm:"primaryKey;autoIncrement"`
	BookId     *uint64 `gorm:"not null"`
	FilePath   *string `gorm:"type:VARCHAR(256);not null"` // ที่เก็บรูปภาพ
}