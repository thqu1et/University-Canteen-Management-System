package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"postgresSQLProject/pkg/database"
	"postgresSQLProject/pkg/routers"
	"postgresSQLProject/pkg/utils"
)

func main() {
	fmt.Println("Hello!")

	router := gin.New()
	router.Use(gin.Logger())
	routers.UserRoutes(router)

	router.Run()
}

func init() {
	utils.LoadEnvFile()
	database.ConnectDB()
	database.SyncDataBase()
}
