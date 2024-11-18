package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Admin struct {
	BaseModel
	FullName       string    `json:"full_name" gorm:"not null"`
	Username       string    `json:"username" gorm:"not null"`
	Email          string    `json:"email" gorm:"not null;uniqueIndex;size:191"`
	Password       string    `json:"password" gorm:"not null"`
	Role           string    `json:"role" gorm:"not null"`
	RoomingHouseID uuid.UUID `json:"rooming_house_id" gorm:"not null;size:191"`
}

type AdminRegisterBody struct {
	FullName       string    `json:"full_name" gorm:"not null"`
	Username       string    `json:"username" gorm:"not null"`
	Email          string    `json:"email" gorm:"not null"`
	Password       string    `json:"password" gorm:"not null"`
	RoomingHouseID uuid.UUID `json:"rooming_house_id"`
}

type AdminResponse struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	Role           string    `json:"role"`
	RoomingHouseID uuid.UUID `json:"rooming_house_id"`
}

type GetAllAdminResponse struct {
	ID           uuid.UUID                  `json:"id"`
	Username     string                     `json:"username"`
	Role         string                     `json:"role"`
	RoomingHouse TenantRoomingHouseResponse `json:"rooming_house"`
}

func (a *Admin) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.New()
	a.CreatedAt = time.Now()

	hasedPassword, _ := bcrypt.GenerateFromPassword([]byte(a.Password), 14)

	a.Password = string(hasedPassword)
	return
}
