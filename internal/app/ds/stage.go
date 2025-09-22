package ds

type Stage struct {
	ID                   uint64                `gorm:"primary_key"`
	Title                string                `gorm:"type:varchar(100);not null"`
	ImageURL             string                `gorm:"type:varchar(100);"`
	Description          string                `gorm:"type:varchar(300);"`
	FirstDimensionName   string                `gorm:"type:varchar(30);"`
	FirstDimensionConst  float64               `gorm:"type:float;not null"`
	SecondDimensionName  string                `gorm:"type:varchar(30);"`
	SecondDimensionConst float64               `gorm:"type:float;not null"`
	IsDeleted            bool                  `gorm:"boolean; not null; default: false"`
	StageRequestToStages []StageRequestToStage `gorm:"foreignKey:StageID; references:ID"`
}
