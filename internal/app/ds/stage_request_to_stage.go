package ds

type StageRequestToStage struct {
	RequestID    uint64       `gorm:"primary_key"`
	StageID      uint64       `gorm:"primary_key"`
	StageRequest StageRequest `gorm:"foreignKey:RequestID; references:ID"`
	Stage        Stage        `gorm:"foreignKey:StageID; references:ID"`
	ProducName   string       `gorm:"type:varchar(100); not null"`
}
