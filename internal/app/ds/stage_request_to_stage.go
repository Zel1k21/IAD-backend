package ds

type StageRequestToStage struct {
	RequestID    uint64       `gorm:"primary_key"`
	StageID      uint64       `gorm:"primary_key"`
	StageRequest StageRequest `gorm:"foreignKey:RequestID; references:ID"`
	Stage        Stage        `gorm:"foreignKey:StageID; references:ID"`
	InputField1  float64      `gorm:"not null; default:0"`
	InputField2  float64      `gorm:"not null; default:0"`
}
