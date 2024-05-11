package tokencontroller

import (
	"backend-atmakitchen/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// func CreateToken(c *gin.Context) {
// 	var token models.Token
// 	db := c.MustGet("db").(*gorm.DB)
// 	userId := c.MustGet("userId").(int)
// 	token.UserId = userId
// 	token.LogoutToken = uuid.New().String()
// 	token.Expiration = "1h"
// 	db.Create(&token)
// 	c.JSON(http.StatusOK, gin.H{"token": token})
// }

func CreateToken(c *gin.Context) {
	var token models.Token

	// Extract user_id from URL parameters
	userIdStr := c.Param("user_id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	token.UserId = userId
	token.LogoutToken = uuid.New().String()
	token.Expiration = "1h"
	// Save the product to the database
	if err := models.DB.Where("user_id").Create(&token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Index(c *gin.Context) {
	var token []models.Token
	models.DB.Find(&token)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func DeleteToken(c *gin.Context) {
	// Extract the product ID from the request parameters
	id := c.Param("user_id")

	// Convert the ID string to an integer
	tokenID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token ID"})
		return
	}

	if err := models.DB.Where("user_id = ?", tokenID).Delete(&models.Token{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token deleted successfully"})
}

func CheckToken(c *gin.Context) {
	// var token models.Token
	// db := c.MustGet("db").(*gorm.DB)
	// userId := c.MustGet("userId").(int)
	// db.Where("user_id = ?", userId).First(&token)
	// c.JSON(http.StatusOK, gin.H{"token": token})
}
