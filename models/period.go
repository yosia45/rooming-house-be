package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Period struct {
	BaseModel
	Name           string          `json:"name" gorm:"not null"`
	Unit           string          `json:"unit" gorm:"not null"`
	PeriodPackages []PeriodPackage `json:"period_packages" gorm:"foreignKey:PeriodID"`
	Tenants        []Tenant        `json:"tenants" gorm:"foreignKey:PeriodID"`
}

func (p *Period) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	p.CreatedAt = time.Now()

	return
}
