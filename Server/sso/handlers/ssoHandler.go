package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sso/data"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type SSOHandler struct {
	repo *data.SSORepo
}

var secretKey = []byte("eUpravaT2")

const InvalidRequestBody = "Invalid request body"

// Constructor
func NewSSOHandler(r *data.SSORepo) *SSOHandler {
	return &SSOHandler{r}
}

// Handler methods

// Logins requested user and provides JWT token
func (sh *SSOHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials data.Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, InvalidRequestBody, http.StatusBadRequest)
		log.Println("Error while decoding body")
		return
	}

	log.Printf("Recieved login request from '%s' for user '%s'", r.RemoteAddr, credentials.Email)

	dbUser, err := sh.repo.FindAccountByEmail(credentials.Email)
	if err != nil {
		http.Error(w, "User not found with email "+credentials.Email, http.StatusBadRequest)
		log.Printf("User '%s' does not exist", credentials.Email)
		return
	}

	if err := sh.validateCredentials(credentials.Email, credentials.Password); err != nil {
		if err.Error() == "account not activated" {
			http.Error(w, "Account not activated", http.StatusForbidden)
			log.Printf("User '%s' account is not activated", credentials.Email)
			return
		}
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		log.Printf("Failed to log in '%s'", credentials.Email)
		return
	}

	token, err := sh.generateToken(credentials.Email, dbUser.Role)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		log.Printf("Failed to generate token for '%s'", credentials.Email)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})

	log.Printf("User '%s' successfully logged in from '%s'", credentials.Email, r.RemoteAddr)
}

func (sh *SSOHandler) validateCredentials(email, password string) error {
	account, err := sh.repo.FindAccountByEmail(email)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		return errors.New("invalid password")
	}

	if !account.Activated {
		return errors.New("account not activated")
	}

	return nil
}

func (sh *SSOHandler) generateToken(email, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"sub":  email,
		"role": role,
		"exp":  expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
