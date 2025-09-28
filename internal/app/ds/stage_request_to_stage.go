package ds

type StageRequestToStage struct {
	RequestID              uint64       `gorm:"primary_key"`
	StageID                uint64       `gorm:"primary_key"`
	StageRequest           StageRequest `gorm:"foreignKey:RequestID; references:ID"`
	Stage                  Stage        `gorm:"foreignKey:StageID; references:ID"`
	InputField1            uint64       `gorm:"not null; default:0"`
	InputField2            uint64       `gorm:"not null; default:0"`
	StageCalculationResult uint64       `gorm:"not null; default:0"`
}
