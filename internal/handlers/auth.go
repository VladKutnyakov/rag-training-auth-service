package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"rag-training-auth-service/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *pgxpool.Pool
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req UserCredentials
	json.NewDecoder(r.Body).Decode(&req)

	var exists bool
	args := pgx.NamedArgs{
		"username": req.Username,
	}
	err := h.DB.QueryRow(r.Context(), "SELECT EXISTS(SELECT 1 FROM users WHERE username = @username)", args).Scan(&exists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Пользователь уже существует", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	args = pgx.NamedArgs{
		"username":      req.Username,
		"password_hash": hash,
	}
	_, err = h.DB.Exec(r.Context(), "INSERT INTO users (username, password_hash) VALUES (@username, @password_hash)", args)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь создан"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req UserCredentials
	json.NewDecoder(r.Body).Decode(&req)

	var hash string
	args := pgx.NamedArgs{
		"username": req.Username,
	}
	err := h.DB.QueryRow(r.Context(), "SELECT password_hash FROM users WHERE username = @username", args).Scan(&hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		http.Error(w, "Неверный пароль", http.StatusUnauthorized)
		return
	}

	var nextDay = time.Now().Add(time.Hour * 24)

	signedJwt, err := utils.GenerateJWT(jwt.MapClaims{
		"username": req.Username,
		"exp":      nextDay.Unix(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    signedJwt,
		Expires:  nextDay,
		MaxAge:   24 * 60 * 60,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Validate(w http.ResponseWriter, r *http.Request) {
	tokenStr := ""

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenStr = authHeader[7:]
	}

	if tokenStr == "" {
		http.Error(w, "Header \"Authorization\" or token not found", http.StatusUnauthorized)
		return
	}

	token, err := utils.ValidateJWT(tokenStr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
}
