package routes

import (
	"net/http"

	"github.com/abdullahelwalid/tradelog-go/pkg/controllers"
	"github.com/abdullahelwalid/tradelog-go/pkg/middleware"
)


var Mux = func() (*http.ServeMux){
	mux := http.NewServeMux()

	//public routes
	mux.HandleFunc("/", middleware.MethodCheckMiddleware(http.HandlerFunc(controllers.TestHandler), []string{http.MethodGet}))
	mux.HandleFunc("/signup", middleware.MethodCheckMiddleware(http.HandlerFunc(controllers.SignUp), []string{http.MethodPost}))
	mux.HandleFunc("/confirmsignup", middleware.MethodCheckMiddleware(http.HandlerFunc(controllers.ConfirmSignUp), []string{http.MethodPost}))
	mux.HandleFunc("/resendConfirmationCode", middleware.MethodCheckMiddleware(http.HandlerFunc(controllers.ResendConfirmationCode), []string{http.MethodPost}))
	mux.HandleFunc("/login", middleware.MethodCheckMiddleware(http.HandlerFunc(controllers.Login), []string{http.MethodPost}))

	// Protected routes
	mux.Handle("/auth", middleware.MethodCheckMiddleware(middleware.AuthenticationMiddleware(http.HandlerFunc(controllers.AuthHandler)), []string{http.MethodGet}))
	mux.Handle("/trade", middleware.MethodCheckMiddleware(middleware.AuthenticationMiddleware(http.HandlerFunc(controllers.AddTrade)), []string{http.MethodPost}))
		
	return mux
}
