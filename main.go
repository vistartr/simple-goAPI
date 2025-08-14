// import (
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"

// 	_ "github.com/lib/pq"
// 	"github.com/rs/cors"
// )

// // Struct Food dengan ID
// type Food struct {
// 	ID          int     `json:"id"`
// 	Name        string  `json:"name"`
// 	Description string  `json:"description"`
// 	Price       float64 `json:"price"`
// }

// // Variabel global untuk koneksi database
// var db *sql.DB

// func main() {
// 	// Ganti dengan password PostgreSQL Anda
// 	connStr := "user=vistartr password=Password1! dbname=postgres sslmode=disable"
// 	var err error
// 	db, err = sql.Open("postgres", connStr)
// 	if err != nil {
// 		log.Fatal("Gagal terhubung ke database:", err)
// 	}
// 	defer db.Close()

// 	err = db.Ping()
// 	if err != nil {
// 		log.Fatal("Database tidak merespon:", err)
// 	}
// 	fmt.Println("Berhasil terhubung ke database!")

// 	// Setup router dan server
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/food", handleFood)
// 	handler := cors.Default().Handler(mux)

// 	fmt.Println("Server is listening on port 8080...")
// 	log.Fatal(http.ListenAndServe(":8080", handler))
// }

// func handleFood(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	switch r.Method {
// 	case http.MethodGet:
// 		getFoods(w, r)
// 	case http.MethodPost:
// 		postFood(w, r)
// 	default:
// 		http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
// 	}
// }

// func getFoods(w http.ResponseWriter, r *http.Request) {
// 	rows, err := db.Query("SELECT id, name, description, price FROM foods ORDER BY id ASC")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	// PERBAIKAN DI BARIS INI
// 	foods := []Food{}

// 	for rows.Next() {
// 		var food Food
// 		if err := rows.Scan(&food.ID, &food.Name, &food.Description, &food.Price); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		foods = append(foods, food)
// 	}

// 	json.NewEncoder(w).Encode(foods)
// }

// func postFood(w http.ResponseWriter, r *http.Request) {
// 	var newFood Food
// 	if err := json.NewDecoder(r.Body).Decode(&newFood); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	sqlStatement := `
// 		INSERT INTO foods (name, description, price)
// 		VALUES ($1, $2, $3)
// 		RETURNING id`

// 	var insertedID int
// 	err := db.QueryRow(sqlStatement, newFood.Name, newFood.Description, newFood.Price).Scan(&insertedID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	newFood.ID = insertedID
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(newFood)
// }

package main

import (
	// "database/sql"
	"fmt"
	"my-first-api/config"
	"my-first-api/handlers"
	"my-first-api/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var jwtSecret = []byte("kunci_rahasia_yang_sangat_aman_dan_panjang")

// token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIiwiZXhwIjoxNzU1MjYzNjM1fQ._yoHSUDt0dhHyg42JDeFVMuSTvxza-N8DZdOXJGjVE0"

func main() {

	// 1. Connect database
	db := config.ConnectDatabase()
	defer db.Close()

	// 2. Buat instance dari Repository
	foodRepo := repository.NewFoodRepository(db)
	userRepo := repository.NewUserRepository(db)

	// 3. Buat instance dari Handler
	foodHandler := handlers.NewFoodHandler(foodRepo)
	userHandler := handlers.NewUserHandler(userRepo, jwtSecret)

	// 4. Setup Router Gin
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// 5. Definisikan Rute, menunjuk ke fungsi di dalam Handler
	router.GET("/food", foodHandler.GetFoods)
	router.GET("food/:id", foodHandler.GetFoodByID)
	router.POST("/food", foodHandler.PostFood)
	router.PUT("/food/:id", foodHandler.UpdateFood)
	router.DELETE("/food/:id", foodHandler.DeleteFood)

	router.POST("/register", userHandler.RegisterUser)
	router.POST("/login", userHandler.LoginUser)

	fmt.Println("Server Gin listening on port 8080...")
	router.Run(":8080")
}
