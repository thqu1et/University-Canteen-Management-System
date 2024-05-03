package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"postgresSQLProject/pkg/database"
	"postgresSQLProject/pkg/models"
	"strconv"
)

// ProcessPayment godoc
// @Summary Process a payment for an order
// @Description Process a payment by matching it with an order total and updating the payment status.
// @Tags payment
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param payment body models.Payment required "Payment Details"
// @Success 200 {object} object{"message": "Payment processed successfully", "order": "Order details"}
// @Failure 400 {object} object{"error": "Invalid order ID or Payment amount does not match the order total"}
// @Failure 404 {object} object{"error": "Order not found"}
// @Failure 500 {object} object{"error": "Internal server error"}
// @Router /payment/process/{id} [post]
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
