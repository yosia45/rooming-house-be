package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Period struct {
	BaseModel
	Name              string             `json:"name" gorm:"not null"`
	Unit              string             `json:"unit" gorm:"not null"`
	PeriodPackages    []PeriodPackage    `json:"period_packages" gorm:"foreignKey:PeriodID"`
	Tenants           []Tenant           `json:"tenants" gorm:"foreignKey:PeriodID"`
	AdditionalPeriods []AdditionalPeriod `json:"additional_periods" gorm:"foreignKey:PeriodID"`
}

type PeriodResponse struct {
	ID   uuid.UUID `json:"id" gorm:"column:period_id"`
	Name string    `json:"name" gorm:"column:period_name"`
	Unit string    `json:"unit" gorm:"column:period_unit"`
}

func (p *Period) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	p.CreatedAt = time.Now()

	return
}
