package table

type Log struct {
	LogId      *uint64 `gorm:"primaryKey;autoIncrement"`
	BookLikeId *uint64 `gorm:"not null"` // ใช้ BookId เป็น foreign key
	Book       *Book   `gorm:"foreignKey:BookLikeId; reference:BookId"`
	UserId     *uint64 `gorm:"not null"`                            // User ที่ถูกใจหนังสือ
	User       *User   `gorm:"foreignKey:UserId;references:UserId"` // สอดคล้องกับ UserId
	Liked      *bool   `gorm:"default:false"`
}
