package models

import (
	"time"
	"github.com/google/uuid"
)
type Plan string 

const (
	Monthly Plan = "monthly"
	HalfYearly Plan = "half_yearly"
	Yearly Plan = "yearly"
)

type Subscription struct {
	ID     uuid.UUID 	 `json:"id"`
	ServiceName string   `json:"service_name"`
	UserID     uuid.UUID `json:"user_id"`
	Price int            `json:"price"`
	Plan Plan            `json:"plan"`
	StartDate  time.Time `json:"start_date"`
	EndDate time.Time    `json:"end_date"`
}