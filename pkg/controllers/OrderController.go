package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"postgresSQLProject/pkg/database"
	"postgresSQLProject/pkg/models"
	"strconv"
)

// GetOrders godoc
// @Summary Retrieve all orders
// @Description Get details of all orders with their associated items
// @Tags orders
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Order
// @Failure 500 {object} object {"error": "Failed to retrieve orders"}
// @Router /orders [get]
func GetOrders(c *gin.Context) {
	var orders []models.Order
	result := database.DB.Preload("OrderItems").Find(&orders)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// GetOrder godoc
// @Summary Retrieve a single order
// @Description Get details of a specific order by ID, including its items
// @Tags orders
// @Accept  json
// @Produce  json
// @Param id path int true "Order ID"
// @Success 200 {object} models.Order
// @Failure 404 {object} object {"error": "Order not found"}
// @Router /orders/{id} [get]
func GetOrder(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var order models.Order
	result := database.DB.Preload("OrderItems").First(&order, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Add a new order and calculate total from order items
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order body models.Order required "Order Info"
// @Success 201 {object} models.Order
// @Failure 400 {object} object {"error": "Invalid input, object invalid"}
// @Failure 500 {object} object {"error": "Failed to create order"}
// @Router /orders [post]
func CreateOrder(c *gin.Context) {
	var order models.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var total float64
	for _, item := range order.OrderItems {
		total += item.Price * float64(item.Quantity)
	}
	order.Total = total

	result := database.DB.Create(&order)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	database.DB.Preload("OrderItems.MenuItem").Find(&order)

	c.JSON(http.StatusCreated, order)
}

// UpdateOrder godoc
// @Summary Update an existing order
// @Description Update details of an existing order by ID
// @Tags orders
// @Accept  json
// @Produce  json
// @Param id path int true "Order ID"
// @Param order body models.Order required "Order Update Data"
// @Success 200 {object} models.Order
// @Failure 400 {object} object {"error": "Invalid input, object invalid"}
// @Failure 404 {object} object {"error": "Order not found"}
// @Failure 500 {object} object {"error": "Failed to update order"}
// @Router /orders/{id} [put]
func UpdateOrder(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := database.DB.Save(&order)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

// DeleteOrder godoc
// @Summary Delete an order
// @Description Delete an order by ID
// @Tags orders
// @Accept  json
// @Produce  json
// @Param id path int true "Order ID"
// @Success 200 {object} object {"message": "Order deleted"}
// @Failure 404 {object} object {"error": "No order found"}
// @Failure 500 {object} object {"error": "Failed to delete order"}
// @Router /orders/{id} [delete]
func DeleteOrder(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	result := database.DB.Delete(&models.Order{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No order found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order deleted"})
}
