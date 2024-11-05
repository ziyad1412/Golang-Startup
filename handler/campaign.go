package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
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
