package ds

type StageRequest struct {
	ID                  uint64                `gorm:"primary_key"`
	Status              uint8                 `gorm:"not null; default:1"`
	UserID              uint64                `gorm:"not null"`
	User                User                  `gorm:"foreignKey:UserID; references:ID"`
	ModeratorID         uint64                `gorm:"default:null"`
	Moderator           User                  `gorm:"foreignKey:ModeratorID; references:ID"`
	CreatedAt           string                `gorm:"not null; default:now()"`
	FormedAt            string                `gorm:"default:null"`
	ClosedAt            string                `gorm:"default:null"`
	StageRequestToStage []StageRequestToStage `gorm:"foreignKey:RequestID; references:ID"`
	ProductName         string                `gorm:"type:varchar(100); not null"`
	CalculationResult   float64
}
