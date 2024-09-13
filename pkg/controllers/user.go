package controllers

import (
	"fmt"
	"strings"

	"github.com/abdullahelwalid/tradelog-go/pkg/models"
	"github.com/abdullahelwalid/tradelog-go/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func SignUp(c *fiber.Ctx) error {
	type FormData struct {
			Email    string `form:"email"`
			Password string `form:"password"`
			FullName string `form:"fullName"`
			FirstName string `form:"firstName"`
			LastName string `form:"lastName"`
		}
	data := new(FormData)
	if err := c.BodyParser(data); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse form data",
		})
	}
	fmt.Println(data)
	auth, err := utils.InitAWSConfig()
	if err != nil {
		c.SendStatus(500)
		return c.SendString("Something went wrong")
	}
	err = auth.Signup(data.Email, data.Password, &data.FirstName, &data.LastName, &data.FullName) 
	if err != nil {
		if (strings.Contains(err.Error(), "UsernameExistsException")) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Email already exist",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			})	
		}
	return c.SendString("sign up successful")
}

func ConfirmSignUp(c *fiber.Ctx) error {
	type FormData struct {
		Email string `form:"email"`
		Code string `form:"code"`
	}
	data := new(FormData)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse form data",
		})
	}
	auth, err := utils.InitAWSConfig()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})	
	}
	err = auth.ConfirmSignUp(data.Email, data.Code)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			})	
		}
	resp, err := auth.AdminGetUser(data.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "An error has occurred while verifying your account",
		})
	}
	var firstName, lastName, fullName *string
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
	userId := uuid.New()	
	//Add user info to database
	user := models.User{UserId: userId.String(), Email: data.Email, FirstName: *firstName, FullName: *fullName, LastName: *lastName}
	result := utils.DB.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "An error has occurred while creating your account",
		})
	}

	return c.SendString("sign up successful")
}

func ResendConfirmationCode(c *fiber.Ctx) error {
	type FormData struct {
		Email string `form:"email"`
	}
	data := new(FormData)
	if err := c.BodyParser(&data); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse form data",
			})
	}
	auth, err := utils.InitAWSConfig()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})	}
	err = auth.ResendConfirmationCode(data.Email)
	if err != nil {
		c.SendStatus(fiber.StatusForbidden)
		return c.JSON(fiber.Map{"error": "An error has occurred while resending confirmation code"})
	}

	c.SendStatus(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"message": "code sent successfully",
	})	
}

func Login(c *fiber.Ctx) error {
	type FormData struct {
		Email string `form:"email"`
		Password string `form:"password"`
	}
	var data FormData
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse form data",
			})
		}
	auth, err := utils.InitAWSConfig()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})	}
	resp, err := auth.Login(data.Email, data.Password)
	if err != nil {
		c.SendStatus(401)
		return c.JSON(fiber.Map{"error": "Invalid credentials"})
	}
	return c.JSON(resp)
}
