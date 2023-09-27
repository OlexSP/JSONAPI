package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

const (
	tokenTTL            = 12 * time.Hour
	authorizationHeader = "Authorization"
	loginURL            = "/login"
	accountURL          = "/account"
	accountIdURL        = "/account/{id}"
	transferURL         = "/transfer"
)

type APIServer struct {
	listenAddr string
	storage    Storage
}

func NewAPIServer(listenAddr string, storage Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		storage:    storage,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc(loginURL, errorMiddleware(s.handleLogin)).Methods(http.MethodPost)
	router.HandleFunc(accountURL, errorMiddleware(s.handleGetAccount)).Methods(http.MethodGet)
	router.HandleFunc(accountURL, errorMiddleware(s.handleCreateAccount)).Methods(http.MethodPost)
	router.HandleFunc(accountIdURL, withJWTAuthMiddleware(errorMiddleware(s.handleGetAccountByID), s.storage)).Methods(http.MethodGet)
	router.HandleFunc(accountIdURL, errorMiddleware(s.handleDeleteAccount)).Methods(http.MethodDelete)
	router.HandleFunc(transferURL, errorMiddleware(s.handleTransfer)).Methods(http.MethodPost)

	slog.Info("API server listening on", slog.String("address", s.listenAddr))

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	account, err := s.storage.GetAccountByNumber(req.Number)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Encrypted), []byte(req.Password)); err != nil {
		return err
	}

	tokenString, err := createJWT(account)
	if err != nil {
		return err
	}

	response := loginResponse{
		Number: account.Number,
		Token:  tokenString,
	}
	return WriteJSON(w, http.StatusOK, response)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	account, err := NewAccount(req.FirstName, req.LastName, req.Password)
	if err != nil {
		return err
	}

	if err := s.storage.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, account.UUID.String())
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.storage.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	account, err := s.storage.GetAccountByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.storage.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]string{"deleted ID": id.String()})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferRequest := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferRequest); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, transferRequest)
}

func WriteJSON(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func errorMiddleware(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle error
			slog.Error("Error", slog.String("error", err.Error()))
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func getID(r *http.Request) (uuid.UUID, error) {
	idStr := mux.Vars(r)["id"]

	id, err := uuid.FromString(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid ID %s", idStr)
	}

	return id, nil
}

func createJWT(account *Account) (string, error) {
	secret := os.Getenv("JWT_SECRET")

	claims := &jwt.MapClaims{
		"ExpiresAt":     time.Now().Add(tokenTTL).Unix(),
		"AccountNumber": account.Number,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "permission denied"})
}

func withJWTAuthMiddleware(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("withJWTAuthMiddleware")

		tokenString := r.Header.Get(authorizationHeader)
		slog.Info(tokenString)
		token, err := validateJWT(tokenString)
		if err != nil {
			permissionDenied(w)
			slog.Info("withJWTAuthMiddleware validation", slog.String("error", err.Error()))
			return
		}
		if !token.Valid {
			permissionDenied(w)
			return
		}

		userID, err := getID(r)
		if err != nil {
			permissionDenied(w)
			return
		}

		account, err := s.GetAccountByID(userID)
		if err != nil {
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		slog.Info(fmt.Sprintf("clames type: %T, claims: -%[1]v-", claims["AccountNumber"]))
		slog.Info(fmt.Sprintf("accNumber type: %T, accNumber: -%[1]v-", account.Number))

		if account.Number != claims["AccountNumber"] {
			permissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// algorithm validation
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}
