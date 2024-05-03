package routers

import (
	"github.com/gin-gonic/gin"
	"postgresSQLProject/pkg/controllers"
)

func MenuItemRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/menu_item", controllers.GetMenuItems)
	incomingRoutes.GET("/menu_item/:menuitem_id", controllers.GetMenuItem)
	incomingRoutes.POST("/menu/create_menu_item", controllers.CreateMenuItem)
	incomingRoutes.PUT("/menu/update_menu_item/:menuitem_id", controllers.UpdateMenuItem)
	incomingRoutes.DELETE("/menu/delete_menu_item/:menuitem_id", controllers.DeleteMenuItem)
}
