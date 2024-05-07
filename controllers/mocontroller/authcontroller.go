package moauthcontroller

import (
	"backend-atmakitchen/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)


func Login(c *gin.Context){
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

	if user.Id == 0{
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Email tidak ditemukan",
		})

		return
	}

	if user.Role.Name != "Manajer Operasional" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User role bukan Manajer Operasional",
		})
		return
    }

	err :=	bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req_user.Password))

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Invalid password",
		})

		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"subject": user.Id,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the SECRET_MO
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_MO")))

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error" : "Failed to create token",
		})

		return
	}

	// send it back

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600 * 24 * 30, "","", false, true)
		
	c.JSON(http.StatusOK, gin.H{
		"token":tokenString,
		"user":user,
	})
}
func Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

func Validate(c *gin.Context) {
    tokenString := c.Param("tokenString")
    log.Printf("Received token string: %s", tokenString)

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrSignatureInvalid
        }
        return []byte(os.Getenv("SECRET_MO")), nil
    })

    if err != nil {
        log.Printf("Error parsing token: %v", err)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    if !token.Valid {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
        return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to extract claims from token"})
        return
    }

    c.JSON(http.StatusOK, claims)
}