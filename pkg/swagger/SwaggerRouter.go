package swagger

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SwaggerRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
