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

// Retrieves user based on provided ID
func (sh *SSOHandler) GetUserByAccountID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	accountID := params["accountID"]

	log.Printf("Retrieving user with id '%s'", accountID)

	person, err := sh.getPersonByID(accountID)
	if err != nil && err.Error() != "person not found" {
		http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		log.Printf("Failed to retrieve user: %s", err.Error())
		return
	} else if err.Error() == "person not found" {
		legalEntity, err := sh.getLegalEntityByID(accountID)
		if err != nil {
			http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
			log.Printf("Failed to retrieve user: %s", err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(legalEntity); err != nil {
			http.Error(w, "Error while encoding body", http.StatusInternalServerError)
			log.Printf("Error while encoding legal entity: %s", err.Error())
		}

		log.Println("Successfully retrieved requested user")
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(person); err != nil {
			http.Error(w, "Error while encoding body", http.StatusInternalServerError)
			log.Printf("Error while encoding person: %s", err.Error())
		}

		log.Println("Successfully retrieved requested user")
	}
}

// Retrieves user based on provided ID
func (sh *SSOHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	email := params["email"]

	log.Printf("Retrieving user with email '%s'", email)

	person, err := sh.getPersonByEmail(email)
	if err != nil && err.Error() != "person not found" {
		http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		log.Printf("Failed to retrieve user: %s", err.Error())
		return
	} else if err != nil && err.Error() == "person not found" {
		legalEntity, err := sh.getLegalEntityByEmail(email)
		if err != nil {
			http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
			log.Printf("Failed to retrieve user: %s", err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(legalEntity); err != nil {
			http.Error(w, "Error while encoding body", http.StatusInternalServerError)
			log.Printf("Error while encoding legal entity: %s", err.Error())
		}

		log.Println("Successfully retrieved requested user")
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(person); err != nil {
			http.Error(w, "Error while encoding body", http.StatusInternalServerError)
			log.Printf("Error while encoding person: %s", err.Error())
		}

		log.Println("Successfully retrieved requested user")
	}
}

// Retrieves person based on provided JMBG
func (sh *SSOHandler) GetPersonByJMBG(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	jmbg := params["jmbg"]

	log.Printf("Retrieving user with JMBG: %s", jmbg)

	person, err := sh.getPersonByJMBG(jmbg)
	if err != nil && err.Error() != "person not found" {
		http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		log.Printf("Failed to retrieve user: %s", err.Error())
		return
	} else if err != nil && err.Error() == "person not found" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Person not found for requested JMBG"))
		log.Printf("Person not found for JMBG: %s", jmbg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(person); err != nil {
		http.Error(w, "Error while encoding body", http.StatusInternalServerError)
		log.Printf("Error while encoding legal entity: %s", err.Error())
	}

	log.Println("Successfully retrieved requested user")
}

// Retrieves legal entity based on provided MB
func (sh *SSOHandler) GetLegalEntityByMB(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	mb := params["mb"]

	log.Printf("Retrieving user with MB: %s", mb)

	legalEntity, err := sh.getLegalEntityByMB(mb)
	if err != nil && err.Error() != "legal entity not found" {
		http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		log.Printf("Failed to retrieve user: %s", err.Error())
		return
	} else if err != nil && err.Error() == "legal entity not found" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Legal entity not found for requested MB"))
		log.Printf("Legal entity not found for JMBG: %s", mb)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(legalEntity); err != nil {
		http.Error(w, "Error while encoding body", http.StatusInternalServerError)
		log.Printf("Error while encoding legal entity: %s", err.Error())
	}

	log.Println("Successfully retrieved requested user")
}

// Logins requested user and provides JWT token
func (sh *SSOHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials data.Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, InvalidRequestBody, http.StatusBadRequest)
		log.Println("Error while decoding body")
		return
	}

	log.Printf("Recieved login request from '%s' for user '%s'", r.RemoteAddr, credentials.Email)

	account, err := sh.repo.FindAccountByEmail(credentials.Email)
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

	person, err := sh.getPersonByEmail(credentials.Email)
	if err != nil && err.Error() != "person not found" {
		http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		log.Printf("Failed to retrieve user: %s", err.Error())
		return
	} else if err != nil && err.Error() == "person not found" {
		legalEntity, err := sh.getLegalEntityByEmail(credentials.Email)
		if err != nil {
			http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
			log.Printf("Failed to retrieve user: %s", err.Error())
			return
		}

		token, err := sh.generateToken(legalEntity.MB, account.Role)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			log.Printf("Failed to generate token for '%s'", credentials.Email)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})

		log.Printf("User '%s' successfully logged in from '%s'", credentials.Email, r.RemoteAddr)
	} else {
		token, err := sh.generateToken(person.JMBG, account.Role)
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

// Helper function for retrieving person based on AccountID
func (sh *SSOHandler) getPersonByID(accountID string) (data.Person, error) {
	person, err := sh.repo.GetPersonByID(accountID)
	if err != nil {
		return data.Person{}, err
	}

	return person, nil
}

// Helper function for retrieving legal entity based on AccountID
func (sh *SSOHandler) getLegalEntityByID(accountID string) (data.LegalEntity, error) {
	legalEntity, err := sh.repo.GetLegalEntityByID(accountID)
	if err != nil {
		return data.LegalEntity{}, err
	}

	return legalEntity, nil
}

// Helper function for retrieving person based on email
func (sh *SSOHandler) getPersonByEmail(email string) (data.Person, error) {
	person, err := sh.repo.GetPersonByEmail(email)
	if err != nil {
		return data.Person{}, err
	}

	return person, nil
}

// Helper function for retrieving legal entity based on email
func (sh *SSOHandler) getLegalEntityByEmail(email string) (data.LegalEntity, error) {
	legalEntity, err := sh.repo.GetLegalEntityByEmail(email)
	if err != nil {
		return data.LegalEntity{}, err
	}

	return legalEntity, nil
}

// Helper function for retrieving person based on JMBG
func (sh *SSOHandler) getPersonByJMBG(jmbg string) (data.Person, error) {
	person, err := sh.repo.GetPersonByJMBG(jmbg)
	if err != nil {
		return data.Person{}, err
	}

	return person, nil
}

// Helper function for retrieving legal entity based on MB
func (sh *SSOHandler) getLegalEntityByMB(mb string) (data.LegalEntity, error) {
	legalEntity, err := sh.repo.GetLegalEntityByMB(mb)
	if err != nil {
		return data.LegalEntity{}, err
	}

	return legalEntity, nil
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
func (sh *SSOHandler) generateToken(jmbg, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"sub":  jmbg,
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

// JWT middleware
func (sh *SSOHandler) AuthorizeRoles(allowedRoles ...string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, rr *http.Request) {
			tokenString := sh.extractTokenFromHeader(rr)
			if tokenString == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return secretKey, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			_, ok1 := claims["sub"].(string)
			role, ok2 := claims["role"].(string)
			if !ok1 || !ok2 {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			for _, allowedRole := range allowedRoles {
				if allowedRole == role {
					next.ServeHTTP(w, rr)
					return
				}
			}

			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}

// Returns token string found in header, otherwise empty string
func (sh *SSOHandler) extractTokenFromHeader(r *http.Request) string {
	token := r.Header.Get("Authorization")
	if token != "" {
		return token[len("Bearer "):]
	}
	return ""
}
