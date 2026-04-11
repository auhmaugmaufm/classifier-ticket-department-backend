// @title Predict Ticket Department API
// @version 1.0
// @description CTD Backend API
// @host localhost:8888
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"fmt"
	"log"

	_ "github.com/auhmaugmaufm/predict-ticket-department-backend/docs"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/auth"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/handler"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/middleware"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/repository"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/service"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/pkg/config"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/pkg/database"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	config.Load()
	cfg := config.Get()

	database.RunMigrations(cfg)
	db := database.NewPostgresDB(cfg)

	jwtManger := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpireHour)

	companyRepository := repository.NewCompanyRepository(db)
	companyService := service.NewCompanyService(companyRepository, jwtManger)
	companyHandler := handler.NewCompanyHandler(companyService, cfg)

	departmentRepository := repository.NewDepartmentRepository(db)
	departmentService := service.NewDepartmentService(departmentRepository)
	departmentHandler := handler.NewDepartmentHandler(departmentService, cfg)

	companyFormRepository := repository.NewCompanyFormRepository(db)
	companyFormService := service.NewCompanyFormService(companyFormRepository)
	companyFormHandler := handler.NewCompanyFormHandler(companyFormService, cfg)

	router := gin.Default()

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r := router.Group("/api/v1")
	r.POST("/register", companyHandler.Register)
	r.POST("/login", companyHandler.Login)

	protected := r.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManger))

	department := protected.Group("/departments")
	department.POST("/add", departmentHandler.AddDepartments)
	department.GET("/company/:company_id", departmentHandler.GetDepartmentsByCompanyID)

	companyForm := protected.Group("/company_form")
	companyForm.POST("/create", companyFormHandler.CreateCompanyForm)
	companyForm.GET("/company_form/:company_id", companyFormHandler.GetCompanyFormByCompanyID)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server running on %s", addr)
	log.Fatal(router.Run(addr))
}
