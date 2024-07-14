package table

type Room struct {
	Id         *uint64 `gorm:"primary_key"`
	SenderId   *uint64 `gorm:"index"`
	ReceiverId *uint64 `gorm:"index"`
}
