// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"strconv" // Kita butuh ini untuk mengubah string ID menjadi integer

// 	"github.com/gin-contrib/cors" // <-- IMPORT CORS YANG BARU
// 	"github.com/gin-gonic/gin"    // Import Gin
// 	_ "github.com/lib/pq"

// )

// //Import Gin

// type Food struct {
// 	ID          int     `json: "id"`
// 	Name        string  `json: "name"`
// 	Description string  `json: "description"`
// 	Price       float64 `json: "price"`
// }

// var db *sql.DB

// func main() {
// 	connStr := "user=vistartr password=Password1! dbname=postgres sslmode=disable"
// 	var err error
// 	db, err = sql.Open("postgres", connStr)
// 	if err != nil {
// 		log.Fatal("Gagal terhubung ke database:", err)
// 	}
// 	defer db.Close()

// 	err = db.Ping()
// 	if err != nil {
// 		log.Fatal("Database tidak meresponse:", err)
// 	}
// 	fmt.Println("Berhasil terhubung ke database!")

// 	// 1. Buat router Gin
// 	router := gin.Default()

// 	// 2. Gunakan middleware CORS
// 	// (Cara Gin sedikit berbeda)
// 	config := cors.DefaultConfig()
// 	config.AllowOrigins = true
// 	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
// 	router.Use(cors.New(config))

// 	// 3. Definisikan route pake Gin
// 	router.GET("/food", getFoods)
// 	router.POST("/food", postFood)
// 	router.PUT("/food/:id", updateFood)
// 	router.DELETE("food/:id", deleteFood)

// 	// 4. Jalankan server Gin
// 	fmt.Println("Server Gin listening on port 8080...")
// 	router.Run(":8080")
// }

// func getFoods(c *gin.Context) {
// 	rows, err := db.Query("Select id, name, description, price FROM foods ORDER BY id ASC")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	defer rows.Close()

// 	foods := []Food{}
// 	for rows.Next() {
// 		var food Food
// 		if err := rows.Scan(&food.ID, &food.Name, &food.Description, &food.Price); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 			return
// 		}
// 		foods = append(foods, food)
// 	}
// 	c.JSON(http.StatusOK, foods)
// }

// func postFood(c *gin.Context) {
// 	var newFood Food
// 	if err := c.ShouldBindJSON(&newFood); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
// 		return
// 	}

// 	sqlStatement := `INSER INTO foods (name, deescription, price) VALUES ($1, $2, $3) RETURNING id`
// 	var insertedID int
// 	err := db.QueryRow(sqlStatement, newFood.Name, newFood.Description, newFood.Price).Scan(&insertedID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	newFood.ID = insertedID
// 	c.JSON(http.StatusCreated, newFood)
// }

// func updateFood(c *gin.Context) {
// 	// 1. Ambil ID dari URL, misalnya "/food/1" -> id = "1"
// 	id := c.Param("id")
// 	var updatedFood Food

// 	// 2. Ambil data baru dari body JSON
// 	if err := c.ShouldBindJSON(&updatedFood); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// 3. Siapkan dan jalankan perintah SQL UPDATE
// 	sqlStatement := `UPDATE foods SET name = $1, description = $2, price = $3 WHERE id = $4`
// 	_, err := db.Exec(sqlStatement, updatedFood.Name, updatedFood.Description, updatedFood.Price, id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Set ID dari URL ke objek yang akan dikirim kembali
// 	// (Karena c.Param mengembalikan string, kita tidak usah repot mengubahnya jadi int)
// 	updatedFood.ID = 0 // Kita tidak tahu ID nya, bisa di-parse dari string 'id' jika perlu
// 	c.JSON(http.StatusOK, updatedFood)
// }

// // --- HANDLER BARU UNTUK DELETE ---
// func deleteFood(c *gin.Context) {
// 	// 1. Ambil ID dari URL
// 	id := c.Param("id")

// 	// 2. Siapkan dan jalankan perintah SQL DELETE
// 	sqlStatement := `DELETE FROM foods WHERE id = $1`
// 	_, err := db.Exec(sqlStatement, id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// 3. Kirim pesan sukses
// 	c.JSON(http.StatusOK, gin.H{"message": "Food deleted successfully"})
// }

package main

import (
	"database/sql"
	// "encoding/json" // tidak lagi dibutuhkan di sini
	"fmt"
	"log"
	"net/http"
	"strconv" // Kita butuh ini untuk mengubah string ID menjadi integer

	"github.com/gin-contrib/cors" // <-- IMPORT CORS YANG BARU
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	// "github.com/rs/cors" // <-- HAPUS IMPORT CORS YANG LAMA
)

type Food struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

var db *sql.DB

func main() {
	connStr := "user=vistartr password=Password1! dbname=postgres sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Database tidak merespon:", err)
	}
	fmt.Println("Berhasil terhubung ke database!")

	router := gin.Default()

	// --- KODE CORS BARU YANG BENAR UNTUK GIN ---
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Izinkan semua origin, bisa diganti dengan alamat React App Anda
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	router.GET("/food", getFoods)
	router.GET("/food/:id", getFoodByID)
	router.POST("/food", postFood)
	router.PUT("/food/:id", updateFood)
	router.DELETE("/food/:id", deleteFood)

	fmt.Println("Server Gin listening on port 8080...")
	router.Run(":8080")
}

// Handler GET
func getFoods(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, description, price FROM foods ORDER BY id ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	foods := []Food{}
	for rows.Next() {
		var food Food
		if err := rows.Scan(&food.ID, &food.Name, &food.Description, &food.Price); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		foods = append(foods, food)
	}
	c.JSON(http.StatusOK, foods)
}

// Handler POST
func postFood(c *gin.Context) {
	var newFood Food
	if err := c.ShouldBindJSON(&newFood); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sqlStatement := `INSERT INTO foods (name, description, price) VALUES ($1, $2, $3) RETURNING id`
	var insertedID int
	err := db.QueryRow(sqlStatement, newFood.Name, newFood.Description, newFood.Price).Scan(&insertedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newFood.ID = insertedID
	c.JSON(http.StatusCreated, newFood)
}

// Handler UPDATE
func updateFood(c *gin.Context) {
	idStr := c.Param("id")
	var updatedFood Food

	if err := c.ShouldBindJSON(&updatedFood); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sqlStatement := `UPDATE foods SET name = $1, description = $2, price = $3 WHERE id = $4`
	_, err := db.Exec(sqlStatement, updatedFood.Name, updatedFood.Description, updatedFood.Price, idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Konversi id dari string ke int untuk dimasukkan kembali ke objek
	id, _ := strconv.Atoi(idStr)
	updatedFood.ID = id
	c.JSON(http.StatusOK, updatedFood)
}

// Handler DELETE
func deleteFood(c *gin.Context) {
	id := c.Param("id")

	sqlStatement := `DELETE FROM foods WHERE id = $1`
	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Food deleted successfully"})
}

func getFoodByID(c *gin.Context) {
	id := c.Param("id")

	var food Food

	sqlStatement := `SELECT id, name, description, price FROM foods WHERE id = $1`

	//Query row untuk mengharapkan suatu hasil
	err := db.QueryRow(sqlStatement, id).Scan(&food.ID, &food.Name, &food.Description, &food.Price)
	if err != nil {
		// Jika tidak ada baris yang ditemukan, kirim error 404 Not Found
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Food not found"})
			return
		}
		// Untuk Error lainnya
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, food)
}
