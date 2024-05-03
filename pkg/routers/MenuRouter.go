package routers

import (
	"github.com/gin-gonic/gin"
	"postgresSQLProject/pkg/controllers"
)

func MenuRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/menu", controllers.GetMenu)
	incomingRoutes.POST("/menu/create_menu", controllers.CreateMenu)
	incomingRoutes.PATCH("/menu/update_menu", controllers.UpdateMenu)
	incomingRoutes.DELETE("/menu/delete_menu", controllers.DeleteMenu)
}
