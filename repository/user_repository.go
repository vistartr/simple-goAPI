// // file: repository/user_repository.go
// package repository

// import (
// 	"database/sql"
// 	"my-first-api/models"
// )

// type UserRepository struct {
// 	DB *sql.DB
// }

// func NewUserRepository(db *sql.DB) *UserRepository {
// 	return &UserRepository{DB: db}
// }

// // Fungsi untuk menyimpan user baru ke database
// func (r *UserRepository) SaveUser(user models.User) (int, error) {
// 	sqlStatement := `INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3) RETURNING id`

// 	var userID int
// 	// Kita simpan password_hash dari struct User, dan role dari struct User
// 	err := r.DB.QueryRow(sqlStatement, user.Username, user.PasswordHash, user.Role).Scan(&userID)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return userID, nil
// }

// func (r *UserRepository) FindByUsername(username string) (models.User, error) {
// 	var user models.User
// 	sqlStatement := `SELECT id, username, password_hash, role FROM users WHERE username = $1`

// 	err := r.DB.QueryRow(sqlStatement, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
// 	if err != nil {
// 		return user, err
// 	}
// 	return user, nil
// }

package repository

import (
	"database/sql"
	"my-first-api/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Fungsi untuk menyimpan user baru ke database
func (r *UserRepository) SaveUser(user models.User) (int, error) {
	sqlStatement := `INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3) RETURNING id`

	var userID int
	err := r.DB.QueryRow(sqlStatement, user.Username, user.PasswordHash, user.Role).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// Fungsi untuk mencari pengguna berdasarkan username-nya
func (r *UserRepository) FindByUsername(username string) (models.User, error) {
	var user models.User
	sqlStatement := `SELECT id, username, password_hash, role FROM users WHERE username = $1`

	err := r.DB.QueryRow(sqlStatement, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
	if err != nil {
		return user, err
	}
	return user, nil
}
