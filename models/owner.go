package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Owner struct {
	BaseModel
	FullName      string         `json:"full_name" gorm:"not null"`
	Username      string         `json:"username" gorm:"not null"`
	Email         string         `json:"email" gorm:"not null;uniqueIndex;size:191"`
	Password      string         `json:"password" gorm:"not null"`
	Role          string         `json:"role" gorm:"not null"`
	RoomingHouses []RoomingHouse `json:"rooming_houses" gorm:"foreignKey:OwnerID"`
}

type OwnerRegisterBody struct {
	FullName string `json:"full_name" gorm:"not null"`
	Username string `json:"username" gorm:"not null"`
	Email    string `json:"email" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
}

type OwnerResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
}

func (o *Owner) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.New()
	o.CreatedAt = time.Now()

	hasedPassword, _ := bcrypt.GenerateFromPassword([]byte(o.Password), 14)

	o.Password = string(hasedPassword)
	return
}
