package data

import (
	"log"
	"net/smtp"
	"os"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// Checks whether password fulfills requirements
func CheckPassword(password string) bool {
	if len(password) < 6 {
		return false
	}

	var (
		hasUpperCase bool
		hasLowerCase bool
		hasNumber    bool
	)

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpperCase = true
		}
		if unicode.IsLower(char) {
			hasLowerCase = true
		}
		if unicode.IsNumber(char) {
			hasNumber = true
		}
	}

	return hasUpperCase && hasLowerCase && hasNumber
}

// Generates BCrypt password hash for provided password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func SendEmail(email, code, intent string) bool {
	accountActivationPath := os.Getenv("ACCOUNT_ACTIVATION_PATH")
	accountRecoveryPath := "/todo"
	// TODO: Recovery path

	// Sender data
	from := os.Getenv("MAIL_ADDRESS")
	password := os.Getenv("MAIL_PASSWORD")

	// Receiver email
	to := []string{
		email,
	}

	// SMTP server config
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	address := smtpHost + ":" + smtpPort
	var subject string
	var body string

	if intent == "ACTIVATION" {
		subject = "eUprava Account Activation"
		body = "Follow the verification link to activate your eUprava account: \n" + accountActivationPath + code
	} else if intent == "RECOVERY" {
		subject = "eUprava Password Recovery"
		body = "To reset your password, copy the given code & then follow the recovery link: \n" + code + "\n" + accountRecoveryPath
	}
	// Text
	stringMsg :=
		"From: " + from + "\n" +
			"To: " + to[0] + "\n" +
			"Subject: " + subject + "\n\n" +
			body

	message := []byte(stringMsg)

	// Email Sender Auth
	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		log.Println("Error while sending email:", err)
		return false
	}
	log.Println("Email successfully sent")
	return true
}
