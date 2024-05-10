package rolecontroller

import (
	"backend-atmakitchen/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)



func Index(c *gin.Context) {
	var role []models.Role
	models.DB.Find(&role)
	c.JSON(http.StatusOK, gin.H{"role": role})
}

func Create(c *gin.Context) {
	var role models.Role

	
	// Bind JSON data to user struct
	if err := c.BindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	
	if role.Name == "" || strconv.Itoa(int(role.GajiHarian)) == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Pastikan semua input terisi"})
		return
	}



	var existingRole models.Role

	if err := models.DB.Where("name = ?", role.Name).First(&existingRole).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role dengan nama tersebut sudah ada"})
		return
	}

	
	if err := models.DB.Create(&role).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": role})
}

func Update(c *gin.Context) {
	var role models.Role

	// Get role ID from the URL parameter
	id := c.Param("id")

	// Check if the role exists
	if err := models.DB.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	// Bind JSON data to role struct
	if err := c.BindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the role in the database
	if err := models.DB.Save(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role updated successfully", "role": role})
}



func Show(c *gin.Context) {
	var role models.Role
	id := c.Param("id")
	if err := models.DB.First(&role, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Role not found"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"role": role})
}

func Delete(c *gin.Context) {
    // Extract the product ID from the request parameters
    id := c.Param("id")

    // Convert the ID string to an integer
    roleID, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
        return
    }

    // Find the product in the database
    var role models.Role
    if err := models.DB.First(&role, roleID).Error; err != nil {
        switch err {
        case gorm.ErrRecordNotFound:
            c.JSON(http.StatusNotFound, gin.H{"message": "Role not found"})
            return
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
            return
        }
    }

    // Delete the product from the database
    if err := models.DB.Delete(&role).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role"})
        return
    }

    // Respond with a success message
    c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}
