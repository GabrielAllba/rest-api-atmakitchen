package transactiondetailcontroller

import (
	"net/http"
	"strconv"

	"backend-atmakitchen/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)



func GetAllBahan(c *gin.Context) {
	transactionDetailIDStr := c.Param("id")

	transactionDetailID, err := strconv.Atoi(transactionDetailIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction detail ID"})
		return
	}

	var transactionDetail models.TransactionDetail
	if err := models.DB.Preload("Product").First(&transactionDetail, transactionDetailID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "transaction detail not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction detail"})
		}
		return
	}

	if transactionDetail.ProductId == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction detail does not have a product ID"})
		return
	}

	var bahanResep []models.BahanResep
	if err := models.DB.Preload("Bahan").Where("product_id = ?", *transactionDetail.ProductId).Find(&bahanResep).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bahan for product"})
		return
	}

	var allBahan []models.Bahan
	for _, br := range bahanResep {
		allBahan = append(allBahan, br.Bahan)
	}

	c.JSON(http.StatusOK, gin.H{"bahan": allBahan})
}

func GetAllBahanByProductID(c *gin.Context) {
	productIDStr := c.Param("product_id")

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := models.DB.First(&product, productID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		}
		return
	}

	var bahanResep []models.BahanResep
	if err := models.DB.Preload("Bahan").Where("product_id = ?", productID).Find(&bahanResep).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bahan for product"})
		return
	}

	var allBahan []models.Bahan
	for _, br := range bahanResep {
		allBahan = append(allBahan, br.Bahan)
	}

	c.JSON(http.StatusOK, gin.H{"bahan": allBahan})
}


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

    if err := models.DB.Preload("Product").Preload("Hampers").Find(&transaction_details).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"transaction_details": transaction_details})
}

func ProsesToday(c *gin.Context) {
    var transaction_details []models.TransactionDetail

    date := c.Param("tanggal_pengiriman")
    if date == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "tanggal_pengiriman is required"})
        return
    }

    // Preload Product and Hamper and filter by the given date
    if err := models.DB.Preload("Product").Preload("Hampers").Where("DATE(tanggal_pengiriman) = ?", date).Find(&transaction_details).Error; err != nil {
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

func GetByUserID(c *gin.Context) {
    userID := c.Param("userId")

    // Retrieve transactions for the given user ID
    var transactions []models.Transaction
    if err := models.DB.Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "transactions not found for this user"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions for this user"})
        }
        return
    }

    // Retrieve transaction details for each transaction
    var allTransactionDetails []models.TransactionDetail
    for _, transaction := range transactions {
        var transactionDetails []models.TransactionDetail
        if err := models.DB.Preload("Product").Preload("Hampers").Where("invoice_number = ?", transaction.InvoiceNumber).Find(&transactionDetails).Error; err != nil {
            if err == gorm.ErrRecordNotFound {
                c.JSON(http.StatusNotFound, gin.H{"error": "transaction details not found for invoice number"})
                return
            } else {
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction details for invoice number"})
                return
            }
        }
        allTransactionDetails = append(allTransactionDetails, transactionDetails...)
    }

    c.JSON(http.StatusOK, gin.H{"transaction_details": allTransactionDetails})
}

func GetPhotosByInvoiceNumber(c *gin.Context) {
	invoiceNumber := c.Param("invoiceNumber")

	var transactionDetails []models.TransactionDetail
	if err := models.DB.Preload("Product").Preload("Hampers").Where("invoice_number = ?", invoiceNumber).Find(&transactionDetails).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "transaction details not found for this invoice number"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction details for this invoice number"})
		}
		return
	}

	var photos []string
	for _, detail := range transactionDetails {
		if detail.Product != nil && detail.Product.Photo != "" {
			photos = append(photos, detail.Product.Photo)
		}
		if detail.Hampers != nil && detail.Hampers.Photo != "" {
			photos = append(photos, detail.Hampers.Photo)
		}
	}

	c.JSON(http.StatusOK, gin.H{"photos": photos})
}

func GetPhotosByUserID(c *gin.Context) {
	userId := c.Param("userId")

	var transactions []models.Transaction
	if err := models.DB.Where("user_id = ?", userId).Find(&transactions).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "transactions not found for this user"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions for this user"})
		}
		return
	}

	var photos []string
	for _, transaction := range transactions {
		var transactionDetails []models.TransactionDetail
		if err := models.DB.Preload("Product").Preload("Hampers").Where("invoice_number = ?", transaction.InvoiceNumber).Find(&transactionDetails).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "transaction details not found for invoice number"})
				return
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction details for invoice number"})
				return
			}
		}

		for _, detail := range transactionDetails {
			if detail.Product != nil && detail.Product.Photo != "" {
				photos = append(photos, detail.Product.Photo)
			}
			if detail.Hampers != nil && detail.Hampers.Photo != "" {
				photos = append(photos, detail.Hampers.Photo)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"photos": photos})
}