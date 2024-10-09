package endpoint

import (
	//"strconv"
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

	Genres []*Genre `gorm:"many2many:user_genres;"`
}

type Book struct {
	BookId          *uint64    `gorm:"primaryKey;autoIncrement"`
	OwnerId         *uint64    `gorm:"not null"`
	Owner           *User      `gorm:"foreignKey:OwnerId;references:UserId"`
	BookName        *string    `gorm:"type:VARCHAR(256);not null"`
	BookPicture     *string    `gorm:"type:VARCHAR(256);not null"`
	BookDescription *string    `gorm:"type:VARCHAR(256);not null"`
	GenreId         *uint64    `gorm:"not null"`
	Quality         *uint64    `gorm:"not null"`
	IsTraded        *bool      `gorm:"default:true"`
	CreatedAt       *time.Time `gorm:"precision:6"`
	UpdatedAt       *time.Time `gorm:"precision:6"`
}

type Genre struct {
	GenreId *uint64 `gorm:"primaryKey;autoIncrement"`
	Name    *string `gorm:"type:VARCHAR(256); unique_index; not null"`
	Users   []*User `gorm:"many2many:user_genres;"` // Many-to-Many relationship
}

type History struct {
	Id           *uint64    `gorm:"primaryKey;autoIncrement"`
	OwnerId      *uint64    `gorm:"index;not null"`
	Owner        *User      `gorm:"foreignKey:OwnerId;references:UserId"`
	OwnerMatchId *uint64    `gorm:"index"`
	BookMatchId  *uint64    `gorm:"not null;unique"`
	Book         *Book      `gorm:"foreignKey:BookMatchId;references:BookId"`
	MatchTime    *time.Time `gorm:"precision:6"`
}

type Log struct {
	//id
	LogId      *uint64 `gorm:"primaryKey;autoIncrement"`
	BookLikeId *uint64 `gorm:"not null"` // ใช้ BookId เป็น foreign key
	Book       *Book   `gorm:"foreignKey:BookLikeId; references:BookId"`
	LikerId *uint64 `gorm:"not null"`                                // User ที่ถูกใจหนังสือ
	User       *User   `gorm:"foreignKey:LikerId;references:UserId"` // สอดคล้องกับ UserId
	Liked      *bool   `gorm:"default:false"`
}

type Match struct {
	MatchId       *uint64    `gorm:"primaryKey;autoIncrement"`
	OwnerId       *uint64    `gorm:"not null"` // เจ้าของแมทช์
	Owner         *User      `gorm:"foreignKey:OwnerId;references:UserId"`
	MatchedUserId *uint64    `gorm:"not null"` // ผู้ใช้ที่ถูกแมทช์กับเจ้าของ
	MatchedUser   *User      `gorm:"foreignKey:MatchedUserId;references:UserId"`
	OwnerBookId   *uint64    `gorm:"not null; uniqueIndex"` // หนังสือของเจ้าของ
	OwnerBook     *Book      `gorm:"foreignKey:OwnerBookId;references:BookId"`
	MatchedBookId *uint64    `gorm:"not null; uniqueIndex"` // หนังสือของผู้ใช้ที่ถูกแมทช์
	MatchedBook   *Book      `gorm:"foreignKey:MatchedBookId;references:BookId"`
	MatchTime     *time.Time `gorm:"precision:6"` // เวลาที่เกิดการแมทช์
	TradeTime     *time.Time `gorm:"precision:6"` // เวลาที่เกิดการแลกเปลี่ยน
}

type Message struct {
	MessageId *uint64    `gorm:"primaryKey;autoIncrement"`
	RoomId    *uint64    `gorm:"not null"` // ใช้ RoomId เป็น foreign key
	Room      *Room      `gorm:"foreignKey:RoomId;references:RoomId"`
	SenderId  *uint64    `gorm:"index;not null"`
	Message   *string    `gorm:"type:VARCHAR(256);not null"`
	CreateAt  *time.Time `gorm:"precision:6"`
}

type Rating struct {
	RatingId   *uint64    `gorm:"primaryKey;autoIncrement"`
	GiverId    *uint64    `gorm:"not null"` // FK ชี้ไปที่ User ที่ให้ rating
	Giver      *User      `gorm:"foreignKey:GiverId;references:UserId"`
	ReceiverId *uint64    `gorm:"not null"` // FK ชี้ไปที่ User ที่ได้รับ rating
	Receiver   *User      `gorm:"foreignKey:ReceiverId;references:UserId"`
	Rating     *float64   `gorm:"not null"`
	Score      *uint64    `gorm:"not null"`
	NumRate    *uint64    `gorm:"not null"`
	CreatedAt  *time.Time `gorm:"precision:6"`
	UpdatedAt  *time.Time `gorm:"precision:6"`
}

type Review struct {
	ReviewId   *uint64    `gorm:"primaryKey;autoIncrement"`
	GiverId    *uint64    `gorm:"not null"` // FK ชี้ไปที่ User ที่ให้ rating
	Giver      *User      `gorm:"foreignKey:GiverId;references:UserId"`
	ReceiverId *uint64    `gorm:"not null"` // FK ชี้ไปที่ User ที่ได้รับ rating
	Receiver   *User      `gorm:"foreignKey:ReceiverId;references:UserId"`
	TextReview string     `gorm:"type:VARCHAR(256); not null"`
	CreatedAt  *time.Time `gorm:"precision:6"`
	UpdatedAt  *time.Time `gorm:"precision:6"`
}

type Room struct {
	RoomId     *uint64 `gorm:"primaryKey;autoIncrement"`
	SenderId   *uint64 `gorm:"index"`
	Sender     *User   `gorm:"foreignKey:SenderId;references:UserId"`
	ReceiverId *uint64 `gorm:"index"`
	Receiver   *User   `gorm:"foreignKey:ReceiverId;references:UserId"`

	Messages []*Message `gorm:"foreignKey:RoomId"`
}
