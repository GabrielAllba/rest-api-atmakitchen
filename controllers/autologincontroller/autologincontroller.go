package autologincontroller

import (
	"backend-atmakitchen/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func CekRole(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email harus diisi"})
		return
	}

	var user models.User // Gantilah dengan model user yang sesuai
	models.DB.Preload("Role").First(&user, "email = ?", email)

	if user.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": user.Role.Name})
}

func Login(c *gin.Context) {
	var req_user models.User

	// Bind other form data fields
	if err := c.BindJSON(&req_user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req_user.Email == "" || req_user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email dan Password harus diisi"})
		return
	}

	var user models.User
	models.DB.Preload("Role").First(&user, "email = ?", req_user.Email)

	if user.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email tidak ditemukan",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req_user.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid password",
		})
		return
	}

	var tokenString string
	var secret []byte
	if user.Role.Name == "Admin" {
		secret = []byte(os.Getenv("SECRET_ADMIN"))
	} else if user.Role.Name == "Customer" {
		secret = []byte(os.Getenv("SECRET"))
	} else if user.Role.Name == "Manajer Operasional" {
		secret = []byte(os.Getenv("SECRET_MO"))
	} else if user.Role.Name == "Owner" {
		secret = []byte(os.Getenv("SECRET_OWNER"))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"subject": user.Id,
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err = token.SignedString(secret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user":  user,
	})
}

func Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}
