package quotacontroller

import (
	"net/http"
	"strconv"

	"backend-atmakitchen/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetByProductAndDate(c *gin.Context) {
    productID := c.Query("product_id")
    tanggal := c.Query("tanggal")

    var quotas []models.Quota
    if err := models.DB.Where("product_id = ? AND tanggal = ?", productID, tanggal).Find(&quotas).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"quotas": quotas})
}

func GetByHampersAndDate(c *gin.Context) {
    hampersID := c.Query("hampers_id")
    tanggal := c.Query("tanggal")

    var quotas []models.Quota
    if err := models.DB.Where("hampers_id = ? AND tanggal = ?", hampersID, tanggal).Find(&quotas).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"quotas": quotas})
}
    
func Create(c *gin.Context) {
    var quota models.Quota

    if err := c.BindJSON(&quota); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := models.DB.Create(&quota).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create quota"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"quota": quota})
}

func Show(c *gin.Context) {
    id := c.Param("id")
    var quota models.Quota

    if err := models.DB.Preload("Product").Preload("Hampers").First(&quota, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "quota not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch quota"})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{"quota": quota})
}

func Index(c *gin.Context) {
    var quota []models.Quota

    if err := models.DB.Preload("Hampers").Preload("Product").Find(&quota).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"quota": quota})
}

func Delete(c *gin.Context) {
    id := c.Param("id")
    intId, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    var quota models.Quota
    if err := models.DB.First(&quota, intId).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "quota not found"})
        return
    }

    if err := models.DB.Delete(&quota).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete quota"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "quota deleted successfully"})
}

func Update(c *gin.Context) {
    var quota models.Quota
    id := c.Param("id")

    if err := c.BindJSON(&quota); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := models.DB.Model(&models.Quota{}).Where("id = ?", id).Updates(&quota).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update quota"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"quota": quota})
}
