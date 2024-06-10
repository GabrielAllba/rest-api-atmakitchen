package transactioncontroller

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"backend-atmakitchen/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func UpdateStatus(c *gin.Context) {
	var transaction models.Transaction
	var transaction_detail models.TransactionDetail

	id := c.Param("id")
	transaction_status := c.Param("transaction_status")

	// Fetch the transaction to get the InvoiceNumber
	if err := models.DB.Where("id = ?", id).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Transaction not found"})
		return
	}

	// Update the status of the transaction
	if err := models.DB.Model(&transaction).Where("id = ?", id).Update("transaction_status", transaction_status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update transaction status"})
		return
	}

	// Update the status of the related transaction detail
	if err := models.DB.Model(&transaction_detail).Where("invoice_number = ?", transaction.InvoiceNumber).Update("transaction_status", transaction_status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update transaction detail status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction": transaction})
}

func UpdateStatusByInvoice(c *gin.Context) {
	var transaction models.Transaction
	var transaction_detail models.TransactionDetail

	invoice_number := c.Param("invoice_number")
	transaction_status := c.Param("transaction_status")

	// Fetch the transaction to get the InvoiceNumber
	if err := models.DB.Where("invoice_number = ?", invoice_number).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Transaction not found"})
		return
	}

	// Update the status of the transaction
	if err := models.DB.Model(&transaction).Where("invoice_number = ?", invoice_number).Update("transaction_status", transaction_status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update transaction status"})
		return
	}

	// Update the status of the related transaction detail
	if err := models.DB.Model(&transaction_detail).Where("invoice_number = ?", transaction.InvoiceNumber).Update("transaction_status", transaction_status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update transaction detail status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction": transaction})
}

func UpdateTotalAfterDeliveryFee(c *gin.Context) {
	var transaction models.Transaction

	id := c.Param("id")

	// Fetch the transaction to get the current values
	if err := models.DB.Where("id = ?", id).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Transaction not found"})
		return
	}

	// Calculate the new transfer_nominal
	newTransferNominal := transaction.TransferNominal + transaction.DeliveryFee

	// Update the transfer_nominal of the transaction
	if err := models.DB.Model(&transaction).Where("id = ?", id).Update("transfer_nominal", newTransferNominal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update transaction total"})
		return
	}

	// Update the transaction object to reflect the new value
	transaction.TransferNominal = newTransferNominal

	c.JSON(http.StatusOK, gin.H{"transaction": transaction})
}

func UpdateBuktiPembayaran(c *gin.Context) {
	// Get the transaction ID from the URL parameters
	invoice_number := c.Param("invoice_number")

	// Handle file upload for photo
	file, fileHeader, err := c.Request.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Generate a unique filename for the uploaded file
	filename := uuid.New().String() + filepath.Ext(fileHeader.Filename)
	filePath := filepath.Join("images", filename) // Assuming "images" folder exists in your project directory
	fullFilePath := filepath.Join("/", filePath)  // Full file path including the leading "/"

	// Save the uploaded file
	if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	fixFullFilePath := strings.Replace(fullFilePath, "\\", "/", -1)

	// Fetch the transaction to ensure it exists
	var transaction models.Transaction
	if err := models.DB.Where("invoice_number = ?", invoice_number).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Transaction not found"})
		return
	}

	// Update the payment_proof column
	if err := models.DB.Model(&transaction).Where("invoice_number = ?", invoice_number).Update("payment_proof", fixFullFilePath).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update payment proof"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction": transaction, "payment_proof": fixFullFilePath})
}

func RejectTransaction(c *gin.Context) {
	// Retrieve user ID from the URL parameters
	userID := c.Param("userId")

	// Parse the user ID to an integer
	id, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Find the transaction by user ID and status
	var transaction models.Transaction
	if err := models.DB.Where("user_id = ? AND transaction_status = ?", id, "Pending").First(&transaction).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found or not pending"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Find all transaction details by invoice number
	var transactionDetails []models.TransactionDetail
	if err := models.DB.Where("invoice_number = ?", transaction.InvoiceNumber).Find(&transactionDetails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Iterate over each transaction detail and return the stock
	for _, detail := range transactionDetails {
		if detail.ProductId != nil && detail.ProductQuantity != nil {
			var product models.Product
			if err := models.DB.First(&product, *detail.ProductId).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Return the stock
			product.Stock += *detail.ProductQuantity
			if err := models.DB.Save(&product).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	// Update transaction status to "Ditolak"
	transaction.TransactionStatus = "Ditolak"
	if err := models.DB.Save(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction rejected and stock returned successfully"})
}
