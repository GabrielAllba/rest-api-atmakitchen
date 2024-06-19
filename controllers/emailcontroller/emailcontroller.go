package emailcontroller

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

const (
	CONFIG_SMTP_HOST     = "smtp.gmail.com"
	CONFIG_SMTP_PORT     = 587
	CONFIG_SENDER_NAME   = "gabrielallbasy@gmail.com"
	CONFIG_AUTH_EMAIL    = "gabrielallbasy@gmail.com"
	CONFIG_AUTH_PASSWORD = "titu yeti vmrw uwcu"
)

type ResetPasswordRequest struct {
	Email string `json:"email"`
}

func generateResetToken() (string, error) {
	token := make([]byte, 16)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

func SendResetPasswordEmail(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	emailReceiver := req.Email
	resetToken, err := generateResetToken()
	if err != nil {
		log.Println("Failed to generate reset token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset token"})
		return
	}

	resetURL := "http://localhost:3000/reset_password/change_password?token=" + resetToken

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", emailReceiver)
	mailer.SetHeader("Subject", "Reset Password Anda")
	mailer.SetBody("text/html", "Hello, <b>Silakan klik link berikut untuk mereset password Anda: </b><a href='"+resetURL+"'>"+resetURL+"</a>")

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	if err := dialer.DialAndSend(mailer); err != nil {
		log.Println("Failed to send email:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	// Save the reset token to the database with an expiry time
	// This part depends on your database setup and is not shown here

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}
