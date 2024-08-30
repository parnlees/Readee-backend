package table

type Room struct {
	RoomId     *uint64 `gorm:"primaryKey;autoIncrement"`
	SenderId   *uint64 `gorm:"index"`
	Sender     *User   `gorm:"foreignKey:SenderId;references:UserId"`
	ReceiverId *uint64 `gorm:"index"`
	Receiver   *User   `gorm:"foreignKey:ReceiverId;references:UserId"`

	Messages []*Message `gorm:"foreignKey:RoomId"`
}
