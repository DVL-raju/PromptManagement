package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// 1. Health check (GET/POST)
	mux.HandleFunc("/health", app.healthCheckHandler)

	// 2. Authentication (POST-only per service handler)
	mux.HandleFunc("/auth/register", app.auth.Register)
	mux.HandleFunc("/auth/login", app.auth.Login)

	// wrapping with middleware
	return app.recoverPanic(app.logRequest(mux))
}
