package gateway

import (
	"github.com/saleh-ghazimoradi/GoJobs/config"
	"net/http"
)

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
