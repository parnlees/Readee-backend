package table

type Genre struct {
	Id   *uint64 `gorm:"primarykey"`
	Name *string `gorm:"type:VARCHAR(256); unique_index; not null"`
}
