package productcontroller

import (
	"backend-atmakitchen/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Create(c *gin.Context){
	var product models.Product

	// Bind JSON data to product struct
	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	// validitas semua input terisi
	if product.Description == "" || 
	product.Name == "" ||
	product.Photo == "" || 
	product.Status== "" || 
	product.ProductTypeId == 0 || 
	product.Price == 0 || product.Stock == 0 || 
	product.DailyQuota == 0 || product.RewardPoin == 0  {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pastikan semua input terisi"})
		return
	}

	// check role id in database
	var product_type models.ProductType
	if err := models.DB.Where("id = ?", product.Id).First(&product_type).Error; err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipe Product tidak tersedia"})
		return
	}

	if product.ProductType.Name == "Titipan"{
		var consignation models.Consignation
		if err := models.DB.Where("id = ?", product.ConsignationId).First(&consignation).Error; err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": "Penitip tidak tersedia"})
			return
		}	
	}
	
	if err := models.DB.Create(&product).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat Product"})
		return
	}
	var returnProduct models.Product
	if err := models.DB.Preload("Consignation").Preload("ProductType").First(&returnProduct, product.Id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": returnProduct})

}

