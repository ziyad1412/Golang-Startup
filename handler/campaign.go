package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

// api/v1/campaigns
func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userID)

	if err != nil {
		response := helper.APIResponse("Error to get campaigns", 500, "error", nil)
		c.JSON(500, response)
		return
	}

	response := helper.APIResponse("List of campaigns", 200, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(200, response)
}

// api/v1/campaigns/:id
func (h *campaignHandler) GetCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)

	if err != nil {
		response := helper.APIResponse("Failed to get detail campaign", 400, "error", nil)
		c.JSON(400, response)
		return
	}

	campaignDetail, err := h.service.GetCampaignByID(input)

	if err != nil {
		response := helper.APIResponse("Failed to get detail campaign", 500, "error", nil)
		c.JSON(500, response)
		return
	}

	response := helper.APIResponse("Campaign detail", 200, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(200, response)
}

// api/v1/campaigns post
func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to create campaign", 422, "error", errorMessage)
		c.JSON(422, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	newCampaign, err := h.service.CreateCampaign(input)

	if err != nil {
		response := helper.APIResponse("Failed to create campaign", 500, "error", nil)
		c.JSON(500, response)
		return
	}

	response := helper.APIResponse("Campaign has been created", 201, "success", campaign.FormatCampaign(newCampaign))
	c.JSON(201, response)
}

// update
func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	var inputID campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", 400, "error", nil)
		c.JSON(400, response)
		return
	}

	var inputData campaign.CreateCampaignInput

	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to update campaign", 422, "error", errorMessage)
		c.JSON(422, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	inputData.User = currentUser

	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", 500, "error", nil)
		c.JSON(500, response)
		return
	}
	response := helper.APIResponse("Campaign has been updated", 200, "success", campaign.FormatCampaign(updatedCampaign))
	c.JSON(200, response)
}

// api/v1/campaign-images
func (h *campaignHandler) UploadImage(c *gin.Context) {
	var input campaign.CreateCampaignImageInput

	err := c.ShouldBind(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to upload campaign image", 422, "error", errorMessage)
		c.JSON(422, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	userID := currentUser.ID

	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image", 400, "error", data)
		c.JSON(400, response)
		return
	}

	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image", 500, "error", data)
		c.JSON(500, response)
	}

	_, err = h.service.SaveCampaignImage(input, path)

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image", 500, "error", data)
		c.JSON(500, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Campaign image has been uploaded", 201, "success", data)
	c.JSON(201, response)
}
