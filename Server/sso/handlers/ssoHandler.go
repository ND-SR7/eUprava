package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sso/data"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
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

// Registers a new person to the system
func (sh *SSOHandler) RegisterPerson(w http.ResponseWriter, r *http.Request) {
	log.Println("Registering new user")

	var newPerson data.NewPerson
	if err := json.NewDecoder(r.Body).Decode(&newPerson); err != nil {
		http.Error(w, InvalidRequestBody, http.StatusBadRequest)
		log.Println("Error while decoding body")
		return
	}

	err := sh.repo.CreatePerson(newPerson)
	if err != nil && err.Error() == "email already taken" {
		http.Error(w, "Email is already in use by an account", http.StatusBadRequest)
		log.Printf("Failed to register new user: email '%s' already in use", newPerson.Email)
		return
	} else if err != nil {
		http.Error(w, "Failed to register new user", http.StatusInternalServerError)
		log.Printf("Failed to register new user: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Println("Successfully registered a new user")
}

// Registers a new legal entity to the system
func (sh *SSOHandler) RegisterLegalEntity(w http.ResponseWriter, r *http.Request) {
	log.Println("Registering new legal entity")

	var newLegalEntity data.NewLegalEntity
	if err := json.NewDecoder(r.Body).Decode(&newLegalEntity); err != nil {
		http.Error(w, InvalidRequestBody, http.StatusBadRequest)
		log.Println("Error while decoding body")
		return
	}

	err := sh.repo.CreateLegalEntity(newLegalEntity)
	if err != nil && err.Error() == "email already taken" {
		http.Error(w, "Email is already in use by an account", http.StatusBadRequest)
		log.Printf("Failed to register new user: email '%s' already in use", newLegalEntity.Email)
		return
	} else if err != nil {
		http.Error(w, "Failed to register new user", http.StatusInternalServerError)
		log.Printf("Failed to register new user: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Println("Successfully registered a new user")
}

// Activates user account after email confirmation
func (sh *SSOHandler) ActivateAccount(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	activationCode := params["activationCode"]

	log.Printf("Activating user account with code '%s'", activationCode)

	err := sh.repo.ActivateAccount(activationCode)
	if err != nil {
		http.Error(w, "Failed to activate user account", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User account successfully activated"))
	log.Printf("Successfully activated user account with code '%s'", activationCode)
}

// Returns error if credentials are not valid
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

// Generates token for logged in user
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

// Returns token string found in header, otherwise empty string
func (sh *SSOHandler) extractTokenFromHeader(r *http.Request) string {
	token := r.Header.Get("Authorization")
	if token != "" {
		return token[len("Bearer "):]
	}
	return ""
}