package routers

import (
	"github.com/gin-gonic/gin"
	"postgresSQLProject/pkg/controllers"
)

func OrderItemRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/order-items/:order_id", controllers.GetOrderItems)
	incomingRoutes.GET("/order-item/:id", controllers.GetOrderItem)
	incomingRoutes.POST("/order-item", controllers.CreateOrderItem)
	incomingRoutes.PUT("/order-item/:id", controllers.UpdateOrderItem)
	incomingRoutes.DELETE("/order-item/:id", controllers.DeleteOrderItem)
}
