package transactioncontroller

import (
	"net/http"
	"strconv"

	"backend-atmakitchen/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetByUserID(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	var transactions []models.Transaction
	if err := models.DB.Where("user_id = ?", userID).Preload("User").Find(&transactions).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "transactions not found for this user"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions for this user"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}
    
func Create(c *gin.Context) {
    var transactions models.Transaction

    if err := c.BindJSON(&transactions); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := models.DB.Create(&transactions).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create transactions"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

func Show(c *gin.Context) {
    id := c.Param("id")
    var transactions models.Transaction

    if err := models.DB.Preload("User").First(&transactions, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "transactions not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

func Index(c *gin.Context) {
    var transactions []models.Transaction

    if err := models.DB.Preload("User").Find(&transactions).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

func Delete(c *gin.Context) {
    id := c.Param("id")
    intId, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    var transactions models.Transaction
    if err := models.DB.First(&transactions, intId).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "transactions not found"})
        return
    }

    if err := models.DB.Delete(&transactions).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transactions"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "transactions deleted successfully"})
}

func Update(c *gin.Context) {
    var transactions models.Transaction
    id := c.Param("id")

    if err := c.BindJSON(&transactions); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := models.DB.Model(&models.Transaction{}).Where("id = ?", id).Updates(&transactions).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update transactions"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}
