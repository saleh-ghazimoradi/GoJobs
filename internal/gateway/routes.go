package gateway

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/GoJobs/config"
	"github.com/saleh-ghazimoradi/GoJobs/docs"
	_ "github.com/saleh-ghazimoradi/GoJobs/docs"
	"github.com/saleh-ghazimoradi/GoJobs/internal/repository"
	"github.com/saleh-ghazimoradi/GoJobs/internal/service"
	"github.com/saleh-ghazimoradi/GoJobs/logger"
	"github.com/saleh-ghazimoradi/GoJobs/utils"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

// @title Golang Web API
// @description This is a web API server.
// @version 1.0
// @host localhost:8080
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func registerRoutes() http.Handler {
	db, err := utils.PostConnection()
	if err != nil {
		logger.Logger.Error(err.Error())
	}

	userDB := repository.NewUserRepository(db, db)
	jobDB := repository.NewJobRepository(db, db)

	userService := service.NewUserService(userDB)
	jobService := service.NewJobService(jobDB)
	authService := service.NewAuthenticateService(userDB)

	userHandler := NewUserHandler(userService)
	jobHandler := NewJob(jobService)
	authHandler := NewAuthenticateHandler(authService)

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(notFoundRouter)
	router.MethodNotAllowed = http.HandlerFunc(methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", healthCheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/forgotpassword", authHandler.ForgotPasswordHandler)

	router.HandlerFunc(http.MethodPost, "/v1/login", authHandler.loginHandler)
	router.HandlerFunc(http.MethodPost, "/v1/register", authHandler.registerHandler)

	router.Handler(http.MethodGet, "/v1/users/:id", AuthMiddleware(http.HandlerFunc(userHandler.getUserByIdHandler)))
	router.Handler(http.MethodPut, "/v1/users/:id", AuthMiddleware(http.HandlerFunc(userHandler.UpdateUserProfileHandler)))
	router.Handler(http.MethodPost, "/v1/users/:id/picture", AuthMiddleware(http.HandlerFunc(userHandler.UpdateUserProfilePictureHandler)))
	router.Handler(http.MethodGet, "/v1/users", AuthMiddleware(http.HandlerFunc(userHandler.GetAllUsersHandler)))
	router.Handler(http.MethodDelete, "/v1/users/:id", AuthMiddleware(http.HandlerFunc(userHandler.DeleteUserHandler)))
	router.Handler(http.MethodPut, "/v1/users/:id/changePassword", AuthMiddleware(http.HandlerFunc(userHandler.ChangePasswordHandler)))

	router.HandlerFunc(http.MethodGet, "/v1/jobs", jobHandler.GetAllJobsHandler)
	router.Handler(http.MethodPost, "/v1/jobs", AuthMiddleware(http.HandlerFunc(jobHandler.CreateJobHandler)))
	router.Handler(http.MethodGet, "/v1/jobsByUser", AuthMiddleware(http.HandlerFunc(jobHandler.GetAllJobsHandler)))
	router.Handler(http.MethodGet, "/v1/jobs/:id", AuthMiddleware(http.HandlerFunc(jobHandler.GetJobByIdHandler)))
	router.Handler(http.MethodPut, "/v1/jobs/:id", AuthMiddleware(http.HandlerFunc(jobHandler.UpdateJobHandler)))
	router.Handler(http.MethodDelete, "/v1/jobs/:id", AuthMiddleware(http.HandlerFunc(jobHandler.DeleteJobHandler)))

	swaggerHandler := SetupSwagger()
	router.Handler(http.MethodGet, "/swagger/*any", swaggerHandler)

	return router
}

// SetupSwagger
// Swagger UI Route
// @Summary Swagger Documentation
// @Description Provides access to the Swagger UI
// @Accept json
// @Produce json
// @Success 200 {string} string "Swagger UI"
// @Router /swagger [get]
func SetupSwagger() http.Handler {
	docs.SwaggerInfo.Title = "Golang Web API"
	docs.SwaggerInfo.Description = "This is a web API server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost%s", config.AppConfig.ServerConfig.Port)
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	return httpSwagger.WrapHandler
}
