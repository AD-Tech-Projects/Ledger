package api

import (
	"net/http"

	logger "github.com/AD-Tech-Projects/Go-Libs/packages/logger"
	handlers "github.com/AD-Tech-Projects/Ledger/packages/handlers"
	"github.com/gin-gonic/gin"
)

func Listen(port string, router *gin.Engine) {
	logger.Info("Starting the application on port: " + port)

	address := ":" + port
	if err := router.Run(address); err != nil {
		logger.Error("Server failed to start:", err)
		panic(err)
	}
}

func NewRouter(handlers *handlers.HandlerResolver) *gin.Engine {
	router := gin.Default()

	router.GET("/health", healthCheckHandler)

	router.POST("/accounts", handlers.Account.Create)
	router.GET("/accounts", handlers.Account.List)
	router.GET("/accounts/:id", handlers.Account.GetById)
	router.PATCH("/accounts/:id", handlers.Account.Update)
	router.DELETE("/accounts/:id", handlers.Account.Delete)

	router.POST("/transactions", handlers.Transaction.Create)
	router.GET("/transactions", handlers.Transaction.List)
	router.GET("/transactions/:id", handlers.Transaction.GetById)

	return router
}

func healthCheckHandler(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
