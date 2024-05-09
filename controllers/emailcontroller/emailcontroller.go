package emailcontroller

import (
	//email dgn library go get -u gopkg.in/gomail.v2
	"log"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
const CONFIG_SENDER_NAME = "PT. Atma Kitchen Gokil <arifkurniawanharrisma@gmail.com>"
const CONFIG_AUTH_EMAIL = "arifkurniawanharrisma@gmail.com"
const CONFIG_AUTH_PASSWORD = "h4rrism4bro."

func SendEmail(c *gin.Context) {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", "andreasmargono.23@gmail.com")
	mailer.SetHeader("Subject", "Selamat Anda Menang!")
	mailer.SetBody("text/html", "Hello, <b>have a nice day</b>")
	// mailer.Attach("./sample.png")

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Mail sent!")
}
