package rolecontroller

import (
	"backend-atmakitchen/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)



func Create(c *gin.Context) {
	var role models.Role

	
	// Bind JSON data to user struct
	if err := c.BindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	
	if role.Name == "" || strconv.Itoa(int(role.GajiHarian)) == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Pastikan semua input terisi"})
		return
	}



	var existingRole models.Role

	if err := models.DB.Where("name = ?", role.Name).First(&existingRole).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role dengan nama tersebut sudah ada"})
		return
	}

	
	if err := models.DB.Create(&role).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": role})
}
