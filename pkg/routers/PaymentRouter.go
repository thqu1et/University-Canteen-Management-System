package routers

import (
	"github.com/gin-gonic/gin"
	"postgresSQLProject/pkg/controllers"
)

func PaymentRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/orders/:id/pay", controllers.ProcessPayment)
}
