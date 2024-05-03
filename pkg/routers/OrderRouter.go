package routers

import (
	"github.com/gin-gonic/gin"
	"postgresSQLProject/pkg/controllers"
)

func OrderRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/orders", controllers.GetOrders)
	incomingRoutes.GET("/orders/:id", controllers.GetOrder)
	incomingRoutes.POST("/orders", controllers.CreateOrder)
	incomingRoutes.PUT("/orders/:id", controllers.UpdateOrder)
	incomingRoutes.DELETE("/orders/:id", controllers.DeleteOrder)
}
