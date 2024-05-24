package transactiondetailcontroller

import (
	"net/http"
	"strconv"

	"backend-atmakitchen/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetByInvoiceNumber(c *gin.Context) {
	invoiceNumber := c.Param("invoiceNumber")
	
	var transaction_details []models.TransactionDetail
	if err := models.DB.Where("invoice_number = ?", invoiceNumber).Find(&transaction_details).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "transaction_details not found for this user"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction_details for this user"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction_details": transaction_details})
}
    
func Create(c *gin.Context) {
    
    var transaction_details models.TransactionDetail
    
    // start generate invoice number
    
    // end generate invoice number

    if err := c.BindJSON(&transaction_details); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := models.DB.Create(&transaction_details).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create transaction_details"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"transaction_details": transaction_details})
}

func Show(c *gin.Context) {
    id := c.Param("id")
    var transaction_details models.TransactionDetail

    if err := models.DB.First(&transaction_details, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "transaction_details not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction_details"})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{"transaction_details": transaction_details})
}

func Index(c *gin.Context) {
    var transaction_details []models.TransactionDetail

    if err := models.DB.Find(&transaction_details).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"transaction_details": transaction_details})
}

func Delete(c *gin.Context) {
    id := c.Param("id")
    intId, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    var transaction_details models.TransactionDetail
    if err := models.DB.First(&transaction_details, intId).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "transaction_details not found"})
        return
    }

    if err := models.DB.Delete(&transaction_details).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction_details"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "transaction_details deleted successfully"})
}

func Update(c *gin.Context) {
    var transaction_details models.TransactionDetail
    id := c.Param("id")

    if err := c.BindJSON(&transaction_details); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := models.DB.Model(&models.Transaction{}).Where("id = ?", id).Updates(&transaction_details).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update transaction_details"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"transaction_details": transaction_details})
}
