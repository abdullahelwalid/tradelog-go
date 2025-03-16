package controllers

import (
	"encoding/json"
	"net/http"

	"time"

	"github.com/abdullahelwalid/tradelog-go/pkg/models"
	"github.com/abdullahelwalid/tradelog-go/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)


var validate = validator.New()

type tradeRequest struct {
	Asset           string    `json:"asset" validate:"required"`
	OpenPositionAt  time.Time `json:"openPositionAt" validate:"required"`
	ClosePositionAt time.Time `json:"closePositionAt" validate:"required"`
	Margin          float32   `json:"margin" validate:"required"`
	OpenPrice       float32   `json:"openPrice" validate:"required"`
	ClosePrice      float32   `json:"closePrice" validate:"required"`
	TradeType       string    `json:"tradeType" validate:"required"`
	Volume          float32   `json:"volume" validate:"required"`
}

func AddTrade(w http.ResponseWriter, r *http.Request) {
	username, _ := r.Context().Value("username").(string)

	var tradeReq tradeRequest
	if err := json.NewDecoder(r.Body).Decode(&tradeReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate input
	if err := validate.Struct(tradeReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Save to DB
	trade := models.Trade{
		TradeId: uuid.New().String(),
		UserId:         username,
		Asset:          tradeReq.Asset,
		OpenPositionAt: tradeReq.OpenPositionAt,
		ClosePositionAt: tradeReq.ClosePositionAt,
		Margin:         tradeReq.Margin,
		OpenPrice:      tradeReq.OpenPrice,
		ClosePrice:     tradeReq.ClosePrice,
		TradeType:      tradeReq.TradeType,
		Volume:         tradeReq.Volume,
	}
	if err := utils.DB.Create(&trade).Error; err != nil {
		http.Error(w, "Could not save trade", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Trade successfully added",
	})
}
