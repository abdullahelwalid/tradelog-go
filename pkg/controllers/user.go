package controllers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/abdullahelwalid/tradelog-go/pkg/models"
	authTypes "github.com/abdullahelwalid/tradelog-go/pkg/types"
	"github.com/abdullahelwalid/tradelog-go/pkg/utils"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
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
		var errTR *types.TooManyRequestsException
		var errTMFA *types.TooManyFailedAttemptsException
		w.Header().Set("Content-Type", "application/json")
		if errors.As(err, &errTR){
			w.WriteHeader(http.StatusTooManyRequests)
			// Send error message in JSON
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Too Many Requests",
			})
			return
		}
		if errors.As(err, &errTMFA){
			w.WriteHeader(http.StatusUnauthorized)
			// Send error message in JSON
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Too Many Failed Attempts",
			})
		}
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
		AuthFlow string `json:"authFlow"`
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
	resp, err := auth.Login(data.Email, data.Password)
	if err != nil {
		var errIC *types.InvalidPasswordException
		var errTR *types.TooManyRequestsException
		var errTMFA *types.TooManyFailedAttemptsException
		w.Header().Set("Content-Type", "application/json")
		if errors.As(err, &errIC){
			w.WriteHeader(http.StatusUnauthorized)
			// Send error message in JSON
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid credentials",
			})
			return
		}
		if errors.As(err, &errTR){
			w.WriteHeader(http.StatusTooManyRequests)
			// Send error message in JSON
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Too Many Requests",
			})
			return
		}
		if errors.As(err, &errTMFA){
			w.WriteHeader(http.StatusUnauthorized)
			// Send error message in JSON
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Too Many Failed Attempts",
			})
		}
		w.WriteHeader(http.StatusUnauthorized)
		// Send error message in JSON
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid credentials",
		})
		return
	}


	user := &models.User{}
	utils.DB.First(user, "email = ?", data.Email)

	if user.Email == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		// Send error message in JSON
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid credentials",
		})
		return
	}



	
	authData := &authTypes.AuthCookies{
		RefreshToken: *resp.AuthenticationResult.RefreshToken,
		AccessToken: *resp.AuthenticationResult.AccessToken,
		Email: user.Email,
		Username: user.UserId,
	}
	authDataSerialized, err := json.Marshal(authData)
	if err != nil{
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		// Send error message in JSON
		json.NewEncoder(w).Encode(map[string]string{
			"error": "An error has occurred while generating auth data",
		})
		return

	}
	authDataCookie := &http.Cookie{
		Name: "authData",
		Value: base64.StdEncoding.EncodeToString(authDataSerialized),
		Expires: time.Now().AddDate(0, 0, 30),
	}
	http.SetCookie(w, authDataCookie)
	// Respond with success message in JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Send success message in JSON
	if (data.AuthFlow == ""){
		json.NewEncoder(w).Encode(map[string]string{
			"userId": user.UserId,
			"firstName": user.FirstName,
			"lastName": user.LastName,
			"fullName": user.FullName,
		})
	} else {
		json.NewEncoder(w).Encode(map[string]string{
			"userId": user.UserId,
			"firstName": user.FirstName,
			"lastName": user.LastName,
			"fullName": user.FullName,
			"accessToken": *resp.AuthenticationResult.AccessToken,
		})
	}
	return
}


func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	type formData struct{
		Email string `json:"email"`
	}

	var data formData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Cannot parse JSON data",
		})
	}
	auth, err := utils.InitAWSConfig()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal Server Error",
		})
	}

	err = auth.ForgotPassword(&data.Email)
	

	if err != nil {
		var errTR *types.TooManyRequestsException
		var errTMFA *types.TooManyFailedAttemptsException
		w.Header().Set("Content-Type", "application/json")
		if errors.As(err, &errTR){
			w.WriteHeader(http.StatusTooManyRequests)
			// Send error message in JSON
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Too Many Requests",
			})
			return
		}
		if errors.As(err, &errTMFA){
			w.WriteHeader(http.StatusUnauthorized)
			// Send error message in JSON
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Too Many Failed Attempts",
			})
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal Server Error",
		})
	}
}

func GetProfile(w http.ResponseWriter, r *http.Request){
	username, _ := r.Context().Value("username").(string)
	user := &models.User{UserId: username}
	utils.DB.First(user)
	w.Header().Set("Content-Type", "application/json")
	if (user.Email == ""){
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal Server Error user not found",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"email": user.Email,
		"firstName": user.FirstName,
		"lastName": user.LastName,
		"fullName": user.FullName,
		"profilePictureURL": user.ProfileUrl,
	})
	return
}
