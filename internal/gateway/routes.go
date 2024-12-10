package gateway

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/GoJobs/logger"
	"github.com/saleh-ghazimoradi/GoJobs/utils"
	"net/http"
)

func registerRoutes() *httprouter.Router {
	db, err := utils.PostConnection()
	if err != nil {
		logger.Logger.Error(err.Error())
	}

	fmt.Println(db)

	router := httprouter.New()
	router.NotFound = http.HandlerFunc(nil)
	router.MethodNotAllowed = http.HandlerFunc(nil)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", healthCheckHandler)

	return router
}
