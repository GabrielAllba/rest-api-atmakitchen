package models

import (
    "time"
)

type WithdrawHistory struct {
    Id        int       `gorm:"primaryKey;" json:"id"`
    UserId    int       `gorm:"index" json:"user_id"`
    User      User      `gorm:"foreignKey:UserId" json:"user"`
    Amount    float64   `json:"amount"`
    CreatedAt time.Time `json:"created_at"`
}
