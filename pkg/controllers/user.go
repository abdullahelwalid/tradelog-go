package controllers

import (
	"strings"
	"encoding/json"
	"net/http"
	"github.com/abdullahelwalid/tradelog-go/pkg/models"
	"github.com/abdullahelwalid/tradelog-go/pkg/utils"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	// Define the struct to map the form data
	type FormData struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FullName  string `json:"fullName"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}

	// Parse the form data
	var data FormData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		// Set response header to application/json
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// Send error message in JSON
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Cannot parse form data",
		})
		return
	}

	// Validate that all fields are provided
	if data.Email == "" || data.Password == "" || data.FullName == "" || data.FirstName == "" || data.LastName == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// Send error message in JSON
		json.NewEncoder(w).Encode(map[string]string{
			"error": "All fields (email, password, fullName, firstName, lastName) are required",
		})
		return
	}

	// Check if the email already exists in the database
	user := &models.User{}
	utils.DB.First(user, "email = ?", data.Email)
	if user.Email != "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		// Send error message in JSON
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Email already exists",
		})
		return
	}

	// Initialize AWS config
	auth, err := utils.InitAWSConfig()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		// Send error message in JSON
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Something went wrong",
		})
		return
	}

	// Perform the signup operation
	err = auth.Signup(data.Email, data.Password, &data.FirstName, &data.LastName, &data.FullName)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if strings.Contains(err.Error(), "UsernameExistsException") {
			w.WriteHeader(http.StatusConflict)
			// Send error message in JSON
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Email already exists",
			})
		} else {
			w.WriteHeader(http.StatusBadRequest)
			// Send error message in JSON
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
		}
		return
	}

	// Respond with success in JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Send success message in JSON
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Sign up successful",
	})
}


func ConfirmSignUp(w http.ResponseWriter, r *http.Request) {
	// Define the struct to map the form data
	type FormData struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	// Parse the form data
	var data FormData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		// Set response header to application/json
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// Send error message in JSON
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Cannot parse form data",
		})
		return
	}

	// Validate that the code is provided
	if data.Code == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// Send error message in JSON
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Code not in request",
		})
		return
	}

	// Validate that the email is provided
	if data.Email == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// Send error message in JSON
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Email not in request",
		})
		return
	}

	// Initialize AWS config
	auth, err := utils.InitAWSConfig()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// Send error message in JSON
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Perform the confirm sign-up operation
	err = auth.ConfirmSignUp(data.Email, data.Code)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// Send error message in JSON
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Retrieve user attributes after confirmation
	resp, err := auth.AdminGetUser(data.Email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		// Send error message in JSON
		json.NewEncoder(w).Encode(map[string]string{
			"error": "An error has occurred while verifying your account",
		})
		return
	}

	// Extract user attributes
	var firstName, lastName, fullName, username *string
	for _, att := range resp.UserAttributes {
		if *att.Name == "given_name" {
			firstName = att.Value
		}
		if *att.Name == "family_name" {
			lastName = att.Value
		}
		if *att.Name == "name" {
			fullName = att.Value
		}
	}
	username = resp.Username

	// Add user info to database
	user := models.User{
		UserId:   *username,
		Email:    data.Email,
		FirstName: *firstName,
		FullName:  *fullName,
		LastName:  *lastName,
	}
	result := utils.DB.Create(&user)
	if result.Error != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		// Send error message in JSON
		json.NewEncoder(w).Encode(map[string]string{
			"error": "An error has occurred while creating your account",
		})
		return
	}

	// Respond with success in JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Send success message in JSON
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Sign up successful",
	})
}

func ResendConfirmationCode(w http.ResponseWriter, r *http.Request) {
	// Define the struct to map the form data
	type FormData struct {
		Email string `json:"email"`
	}

	// Parse the form data
	var data FormData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		// Set response header to application/json
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// Send error message in JSON
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Cannot parse form data",
		})
		return
	}

	// Validate that email is provided
	if data.Email == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// Send error message in JSON
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Email not in request",
		})
		return
	}

	// Initialize AWS config
	auth, err := utils.InitAWSConfig()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// Send error message in JSON
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Resend the confirmation code
	err = auth.ResendConfirmationCode(data.Email)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		// Send error message in JSON
		json.NewEncoder(w).Encode(map[string]string{
			"error": "An error has occurred while resending confirmation code",
		})
		return
	}

	// Respond with success in JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Send success message in JSON
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Code sent successfully",
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	// Define the struct to map the form data
	type FormData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse the form data
	var data FormData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		// Set response header to application/json
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// Send error message in JSON
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Cannot parse form data",
		})
		return
	}

	// Validate that email is provided
	if data.Email == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// Send error message in JSON
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Email is required",
		})
		return
	}

	// Initialize AWS config
	auth, err := utils.InitAWSConfig()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		// Send error message in JSON
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Perform the login operation
	_, err = auth.Login(data.Email, data.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		// Send error message in JSON
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid credentials",
		})
		return
	}

	// Respond with success message in JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Send success message in JSON
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful! Welcome, " + data.Email + ".",
	})
}
