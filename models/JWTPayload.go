package models

import "github.com/google/uuid"

type JWTPayload struct {
	UserID         uuid.UUID `json:"user_id"`
	Role           string    `json:"role"`
	RoomingHouseID uuid.UUID `json:"rooming_house_id"`
}
