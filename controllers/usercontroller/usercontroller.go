package usercontroller

import (
	"backend-atmakitchen/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Get balance
func GetBalance(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := models.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": user.Balance})
}

type WithdrawRequest struct {
	Amount    float64 `json:"amount"`
	BankName  string  `json:"bank_name"`
	AccountNo string  `json:"account_no"`
}

func WithdrawBalance(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := models.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var request WithdrawRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Amount <= 0 || request.Amount > user.Balance {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	// Optional: Validate bankName and accountNo if required

	user.Balance -= request.Amount
	if err := models.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user balance"})
		return
	}

	// Create a withdrawal history with status "Pending Approval"
	withdrawHistory := models.WithdrawHistory{
		UserId:    user.Id,
		Amount:    request.Amount,
		BankName:  request.BankName,
		AccountNo: request.AccountNo,
		Status:    "Pending Approval",
		CreatedAt: time.Now(),
	}
	if err := models.DB.Create(&withdrawHistory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create withdrawal history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Withdrawal request submitted for approval"})
}

// Get withdraw history
func GetWithdrawHistory(c *gin.Context) {
	id := c.Param("id")
	var histories []models.WithdrawHistory
	if err := models.DB.Where("user_id = ?", id).Find(&histories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, histories)
}

func UpdateWithdrawStatus(c *gin.Context) {
	id := c.Param("id")
	status := c.Param("status")

	// Temukan histori penarikan berdasarkan ID
	var history models.WithdrawHistory
	if err := models.DB.First(&history, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Withdraw history not found"})
		return
	}

	// Pastikan histori penarikan berada dalam status "Pending Approval"
	if history.Status != "Pending Approval" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Withdraw history is not pending approval"})
		return
	}

	// Ambil data pengguna berdasarkan UserId di histori penarikan
	var user models.User
	if err := models.DB.First(&user, history.UserId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user details"})
		return
	}

	// Periksa status yang diminta dan lakukan tindakan yang sesuai
	switch status {
	case "approved":
		history.Status = "Approved"
	case "rejected":
		history.Status = "Rejected"
		// Kembalikan jumlah saldo pengguna
		user.Balance += history.Amount
		if err := models.DB.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user balance"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	// Simpan perubahan status histori penarikan
	if err := models.DB.Save(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update withdraw status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Withdrawal status updated"})
}

func GetAllWithdrawHistory(c *gin.Context) {
	var histories []models.WithdrawHistory
	if err := models.DB.Find(&histories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, histories)
}

func Signup(c *gin.Context) {
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
	if err != nil {
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
	if err := models.DB.Where("id = ?", roleId).First(&role).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role tersebut tidak tersedia"})
		return
	}
	user.RoleId = roleId

	// password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
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

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"user": "Gagal membuat user"})
		return
	}

	var returnUser models.User
	models.DB.Preload("Role").First(&returnUser, "id = ?", user.Id)

	c.JSON(http.StatusOK, gin.H{"User": returnUser})
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

	if user.Role.Name != "Customer" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User role bukan Customer",
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"subject": user.Id,
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})

		return
	}

	// send it back

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

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

func Search(c *gin.Context) {
	query := c.Query("query")
	var user []models.User

	query = strings.ToLower(query)
	result := models.DB.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ? OR LOWER(username) LIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%")

	if err := result.Find(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": "No user found"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func Delete(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := models.DB.First(&user, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "user not found"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
	}
	models.DB.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func UpdateUser(c *gin.Context) {
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

	// Update the user in the database
	if err := models.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": user})
}

func Show(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	if err := models.DB.First(&user, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "user not found"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
	}

	// Successfully found the user
	c.JSON(http.StatusOK, gin.H{"user": user})
}
func UpdatePoints(c *gin.Context) {
	var user models.User
	id := c.Param("id")
	pointsStr := c.Param("points")

	// Convert points query parameter to an integer
	points, err := strconv.Atoi(pointsStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid points parameter"})
		return
	}

	// Retrieve the user by ID
	if err := models.DB.First(&user, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
	}

	// Update the user's points
	user.TotalPoint = points

	// Save the updated user record to the database
	if err := models.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update points"})
		return
	}

	// Return the updated user
	c.JSON(http.StatusOK, gin.H{"message": "Points updated successfully", "user": user})
}

func SearchType(c *gin.Context) {
	queryStr := c.Query("query")
	query, err := strconv.Atoi(queryStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query format"})
		return
	}
	var users []models.User

	result := models.DB.Where("role_id = ?", query)

	if err := result.Find(&users).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": "No customers found"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func SearhUserByType(c *gin.Context) {
	searchQuery := c.Query("search_query")
	queryStr := c.Query("query")
	query, err := strconv.Atoi(queryStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query format"})
		return
	}
	var users []models.User

	searchQuery = strings.ToLower(searchQuery)

	result := models.DB.Where("(LOWER(name) LIKE ? OR LOWER(email) LIKE ? OR phone_number LIKE ?) AND role_id = ?", "%"+searchQuery+"%", "%"+searchQuery+"%", "%"+searchQuery+"%", query)

	if err := result.Find(&users).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": "No customers found"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}
