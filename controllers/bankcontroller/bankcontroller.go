package bankcontroller

import (
	"backend-atmakitchen/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	var bank []models.Bank
	models.DB.Find(&bank)
	c.JSON(http.StatusOK, gin.H{"bank": bank})
}

func Create(c *gin.Context) {
	var banks []models.Bank

	
	if err := c.BindJSON(&banks); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	
	for _, bank := range banks {
		if err := models.DB.Create(&bank).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat bank"})
			return
		}
	}
	// Respond with the created banks
	c.JSON(http.StatusOK, gin.H{"banks": banks})
}
