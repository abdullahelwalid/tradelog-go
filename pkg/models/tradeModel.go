package models

import (
	"time"

	"gorm.io/gorm"
)


type Trade struct {
	gorm.Model
	TradId string `gorm:"primaryKey;column:trad_id"`
	UserId string
	Asset string
	OpenPositionAt time.Time
	ClosePositionAt time.Time
	Margin float32
	OpenPrice float32
	ClosePrice float32
}
