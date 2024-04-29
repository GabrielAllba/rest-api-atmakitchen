package titipancontroller

import (
	"backend-atmakitchen/models"
	"net/http"

	"github.com/gin-gonic/gin"
)



func Create(c *gin.Context) {
	var consignation models.Consignation

	
	if err := c.BindJSON(&consignation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
		
	if err := models.DB.Create(&consignation).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat consignation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"consignation": consignation})
}
