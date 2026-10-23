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
	mux.HandleFunc("GET /accounts", handlers.Account.List)
	mux.HandleFunc("GET /accounts/{id}", handlers.Account.GetById)
	mux.HandleFunc("PATCH /accounts/{id}", handlers.Account.Update)
	mux.HandleFunc("DELETE /accounts/{id}", handlers.Account.Delete)

	mux.HandleFunc("POST /transactions", handlers.Transaction.Create)
	mux.HandleFunc("GET /transactions", handlers.Transaction.List)
	mux.HandleFunc("GET /transactions/{id}", handlers.Transaction.GetById)

	return mux
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
