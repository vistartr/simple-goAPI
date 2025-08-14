package repository

import (
	"database/sql"
	"my-first-api/models"
)

type FoodRepository struct {
	DB *sql.DB
}

func NewFoodRepository(db *sql.DB) *FoodRepository {
	return &FoodRepository{DB: db}
}

func (r *FoodRepository) GetAllFoods() ([]models.Food, error) {
	rows, err := r.DB.Query("SELECT id, name, description, price FROM foods ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var foods []models.Food
	for rows.Next() {
		var food models.Food
		if err := rows.Scan(&food.ID, &food.Name, &food.Description, &food.Price); err != nil {
			return nil, err
		}
		foods = append(foods, food)
	}
	return foods, nil
}

func (r *FoodRepository) FindByID(id string) (models.Food, error) {
	var food models.Food
	// Siapkan query SQL untuk mengambil salah satu baris berdasarkan ID
	sqlStatement := `SELECT id, name, description, price FROM foods WHERE id = $1`

	// Gunakan QueryRow karena kita hanya mengharapkan satu hasil.
	// Scan akan memasukkan hasilnya ke struct food.
	err := r.DB.QueryRow(sqlStatement, id).Scan(&food.ID, &food.Name, &food.Description, &food.Price)

	// Kembalikan struct food dan error (jika ada)
	return food, err
}

func (r *FoodRepository) CreateFood(food models.Food) (models.Food, error) {
	sqlStatement := `INSERT INTO foods (name, description, price) VALUES ($1, $2, $3) RETURNING id`

	var insertedID int
	err := r.DB.QueryRow(sqlStatement, food.Name, food.Description, food.Price).Scan(&insertedID)
	if err != nil {
		return food, err
	}

	food.ID = insertedID
	return food, nil
}

func (r *FoodRepository) UpdateFood(id string, food models.Food) (models.Food, error) {
	sqlStatement := `UPDATE foods SET name = $1, description = $2, price = $3 WHERE id = $4`

	_, err := r.DB.Exec(sqlStatement, food.Name, food.Description, food.Price, id)
	if err != nil {
		return food, err
	}

	// Karena Exec tidak mengembalikan data, kita kembalikan saja data yang diterima
	// dengan asumsi update berhasil.
	return food, nil
}

func (r *FoodRepository) DeleteFood(id string) error {
	sqlStatement := `DELETE FROM foods WHERE id = $1`

	_, err := r.DB.Exec(sqlStatement, id)
	return err
}
