package models

import (
	"time"

	"github.com/google/uuid"
)

type Space struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Code      string    `json:"code" db:"code"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Member struct {
	ID      int       `json:"id" db:"id"`
	SpaceID uuid.UUID `json:"space_id" db:"space_id"`
	IsHost  bool      `json:"isHost" db:"is_host"`
}

type Media struct {
	ID       uuid.UUID `json:"id" db:"id"`
	SpaceID  uuid.UUID `json:"space_id" db:"space_id"`
	Location string    `json:"location" db:"location"`
	HashID   string    `json:"hash_id" db:"hash_id"`
}

type MediaHash struct {
	Hash string `json:"hash" db:"hash"`
}