package api

import (
	"net/http"

	logger "github.com/AD-Tech-Projects/Go-Libs/packages/logger"
	handlers "github.com/AD-Tech-Projects/Ledger/packages/handlers"
)

func Listen(port string, router http.Handler) {
	logger.Info("Starting the application on port: " + port)

	address := ":" + port
	if err := http.ListenAndServe(address, router); err != nil {
		logger.Error("Server failed to start:", err)
		panic(err)
	}
}

func NewRouter(handlers *handlers.HandlerResolver) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc(
		"GET /health",
		healthCheckHandler,
	)

	mux.HandleFunc("POST /accounts", handlers.Account.Create)

	mux.HandleFunc("POST /transactions", handlers.Transaction.Create)

	return mux
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
