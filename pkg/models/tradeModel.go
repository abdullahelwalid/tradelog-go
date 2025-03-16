package models

import (
	"time"

	"gorm.io/gorm"
)


type Trade struct {
	gorm.Model
	TradeId         string    `gorm:"primaryKey;column:trade_id"`
	UserId          string    `gorm:"column:user_id"`
	Asset           string    `gorm:"column:asset"`
	OpenPositionAt  time.Time `gorm:"column:open_position_at"`
	ClosePositionAt time.Time `gorm:"column:close_position_at"`
	Margin          float32   `gorm:"column:margin"`
	OpenPrice       float32   `gorm:"column:open_price"`
	ClosePrice      float32   `gorm:"column:close_price"`
	TradeType       string    `gorm:"column:trade_type"`      // e.g., "buy" or "sell"
	Volume          float32   `gorm:"column:volume"`
	StopLoss        float32   `gorm:"column:stop_loss"`
	TakeProfit      float32   `gorm:"column:take_profit"`
	Profit          float32   `gorm:"column:profit"`
	Notes           string    `gorm:"column:notes"`
	StrategyTag     string    `gorm:"column:strategy_tag"`
	IsSimulated     bool      `gorm:"column:is_simulated"`
}
