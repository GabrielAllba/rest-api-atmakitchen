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

func InsertQuota(c *gin.Context) {
	var quota models.Quota

	// Bind the incoming JSON to the quota struct
	if err := c.ShouldBindJSON(&quota); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if a record with the same product_id and tanggal already exists
	var existingQuota models.Quota
	if err := models.DB.Where("product_id = ? AND tanggal = ?", quota.ProductId, quota.Tanggal).First(&existingQuota).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// No record found, so we can insert the new quota
			if err := models.DB.Create(&quota).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Quota inserted successfully", "quota": quota})
			return
		} else {
			// An error occurred while checking the existing record
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// If we reach here, it means the record already exists
	c.JSON(http.StatusConflict, gin.H{"message": "Quota with the same product_id and tanggal already exists"})
}

func UpdateQuota(c *gin.Context) {
	var quota models.Quota

	// Bind the incoming JSON to the quota struct
	if err := c.ShouldBindJSON(&quota); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if a record with the same product_id and tanggal exists
	var existingQuota models.Quota
	if err := models.DB.Where("product_id = ? AND tanggal = ?", quota.ProductId, quota.Tanggal).First(&existingQuota).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// No record found, so we cannot update a non-existing quota
			c.JSON(http.StatusNotFound, gin.H{"error": "Quota not found"})
			return
		} else {
			// An error occurred while checking the existing record
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Update the existing quota with new values
	existingQuota.Quota = quota.Quota
	// Add other fields to update as needed

	if err := models.DB.Save(&existingQuota).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Quota updated successfully", "quota": existingQuota})
}

func GetQuotaByProductIDAndTanggal(c *gin.Context) {
	productID := c.Param("product_id")
	tanggal := c.Param("tanggal")

	var quota models.Quota
	if err := models.DB.Where("product_id = ? AND tanggal = ?", productID, tanggal).First(&quota).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Quota not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"quota": quota})
}
func GetQuotaByHampersIDAndTanggal(c *gin.Context) {
	hampersID := c.Param("hampers_id")
	tanggal := c.Param("tanggal")

	var quota models.Quota
	if err := models.DB.Where("hampers_id = ? AND tanggal = ?", hampersID, tanggal).First(&quota).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Quota not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"quota": quota})
}

func InsertQuotaHampers(c *gin.Context) {
	var quota models.Quota

	// Bind the incoming JSON to the quota struct
	if err := c.ShouldBindJSON(&quota); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if a record with the same product_id and tanggal already exists
	var existingQuota models.Quota
	if err := models.DB.Where("hampers_id = ? AND tanggal = ?", quota.HampersId, quota.Tanggal).First(&existingQuota).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// No record found, so we can insert the new quota
			if err := models.DB.Create(&quota).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Quota inserted successfully", "quota": quota})
			return
		} else {
			// An error occurred while checking the existing record
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// If we reach here, it means the record already exists
	c.JSON(http.StatusConflict, gin.H{"message": "Quota with the same product_id and tanggal already exists"})
}

func UpdateQuotaHampers(c *gin.Context) {
	var quota models.Quota

	// Bind the incoming JSON to the quota struct
	if err := c.ShouldBindJSON(&quota); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if a record with the same product_id and tanggal exists
	var existingQuota models.Quota
	if err := models.DB.Where("hampers_id = ? AND tanggal = ?", quota.HampersId, quota.Tanggal).First(&existingQuota).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// No record found, so we cannot update a non-existing quota
			c.JSON(http.StatusNotFound, gin.H{"error": "Quota not found"})
			return
		} else {
			// An error occurred while checking the existing record
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Update the existing quota with new values
	existingQuota.Quota = quota.Quota
	// Add other fields to update as needed

	if err := models.DB.Save(&existingQuota).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Quota updated successfully", "quota": existingQuota})
}