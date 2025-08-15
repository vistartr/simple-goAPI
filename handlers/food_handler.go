package handlers

import (
	"database/sql"
	"my-first-api/models"
	"my-first-api/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FoodHandler struct {
	FoodRepo *repository.FoodRepository
}

func NewFoodHandler(foodRepo *repository.FoodRepository) *FoodHandler {
	return &FoodHandler{FoodRepo: foodRepo}
}

func (h *FoodHandler) GetFoods(c *gin.Context) {
	foods, err := h.FoodRepo.GetAllFoods()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, foods)
}

func (h *FoodHandler) GetFoodByID(c *gin.Context) {
	// Ambil ID dari parameter URL, contoh /food/5 -> id = "5"
	id := c.Param("id")

	// Panggil metode repository untuk mencari data
	food, err := h.FoodRepo.FindByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Food not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, food)
}

func (h *FoodHandler) PostFood(c *gin.Context) {
	var newFood models.Food
	if err := c.ShouldBindJSON(&newFood); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdFood, err := h.FoodRepo.CreateFood(newFood)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdFood)
}

func (h *FoodHandler) UpdateFood(c *gin.Context) {
	id := c.Param("id")
	var foodToUpdate models.Food

	if err := c.ShouldBindJSON(&foodToUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedFood, err := h.FoodRepo.UpdateFood(id, foodToUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedFood)
}

func (h *FoodHandler) DeleteFood(c *gin.Context) {
	id := c.Param("id")

	err := h.FoodRepo.DeleteFood(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Food deleted successfully"})
}
