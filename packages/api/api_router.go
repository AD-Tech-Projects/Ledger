package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	logger "github.com/AD-Tech-Projects/Go-Libs/packages/logger"
	handlers "github.com/AD-Tech-Projects/Ledger/packages/handlers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Listen(port string, router *gin.Engine, log logger.ContextLogger) {
	log.Info("Starting the application on port: " + port)

	address := ":" + port
	if err := router.Run(address); err != nil {
		log.Error("Server failed to start:", err)
		panic(err)
	}
}

func NewRouter(handlers *handlers.HandlerResolver, log logger.ContextLogger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(requestLogger(log), gin.CustomRecovery(func(c *gin.Context, recovered any) {
		err := fmt.Errorf("panic recovered: %v", recovered)
		log.Error(err)
		_ = c.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

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

func requestLogger(log logger.ContextLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.NewString()
		startedAt := time.Now()
		body := readRequestBody(c)
		request := requestLogData(c, body, requestID)

		log.Http(
			"REQUESTED - ["+c.Request.Method+"] "+c.Request.URL.Path,
			logJSON(request, log),
		)

		c.Next()

		response := map[string]any{
			"status":    c.Writer.Status(),
			"ip":        c.ClientIP(),
			"requestId": requestID,
		}
		action := "RESPONDING"
		if c.Writer.Status() >= http.StatusBadRequest || len(c.Errors) > 0 {
			action = "RESPONDING ERROR"
			if len(c.Errors) > 0 {
				response["error"] = c.Errors.Last().Error()
			}
		}

		elapsed := time.Since(startedAt).Milliseconds()
		log.Http(
			fmt.Sprintf("%s - [%s] %s (%d ms):", action, c.Request.Method, c.Request.URL.Path, elapsed),
			logJSON(response, log),
		)
	}
}

func logJSON(value any, log logger.ContextLogger) string {
	data, err := json.Marshal(value)
	if err != nil {
		log.Error("Failed to serialize HTTP log data:", err)
		return "{}"
	}
	return string(data)
}

func readRequestBody(c *gin.Context) any {
	if c.Request.Body == nil {
		return map[string]any{}
	}

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return map[string]any{}
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(data))
	if len(bytes.TrimSpace(data)) == 0 {
		return map[string]any{}
	}

	var body any
	if err := json.Unmarshal(data, &body); err != nil {
		return string(data)
	}
	return body
}

func requestLogData(c *gin.Context, body any, requestID string) map[string]any {
	pathParams := map[string]string{}
	for _, param := range c.Params {
		pathParams[param.Key] = param.Value
	}

	queryParams := map[string]any{}
	for key, values := range c.Request.URL.Query() {
		if len(values) == 1 {
			queryParams[key] = values[0]
			continue
		}
		queryParams[key] = values
	}

	return map[string]any{
		"pathParams":  pathParams,
		"queryParams": queryParams,
		"body":        body,
		"ip":          c.ClientIP(),
		"requestId":   requestID,
	}
}
