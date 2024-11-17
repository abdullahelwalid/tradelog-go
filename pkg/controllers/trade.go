package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/abdullahelwalid/tradelog-go/pkg/models"
	"github.com/abdullahelwalid/tradelog-go/pkg/utils"
	"github.com/google/uuid"
)

func AddTrade(w http.ResponseWriter, r *http.Request) {
	// Define the struct to map the form data
	type FormData struct {
		Asset          string    `json:"asset"`
		OpenPositionAt time.Time `json:"openPositionAt"`
		ClosePositionAt time.Time `json:"closePositionAt"`
		Margin         float32   `json:"margin"`
		OpenPrice      float32   `json:"openPrice"`
		ClosePrice     float32   `json:"closePrice"`
	}

	// Check if Content-Type is application/x-www-form-urlencoded
	reqHeaders := r.Header
	fmt.Println(reqHeaders)
	if !slices.Contains(reqHeaders["Content-Type"], "application/x-www-form-urlencoded") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// Return error in JSON
		json.NewEncoder(w).Encode(map[string]string{"error": "Unsupported content type"})
		return
	}

	// Parse the request body into FormData
	var data FormData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// Return error in JSON
		json.NewEncoder(w).Encode(map[string]string{"error": "Cannot parse request payload"})
		return
	}

	// Validate required fields
	if data.Asset == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// Return error in JSON
		json.NewEncoder(w).Encode(map[string]string{"error": "Asset is required"})
		return
	}
	if data.OpenPositionAt.IsZero() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// Return error in JSON
		json.NewEncoder(w).Encode(map[string]string{"error": "OpenPositionAt is required"})
		return
	}
	if data.ClosePositionAt.IsZero() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// Return error in JSON
		json.NewEncoder(w).Encode(map[string]string{"error": "ClosePositionAt is required"})
		return
	}

	// Additional validation checks for Margin, OpenPrice, and ClosePrice
	if data.Margin <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// Return error in JSON
		json.NewEncoder(w).Encode(map[string]string{"error": "Margin must be greater than 0"})
		return
	}
	if data.OpenPrice <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// Return error in JSON
		json.NewEncoder(w).Encode(map[string]string{"error": "OpenPrice must be greater than 0"})
		return
	}
	if data.ClosePrice <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// Return error in JSON
		json.NewEncoder(w).Encode(map[string]string{"error": "ClosePrice must be greater than 0"})
		return
	}

	// Generate a trade ID and get the user ID
	tradeId := uuid.New()
	userId, _ := r.Context().Value("username").(string)

	// Create the trade model
	trade := &models.Trade{
		TradId:         tradeId.String(),
		Asset:          data.Asset,
		OpenPositionAt: data.OpenPositionAt,
		ClosePositionAt: data.ClosePositionAt,
		Margin:         data.Margin,
		OpenPrice:      data.OpenPrice,
		ClosePrice:     data.ClosePrice,
		UserId:         userId,
	}

	// Create the trade in the database
	result := utils.DB.Create(trade)
	if result.Error != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		// Return error in JSON
		json.NewEncoder(w).Encode(map[string]string{"error": "An error occurred while adding the trade"})
		return
	}

	// Return success with the trade ID in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// Return the created trade ID in JSON
	json.NewEncoder(w).Encode(map[string]string{"id": tradeId.String()})
}

