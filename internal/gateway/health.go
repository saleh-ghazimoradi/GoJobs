package gateway

import (
	"github.com/saleh-ghazimoradi/GoJobs/config"
	"net/http"
)

// healthCheckHandler handles the health check request and returns the application's status, environment, and version.
// @Summary Health check endpoint
// @Description Returns the current status of the application, including its environment and version details.
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string "Health check status, environment, and version"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /v1/healthcheck [get]
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     config.AppConfig.ServerConfig.Port,
		"version": config.AppConfig.ServerConfig.Version,
	}
	if err := jsonResponse(w, http.StatusOK, data); err != nil {
		internalServerError(w, r, err)
	}
}
