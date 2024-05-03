package routers

func PaymentRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/orders/:id/pay", controllers.ProcessPayment)
}
