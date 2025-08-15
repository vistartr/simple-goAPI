// // file: handlers/user_handler.go
// package handlers

// import (
// 	"database/sql"
// 	"fmt"
// 	"my-first-api/models"
// 	"my-first-api/repository"
// 	"net/http"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt/v5" // Impor baru
// 	"golang.org/x/crypto/bcrypt"
// )

// var jwtSecret = []byte("kunci_rahasia_yang_sangat_aman_dan_panjang")

// type JWTClaims struct {
// 	UserID   string `json:"user_id"`
// 	Username string `json:"username"`
// 	Role     string `json:"role"`
// 	jwt.RegisteredClaims
// }

// type UserHandler struct {
// 	UserRepo  *repository.UserRepository
// 	JwtSecret []byte
// }

// func NewUserHandler(userRepo *repository.UserRepository, jwtSecret []byte) *UserHandler {
// 	return &UserHandler{
// 		UserRepo:  userRepo,
// 		JwtSecret: jwtSecret,
// 	}
// }

// func (h *UserHandler) RegisterUser(c *gin.Context) {
// 	var newUser models.User
// 	if err := c.ShouldBindJSON(&newUser); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	// Hash password sebelum disimpan
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
// 		return
// 	}
// 	newUser.PasswordHash = string(hashedPassword)

// 	// Set role default jika tidak disediakan
// 	if newUser.Role == "" {
// 		newUser.Role = "visitor"
// 	}

// 	// Panggil repository untuk menyimpan user
// 	userID, err := h.UserRepo.SaveUser(newUser)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{
// 		"message": "User created successfully",
// 		"userID":  userID,
// 	})
// }

// func (h *UserHandler) LoginUser(c *gin.Context) {
// 	var loginAttempt models.User
// 	if err := c.ShouldBindJSON(&loginAttempt); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	// 1. Cari user di database berdasarkan username
// 	storedUser, err := h.UserRepo.FindByUsername(loginAttempt.Username)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
// 		return
// 	}

// 	// 2. Bandingkan password yang dikirim dengan hash di database
// 	err = bcrypt.CompareHashAndPassword([]byte(storedUser.PasswordHash), []byte(loginAttempt.Password))
// 	// Jika tidak cocok, errornya bisa `hashed password is not the hash of the given password`
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
// 		return
// 	}

// 	// 3. Jika password cocok, buat tokent JWT
// 	expirationTime := time.Now().Add(24 * time.Hour)

// 	claims := &JWTClaims{
// 		UserID:   strconv.Itoa(storedUser.ID),
// 		Username: storedUser.Username,
// 		Role:     storedUser.Role,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(expirationTime),
// 		},
// 		//Bisa menambahkan claims kustom misalkan role
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString(jwtSecret)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
// 		return
// 	}

// 	// 4. Kirim token sebagai response
// 	c.JSON(http.StatusOK, gin.H{"token": tokenString})

// }

// func (h *UserHandler) AuthMiddleWare() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if authHeader == "" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
// 			return
// 		}

// 		// 2. Header harus dalam format "Bearer <token>"
// 		//    Kita pisahkan untuk mendapatkan tokennya saja
// 		headerParts := strings.Split(authHeader, " ")
// 		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be in 'Bearer <token>' format"})
// 			return
// 		}

// 		tokenString := headerParts[1]
// 		claims := &jwt.RegisteredClaims{}

// 		// 3. Parse dan validasi token
// 		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 			// Pastikan metode signing-nya adalah yang kita gunakan (HS256)
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 			}
// 			return h.JwtSecret, nil
// 		})

// 		// 4. Periksa jika ada error saat parsing atau token tidak valid
// 		if err != nil || !token.Valid {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
// 			return
// 		}

// 		// 5. Jika token valid, simpan ID pengguna di context
// 		//    agar bisa digunakan oleh handler selanjutnya
// 		c.Set("userID", claims.Subject)

// 		// 6. Lanjutkan ke handler berikutnya
// 		c.Next()
// 	}
// }

package handlers

import (
	"database/sql"
	"fmt"
	"my-first-api/models"
	"my-first-api/repository"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Claims adalah data yang kita simpan di dalam JWT
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	newUser.PasswordHash = string(hashedPassword)
	if newUser.Role == "" {
		newUser.Role = "visitor"
	}

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

	storedUser, err := h.UserRepo.FindByUsername(loginAttempt.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.PasswordHash), []byte(loginAttempt.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &JWTClaims{
		UserID:   strconv.Itoa(storedUser.ID),
		Username: storedUser.Username,
		Role:     storedUser.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.JwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Middleware untuk memeriksa apakah pengguna adalah Admin
func (h *UserHandler) AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}
		tokenString := headerParts[1]

		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return h.JwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Periksa apakah peran pengguna adalah 'admin'
		if claims.Role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied: admin role required"})
			return // Hentikan jika bukan admin
		}

		// Lanjutkan jika dia adalah admin
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
