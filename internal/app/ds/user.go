package ds

type User struct {
	ID       uint64 `gorm:"primary_key"`
	Username string `gorm:"type:varchar(50); unique; not null"`
	Password string `gorm:"type:varchar(50); not null"`
	IsMod    bool   `gorm:"boolean; not null;"`
}
