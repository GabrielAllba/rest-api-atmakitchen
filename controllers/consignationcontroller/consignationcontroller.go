package consignationcontroller

import (
	"backend-atmakitchen/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Index(c *gin.Context) {
	var consignation []models.Consignation
	models.DB.Find(&consignation)
	c.JSON(http.StatusOK, gin.H{"consignation": consignation})
}

func Create(c *gin.Context) {
	var consignation models.Consignation

	
	if err := c.BindJSON(&consignation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	
	if consignation.Name == "" || consignation.Address == "" || consignation.PhoneNumber == "" || consignation.BankAccount == "" || consignation.BankNumber == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Pastikan semua input terisi"})
		return
	}

	var existingConsignation models.Consignation

	if err := models.DB.Where("name = ?", consignation.Name).First(&existingConsignation).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Consignation dengan nama tersebut sudah ada"})
		return
	}

	
	if err := models.DB.Create(&consignation).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat consignation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"consignation": consignation})
}



func Search(c *gin.Context) {
	query := c.Query("query")
	var consignation []models.Consignation

	
	query = strings.ToLower(query)
	result := models.DB.Where("LOWER(name) LIKE ? OR LOWER(address) LIKE ? OR LOWER(phone_number) LIKE ? OR LOWER(bank_account) LIKE ? OR LOWER(bank_number) LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%")

	
	if err := result.Find(&consignation).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": "No consignation found"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"consignation": consignation})
}

func Delete(c *gin.Context) {
    
    id := c.Param("id")

    consignationID, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Consignation ID"})
        return
    }

    var consignation models.Consignation
    if err := models.DB.First(&consignation, consignationID).Error; err != nil {
        switch err {
        case gorm.ErrRecordNotFound:
            c.JSON(http.StatusNotFound, gin.H{"message": "consignation not found"})
            return
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
            return
        }
    }



    // Delete the consignation from the database
    if err := models.DB.Delete(&consignation).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete consignation"})
        return
    }

    // Respond with a success message
    c.JSON(http.StatusOK, gin.H{"message": "consignation deleted successfully"})
}

func Update(c *gin.Context) {
    // Get the Consignation ID from the URL parameter
    id := c.Param("id")

    // Convert the Consignation ID to an integer
    consignationID, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Consignation ID"})
        return
    }

    // Fetch the existing Consignation from the database
    var consignation models.Consignation
    if err := models.DB.First(&consignation, consignationID).Error; err != nil {
        switch err {
        case gorm.ErrRecordNotFound:
            c.JSON(http.StatusNotFound, gin.H{"message": "Consignation not found"})
            return
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
            return
        }
    }

    // Bind the updated Consignation data from the request body
    var updatedConsignation models.Consignation
    if err := c.BindJSON(&updatedConsignation); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Update the fields of the existing Consignation with the new values
    consignation.Name = updatedConsignation.Name
    consignation.Address = updatedConsignation.Address
    consignation.PhoneNumber = updatedConsignation.PhoneNumber
    consignation.BankAccount = updatedConsignation.BankAccount
    consignation.BankNumber = updatedConsignation.BankNumber

    // Save the updated Consignation back to the database
    if err := models.DB.Save(&consignation).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update Consignation"})
        return
    }

    // Respond with the updated Consignation data
    c.JSON(http.StatusOK, gin.H{"consignation": consignation})
}
