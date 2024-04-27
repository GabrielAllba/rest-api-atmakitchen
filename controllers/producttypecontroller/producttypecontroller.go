package producttypecontroller

import (
	"backend-atmakitchen/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)



func Create(c *gin.Context) {
	var product_type models.ProductType

	
	if err := c.BindJSON(&product_type); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	
	if product_type.Name == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Pastikan semua input terisi"})
		return
	}

	var existingProductType models.ProductType

	if err := models.DB.Where("name = ?", product_type.Name).First(&existingProductType).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product type dengan nama tersebut sudah ada"})
		return
	}

	
	if err := models.DB.Create(&product_type).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat product type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product_type": product_type})
}

func Index(c *gin.Context) {
	var product_type []models.ProductType
	models.DB.Find(&product_type)
	c.JSON(http.StatusOK, gin.H{"product_type": product_type})
}

func Show(c *gin.Context) {
	var product_type models.ProductType
	id := c.Param("id")
	if err := models.DB.First(&product_type, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Product Type not found"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"product type": product_type})
}