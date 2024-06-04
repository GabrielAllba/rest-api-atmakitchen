package bahanresepcontroller

import (
	"backend-atmakitchen/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func Index(c *gin.Context) {
	var bahan_resep []models.BahanResep
	models.DB.Find(&bahan_resep).Preload("Product").Preload("Bahan")
	c.JSON(http.StatusOK, gin.H{"bahan_resep": bahan_resep})
}

func Show(c *gin.Context) {
    var bahan_resep []models.BahanResep
    product_id := c.Param("product_id")

    // Use Where to filter by product_id and then load all matching records
    if err := models.DB.Preload("Bahan").Preload("Product").Where("product_id = ?", product_id).Find(&bahan_resep).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "bahan_resep not found"})
            return
        } else {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
            return
        }
    }
    c.JSON(http.StatusOK, gin.H{"bahan_resep": bahan_resep})
}
