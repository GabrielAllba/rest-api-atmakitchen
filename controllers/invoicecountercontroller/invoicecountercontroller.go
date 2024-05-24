package invoicecountercontroller

import (
	"fmt"
	"net/http"
	"time"

	"backend-atmakitchen/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Create(c *gin.Context) {
    now := time.Now()
    year := now.Year()
    month := int(now.Month())

    var invoiceCounter models.InvoiceCounter

    if err := models.DB.First(&invoiceCounter).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            invoiceCounter.LastInvoice = 0
            models.DB.Create(&invoiceCounter)
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch invoice counter"})
            return
        }
    }
    newInvoiceNumber := invoiceCounter.LastInvoice + 1
    invoiceCounter.LastInvoice = newInvoiceNumber

    if err := models.DB.Save(&invoiceCounter).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update invoice counter"})
        return
    }

    // Format the invoice number
    formattedInvoiceNumber := fmt.Sprintf("%d.%02d.%d", year, month, newInvoiceNumber)


    c.JSON(http.StatusOK, gin.H{"invoice_number": formattedInvoiceNumber})
}
