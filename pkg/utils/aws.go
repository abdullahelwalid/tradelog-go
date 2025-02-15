package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)


type awsInitializer interface {
	initfunc() (any any)
}

type CognitoAuth struct {
	Cfg aws.Config
	UserPoolID      string
	AppClientID     string
	AppClientSecret string
}

func InitAWSConfig() (*CognitoAuth, error){
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	awsAuth := &CognitoAuth{
		Cfg: cfg,
		UserPoolID:      os.Getenv("COGNITO_USER_POOL_ID"),
		AppClientID:     os.Getenv("COGNITO_APP_CLIENT_ID"),
		AppClientSecret: os.Getenv("COGNITO_APP_CLIENT_SECRET"),
	}
	return awsAuth, nil
}

func computeSecretHash(clientSecret string, username string, clientId string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientId))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (c *CognitoAuth) AdminGetUser(username string) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	client := cognitoidentityprovider.NewFromConfig(c.Cfg)
	userInput := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: &c.UserPoolID,
		Username: &username,
	}
	resp, err := client.AdminGetUser(context.TODO(), userInput)
	if (err != nil) {
		fmt.Println(err)
	}
	return resp, err
}

func (c *CognitoAuth) ValidateToken(token string) (*cognitoidentityprovider.GetUserOutput, error) {
	client := cognitoidentityprovider.NewFromConfig(c.Cfg)
	getUserInputFields := &cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(token),
	}
	resp, err := client.GetUser(context.TODO(), getUserInputFields)
	if err == nil {
		fmt.Println(*resp.Username)	
	}
	return resp, err
}

func (c *CognitoAuth) RefreshToken(refreshToken string, email string) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	client := cognitoidentityprovider.NewFromConfig(c.Cfg)
	authParams := map[string]string{"USERNAME": refreshToken, "SECRET_HASH": computeSecretHash(c.AppClientSecret, email, c.AppClientID)}
	signInInput := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: "REFRESH_TOKEN",
		ClientId: &c.AppClientID,
		AuthParameters: authParams,
	}
	resp, err := client.InitiateAuth(context.TODO(), signInInput)
	return resp, err
}

func (c *CognitoAuth) Login(email string, password string) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	client := cognitoidentityprovider.NewFromConfig(c.Cfg)
	authParams := map[string]string{"USERNAME": email, "PASSWORD": password, "SECRET_HASH": computeSecretHash(c.AppClientSecret, email, c.AppClientID)}
	signInInput := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: "USER_PASSWORD_AUTH",
		ClientId: &c.AppClientID,
		AuthParameters: authParams,
	}
	resp, err := client.InitiateAuth(context.TODO(), signInInput)
	return resp, err
}

func (c *CognitoAuth) ConfirmSignUp(email string, code string) error {
	client := cognitoidentityprovider.NewFromConfig(c.Cfg)
	confirmSignUpInput := &cognitoidentityprovider.ConfirmSignUpInput{
		Username: aws.String(email),
		ClientId: &c.AppClientID,
		ConfirmationCode: aws.String(code),
		SecretHash: aws.String(computeSecretHash(c.AppClientSecret, email, c.AppClientID)),
	}
	resp, err := client.ConfirmSignUp(context.TODO(), confirmSignUpInput)
	fmt.Println(resp)
	return err
}

func (c *CognitoAuth) ResendConfirmationCode(email string) error {
	client := cognitoidentityprovider.NewFromConfig(c.Cfg)
	input := cognitoidentityprovider.ResendConfirmationCodeInput{
		ClientId: &c.AppClientID,
		Username: &email,
		SecretHash: aws.String(computeSecretHash(c.AppClientSecret, email, c.AppClientID)),
	}
	_, err := client.ResendConfirmationCode(context.TODO(), &input)
	return err
}

func (c *CognitoAuth) Signup(email string, password string, firstName *string, lastName *string, fullName *string) error {
	client := cognitoidentityprovider.NewFromConfig(c.Cfg)
	type signUpAttributes struct {
		givenName string
		familyName string
		name string
	}
	signUpAttributesInput := signUpAttributes{givenName: "given_name", familyName: "family_name", name: "name"}
	firstname := types.AttributeType{Name: &signUpAttributesInput.givenName, Value: firstName}	
	lastname := types.AttributeType{Name: &signUpAttributesInput.familyName, Value: lastName}
	fullname := types.AttributeType{Name: &signUpAttributesInput.name, Value: fullName}
	userAttributes := []types.AttributeType{firstname, lastname, fullname}
	signUpInput := &cognitoidentityprovider.SignUpInput{
		ClientId: &c.AppClientID,
		Username: aws.String(email),
		Password: aws.String(password),
		SecretHash: aws.String(computeSecretHash(c.AppClientSecret, email, c.AppClientID)),
		UserAttributes: userAttributes,
	}
	resp, err := client.SignUp(context.TODO(), signUpInput)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(resp)
	return nil
}

func (c *CognitoAuth) ForgotPassword(email *string) error {
	client := cognitoidentityprovider.NewFromConfig(c.Cfg)
	forgotPasswordInput := &cognitoidentityprovider.ForgotPasswordInput{
		Username: email,
		ClientId: &c.AppClientID,
		SecretHash: aws.String(computeSecretHash(c.AppClientSecret, *email, c.AppClientID)),
	}
	resp, err := client.ForgotPassword(context.TODO(), forgotPasswordInput)
	fmt.Println(resp)
	return err
}

func (c *CognitoAuth) ConfirmForgotPassword(code *string, password *string, email *string) error {
	client := cognitoidentityprovider.NewFromConfig(c.Cfg)
	confirmPasswordInput := &cognitoidentityprovider.ConfirmForgotPasswordInput{
		ClientId: &c.AppClientID,
		ConfirmationCode: code,
		Password: password,
		SecretHash: aws.String(computeSecretHash(c.AppClientSecret, *email, c.AppClientID)),
		Username: email,
	}
	resp, err := client.ConfirmForgotPassword(context.TODO(), confirmPasswordInput)
	fmt.Println(resp)
	return err
}
