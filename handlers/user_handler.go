// file: handlers/user_handler.go
package handlers

import (
	"database/sql"
	"my-first-api/models"
	"my-first-api/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5" // Impor baru
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("kunci_rahasia_yang_sangat_aman_dan_panjang")

type UserHandler struct {
	UserRepo  *repository.UserRepository
	JwtSecret []byte
}

func NewUserHandler(userRepo *repository.UserRepository, jwtSecret []byte) *UserHandler {
	return &UserHandler{
		UserRepo:  userRepo,
		JwtSecret: jwtSecret,
	}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	newUser.PasswordHash = string(hashedPassword)

	// Set role default jika tidak disediakan
	if newUser.Role == "" {
		newUser.Role = "visitor"
	}

	// Panggil repository untuk menyimpan user
	userID, err := h.UserRepo.SaveUser(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"userID":  userID,
	})
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var loginAttempt models.User
	if err := c.ShouldBindJSON(&loginAttempt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 1. Cari user di database berdasarkan username
	storedUser, err := h.UserRepo.FindByUsername(loginAttempt.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// 2. Bandingkan password yang dikirim dengan hash di database
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.PasswordHash), []byte(loginAttempt.Password))
	// Jika tidak cocok, errornya bisa `hashed password is not the hash of the given password`
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 3. Jika password cocok, buat tokent JWT
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.RegisteredClaims{
		Subject:   strconv.Itoa(storedUser.ID),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		//Bisa menambahkan claims kustom misalkan role
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// 4. Kirim token sebagai response
	c.JSON(http.StatusOK, gin.H{"token": tokenString})

}
