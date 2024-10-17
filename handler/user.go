package handler

import (
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	//tangkap input dari user
	//map input dari user ke struct RegisterUserInput
	//struct di atas akan dikirim ke service

	var input user.RegisterUserInput

	err := c.ShouldBindBodyWithJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Register account failed", 422, "error", errorMessage)
		c.JSON(422, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)

	if err != nil {
		response := helper.APIResponse("Register account failed", 400, "error", nil)
		c.JSON(400, response)
		return
	}

	//token, err := h.jwtService.GenerateToken()

	formatter := user.FormatUser(newUser, "token")

	response := helper.APIResponse("Account has been registered", 200, "success", formatter)

	c.JSON(200, response)
}

func (h *userHandler) Login(c *gin.Context) {

	var input user.LoginInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Login failed", 422, "error", errorMessage)
		c.JSON(422, response)
		return
	}

	loggedInUser, err := h.userService.Login(input)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Login failed", 422, "error", errorMessage)
		c.JSON(422, response)
		return
	}

	formatter := user.FormatUser(loggedInUser, "token")

	response := helper.APIResponse("Successfuly logged in", 200, "success", formatter)

	c.JSON(200, response)
}

// check email availability
func (h *userHandler) CheckEmailAvailability(c *gin.Context) {

	var input user.CheckEmailInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Email checking failed", 422, "error", errorMessage)
		c.JSON(422, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)

	if err != nil {
		errorMessage := gin.H{"errors": "Server error"}

		response := helper.APIResponse("Email checking failed", 500, "error", errorMessage)
		c.JSON(500, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	metaMessage := "Email has been registered"

	if isEmailAvailable {
		metaMessage = "Email is available"
	}

	response := helper.APIResponse(metaMessage, 200, "success", data)
	c.JSON(200, response)
}

// Upload Avatar
func (h *userHandler) UploadAvatar(c *gin.Context) {
	//input dari user
	//simpan gambar di folder "images/"
	//service -> panggil repo
	//JWT (sementara hardcode, seakan-akan user yang login ID = 1)
	//repo ambil data user yang ID = 1
	//repo update data user simpan lokasi file

	// get file from form-data
	file, err := c.FormFile("avatar")

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", 400, "error", data)
		c.JSON(400, response)
		return
	}

	// get user ID from JWT
	userID := 1

	// save file to path folder "images/"
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)

	if err != nil {
		// data error disimpan di helper
		data := gin.H{"is_uploaded": false}
		// response error disimpan di helper
		response := helper.APIResponse("Failed to upload avatar image", 400, "error", data)
		// response error dikirim ke user
		c.JSON(400, response)
		return
	}

	// call service
	_, err = h.userService.SaveAvatar(userID, path)

	if err != nil {
		// data error disimpan di helper
		data := gin.H{"is_uploaded": false}
		// response error disimpan di helper
		response := helper.APIResponse("Failed to upload avatar image", 400, "error", data)
		// response error dikirim ke user
		c.JSON(400, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Avatar successfuly uploaded", 200, "success", data)

	c.JSON(200, response)

}
