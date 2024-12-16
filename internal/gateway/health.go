package gateway

import (
	"github.com/saleh-ghazimoradi/GoJobs/config"
	"net/http"
)

// healthCheckHandler godoc
// @Summary Health Check Endpoint
// @Description This endpoint returns the health status of the application, including environment and version information.
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "OK"
// @Failure 500 {object} map[string]string "Internal Server Error"
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
