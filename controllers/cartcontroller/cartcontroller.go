package cartcontroller

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
	var carts []models.Cart
	if err := models.DB.Where("user_id = ?", userID).Preload("Product").Preload("Hampers").Find(&carts).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Carts not found for this user"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch carts for this user"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"carts": carts})
}
    
func Create(c *gin.Context) {
	var cart models.Cart

	// Bind the incoming JSON payload to the cart struct
	if err := c.BindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the cart entry in the database
	if err := models.DB.Create(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create cart"})
		return
	}

	// Retrieve the created cart entry with the associated product details preloaded
	if err := models.DB.Preload("Product").Preload("Hampers").First(&cart, cart.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to preload product"})
		return
	}

	// Return the cart with the preloaded product details
	c.JSON(http.StatusOK, gin.H{"cart": cart})
}

func Show(c *gin.Context) {
    id := c.Param("id")
    var cart models.Cart

    if err := models.DB.Preload("Product").Preload("Hampers").First(&cart, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{"cart": cart})
}

func Index(c *gin.Context) {
    var carts []models.Cart

    if err := models.DB.Preload("Product").Preload("Hampers").Find(&carts).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"carts": carts})
}

func Delete(c *gin.Context) {
    id := c.Param("id")
    intId, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    var cart models.Cart
    if err := models.DB.First(&cart, intId).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
        return
    }

    if err := models.DB.Delete(&cart).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cart"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Cart deleted successfully"})
}

func Update(c *gin.Context) {
    var cart models.Cart
    id := c.Param("id")

    if err := c.BindJSON(&cart); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := models.DB.Model(&models.Cart{}).Where("id = ?", id).Updates(&cart).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update cart"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"cart": cart})
}
