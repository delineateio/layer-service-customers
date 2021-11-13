package routes

import (
	"fmt"
	"net/http"

	"delineate.io/customers/src/config"
	"delineate.io/customers/src/logging"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "delineate.io/customers/src/docs"
)

func NewRouter() *gin.Engine {
	mode := config.GetStringOrDefault("server.mode", "release")
	gin.SetMode(mode)

	router := gin.New()
	router.Use(cors.Default())
	router.Use(logging.GetRouterLogger())
	router.Use(gin.CustomRecovery(func(ctx *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			ctx.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}))

	// custom routes
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/healthz", healthz)
	router.POST("/customer", createCustomer)
	router.GET("/customer/:id", getCustomerByID)
	router.GET("/customers", getCustomers)

	return router
}
