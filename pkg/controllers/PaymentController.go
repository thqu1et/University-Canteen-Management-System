package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"postgresSQLProject/pkg/database"
	"postgresSQLProject/pkg/models"
	"strconv"
)

func ProcessPayment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var payment models.Payment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if order.Total != payment.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment amount does not match the order total"})
		return
	}

	payment.OrderID = uint(id) // Make sure the payment is linked to the order
	if err := models.ProcessPayment(database.DB, payment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	database.DB.Preload("OrderItems").Preload("OrderItems.MenuItem").First(&order, id)

	c.JSON(http.StatusOK, gin.H{"message": "Payment processed successfully", "order": order})
}
