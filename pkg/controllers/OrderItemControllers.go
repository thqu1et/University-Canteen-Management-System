package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"postgresSQLProject/pkg/database"
	"postgresSQLProject/pkg/models"
)

// GetOrderItems returns all order items for a specific order
func GetOrderItems(c *gin.Context) {
	var orderItems []models.OrderItem
	orderID := c.Param("order_id") // Assuming you pass order ID in the path

	result := database.DB.Where("order_id = ?", orderID).Find(&orderItems)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, orderItems)
}

// GetOrderItem retrieves a single order item by its ID
func GetOrderItem(c *gin.Context) {
	var orderItem models.OrderItem
	orderItemID := c.Param("id")

	result := database.DB.First(&orderItem, orderItemID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, orderItem)
}

// CreateOrderItem handles the creation of a new order item
func CreateOrderItem(c *gin.Context) {
	var orderItem models.OrderItem
	if err := c.ShouldBindJSON(&orderItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Create(&orderItem)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, orderItem)
}

// UpdateOrderItem modifies an existing order item
func UpdateOrderItem(c *gin.Context) {
	orderItemID := c.Param("id")
	var orderItem models.OrderItem

	if err := database.DB.First(&orderItem, orderItemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order item not found"})
		return
	}

	if err := c.ShouldBindJSON(&orderItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Save(&orderItem)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, orderItem)
}

// DeleteOrderItem removes an order item
func DeleteOrderItem(c *gin.Context) {
	orderItemID := c.Param("id")
	result := database.DB.Delete(&models.OrderItem{}, orderItemID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No order item found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order item deleted"})
}
