package customerauthcontroller

import (
	"backend-atmakitchen/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context){
	var user models.User

	// Bind JSON data to user struct
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// validitas semua input terisi
	if user.Email == "" || user.Name == "" || user.Username == "" || user.Password == "" ||
		user.BornDate == "" || user.PhoneNumber == "" || strconv.Itoa(user.TotalPoint) == "" || strconv.Itoa(user.RoleId) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pastikan semua input terisi"})
		return
	}

	
	

	// total point	
	totalP, err := strconv.Atoi(strconv.Itoa(user.TotalPoint))
	if(err != nil){
		fmt.Print(totalP)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Total point invalid"})
		return
	}
	

	// role id
	roleId, err := strconv.Atoi(strconv.Itoa(user.RoleId))
	if err != nil {
		
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format role id tidak benar"})
		return
	}
	// check role id in database
	var role models.Role
	if err := models.DB.Where("id = ?", roleId).First(&role).Error; err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role tersebut tidak tersedia"})
		return
	}
	user.RoleId = roleId

	// password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error":"Failed to hash password",
		})
		return
	}

	// existing username and email
	var existingUser models.User
	if err := models.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username tersebut sudah ada"})
		return
	}
	if err := models.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email tersebut sudah ada"})
		return
	}


	
	
	user = models.User{Name: user.Name, Email: user.Email, Password: string(hash), Username: user.Username, BornDate: user.BornDate, PhoneNumber: user.PhoneNumber, TotalPoint: user.TotalPoint, RoleId: user.RoleId}
	result := models.DB.Create(&user)

	if result.Error != nil{
		c.JSON(http.StatusBadRequest, gin.H{"user": "Gagal membuat user"})
		return
	}

	var returnUser models.User
	models.DB.Preload("Role").First(&returnUser, "id = ?", user.Id)
	

	c.JSON(http.StatusOK, gin.H{"User": returnUser})
}


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

	if user.Role.Name != "Customer" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User role bukan Customer",
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

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

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
        return []byte(os.Getenv("SECRET")), nil
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

// func GetUser(c *gin.Context) {
//     // Get user ID from JWT token in the request header
//     userID, exists := c.Get("id")
//     if !exists {
//         c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
//         return
//     }

//     // Convert user ID to integer
//     userIDInt, ok := userID.(int)
//     if !ok {
//         c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
//         return
//     }

//     // Fetch user details from the database
//     var user models.User
//     if err := models.DB.Preload("Role").First(&user, userIDInt).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user details"})
//         return
//     }

//     // Return user details
//     c.JSON(http.StatusOK, user)
// }

// func GetUsersByRoleID(c *gin.Context) {
//     // Get role ID from the URL parameter
//     roleID, err := strconv.Atoi(c.Param("role_id"))
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
//         return
//     }

//     // Fetch users with the specified role ID from the database
//     var users []models.User
//     if err := models.DB.Where("role_id = ?", roleID).Find(&users).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
//         return
//     }

//     // Return the list of users with the specified role ID
//     c.JSON(http.StatusOK, users)
// }

// Function signature
func UpdatePassword(c *gin.Context) {
    var user models.User

    // Get email from the URL parameter
    email := c.Param("email")

    // Check if the user exists
    if err := models.DB.Where("email = ?", email).First(&user).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Bind JSON data to user struct
    if err := c.BindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Hash the password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    // Update the user's password with the hashed one
    user.Password = string(hashedPassword)

    // Update the user in the database
    if err := models.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully", "user": user})
}

func Index(c *gin.Context) {
	var user []models.User
	models.DB.Find(&user)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func EmailExists(c *gin.Context) {
	email := c.Query("email")

	var user models.User
	if err := models.DB.Where("email = ?", email).First(&user).Error; err == nil {
		// Email exists
		c.JSON(http.StatusOK, gin.H{"exists": true})
		return
	}

	// Email does not exist
	c.JSON(http.StatusOK, gin.H{"exists": false})
}

