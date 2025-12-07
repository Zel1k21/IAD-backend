package ds

import "time"

type StageRequest struct {
	ID                  uint64                `gorm:"primary_key"`
	Status              uint8                 `gorm:"not null; default:1"`
	UserID              uint64                `gorm:"not null"`
	User                User                  `gorm:"foreignKey:UserID; references:ID"`
	ModeratorID         uint64                `gorm:"default:null"`
	Morderator          User                  `gorm:"foreignKey:ModeratorID; references:ID"`
	CreatedAt           time.Time             `gorm:"not null; default:now()"`
	FormedAt            time.Time             `gorm:"default:null"`
	ClosedAt            time.Time             `gorm:"default:null"`
	StageRequestToStage []StageRequestToStage `gorm:"foreignKey:RequestID; references:ID"`
	ProductName         string                `gorm:"type:varchar(100); not null"`
	CalculationResult   float64
}
