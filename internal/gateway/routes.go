package gateway

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/GoJobs/internal/repository"
	"github.com/saleh-ghazimoradi/GoJobs/internal/service"
	"github.com/saleh-ghazimoradi/GoJobs/logger"
	"github.com/saleh-ghazimoradi/GoJobs/utils"
	"net/http"
)

func registerRoutes() http.Handler {
	db, err := utils.PostConnection()
	if err != nil {
		logger.Logger.Error(err.Error())
	}

	userDB := repository.NewUserRepository(db, db)
	userService := service.NewUserService(userDB)
	userHandler := NewUserHandler(userService)
	authService := service.NewAuthenticateService(userDB)
	authHandler := NewAuthenticateHandler(authService)

	jobDB := repository.NewJobRepository(db, db)
	jobService := service.NewJobService(jobDB)
	jobHandler := NewJob(jobService)

	router := httprouter.New()
	router.NotFound = http.HandlerFunc(notFoundRouter)
	router.MethodNotAllowed = http.HandlerFunc(methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", healthCheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/login", authHandler.loginHandler)
	router.HandlerFunc(http.MethodPost, "/v1/register", authHandler.registerHandler)
	router.Handler(http.MethodGet, "/v1/users/:id", AuthMiddleware(http.HandlerFunc(userHandler.getUserByIdHandler)))
	router.Handler(http.MethodPut, "/v1/users/:id", AuthMiddleware(http.HandlerFunc(userHandler.UpdateUserProfileHandler)))
	router.Handler(http.MethodPost, "/v1/users/:id/picture", AuthMiddleware(http.HandlerFunc(userHandler.UpdateUserProfilePictureHandler)))

	router.HandlerFunc(http.MethodGet, "/v1/jobs", jobHandler.GetAllJobsHandler)
	router.Handler(http.MethodPost, "/v1/jobs", AuthMiddleware(http.HandlerFunc(jobHandler.CreateJobHandler)))
	router.Handler(http.MethodGet, "/v1/jobsByUser", AuthMiddleware(http.HandlerFunc(jobHandler.GetAllJobsHandler)))
	router.Handler(http.MethodGet, "/v1/jobs/:id", AuthMiddleware(http.HandlerFunc(jobHandler.GetJobByIdHandler)))
	router.Handler(http.MethodPut, "/v1/jobs/:id", AuthMiddleware(http.HandlerFunc(jobHandler.UpdateJobHandler)))

	return router
}
