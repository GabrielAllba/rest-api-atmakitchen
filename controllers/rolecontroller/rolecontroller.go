package rolecontroller

import (
	"backend-atmakitchen/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)



func Create(c *gin.Context) {
	var role models.Role

	
	role.Name = c.PostForm("name")
	gajiHarianStr := c.PostForm("gaji_harian")
	
	
	if role.Name == "" || gajiHarianStr == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Pastikan semua input terisi"})
		return
	}

	nominal, err := strconv.Atoi(gajiHarianStr)
	if(err != nil){
		c.JSON(http.StatusBadRequest, gin.H{"error": "Total point invalid"})
		return
	}
	role.GajiHarian = float32(nominal)

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
