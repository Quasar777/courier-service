package healthroute

import (
	"github.com/go-chi/chi/v5"

	"github.com/Quasar777/courier-service/internal/handler/healthhttp"
)

// HealthRoute - роуты для health-проверок
func HealthRoute(mr *chi.Mux, handler *healthhttp.HealthHandler) {
	mr.Get("/ping", handler.Ping)
	mr.Head("/healthcheck", handler.HealthCheck)
}
