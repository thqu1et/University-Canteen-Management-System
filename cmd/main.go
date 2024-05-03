package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"postgresSQLProject/pkg/database"
	"postgresSQLProject/pkg/routers"
	"postgresSQLProject/pkg/utils"
)

func main() {
	fmt.Println("Hello!")

	router := gin.New()
	router.Use(gin.Logger())
	routers.UserRoutes(router)
	routers.MenuItemRoutes(router)
	routers.OrderRoutes(router)
	routers.PaymentRoutes(router)
	//swagger.SwaggerRoutes(router)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run()
}

func init() {
	utils.LoadEnvFile()
	database.ConnectDB()
	database.SyncDataBase()
}
