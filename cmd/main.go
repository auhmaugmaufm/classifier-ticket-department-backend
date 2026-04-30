// @title Predict Ticket Department API
// @version 1.0
// @description CTD Backend API
// @host localhost:8888
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @securityDefinitions.apikey HMACAuth
// @in header
// @name X-HMAC-Signature
package main

import (
	"fmt"
	"log"

	_ "github.com/auhmaugmaufm/predict-ticket-department-backend/docs"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/auth"
	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/cron"
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

	txManager := repository.NewTxManager(db)

	companyRepository := repository.NewCompanyRepository(db)
	departmentRepository := repository.NewDepartmentRepository(db)
	companyService := service.NewCompanyService(companyRepository, departmentRepository, txManager, jwtManger)
	companyHandler := handler.NewCompanyHandler(companyService, cfg)

	departmentService := service.NewDepartmentService(departmentRepository)
	departmentHandler := handler.NewDepartmentHandler(departmentService, cfg)

	formRepository := repository.NewFormRepository(db)
	formService := service.NewFormService(formRepository)
	formHandler := handler.NewFormHandler(formService, cfg)

	LinkRepository := repository.NewLinkRepository(db)
	LinkService := service.NewLinkService(LinkRepository)
	LinkHandler := handler.NewLinkHandler(LinkService, cfg)

	aiService := service.NewAIService(
		cfg.AIBackendUrl,
		cfg.HMACSecret,
	)

	ticketRepository := repository.NewTicketRepositry(db)
	ticketService := service.NewTicketService(ticketRepository)
	ticketHandler := handler.NewTicketHandler(ticketService, cfg)

	formCron := cron.NewFormCron(formService, companyService, aiService)
	formCron.Start()
	defer formCron.Stop()

	router := gin.Default()

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Submit Form: post /form/create
	router.POST("/form/submit", formHandler.SubmitForm)

	r := router.Group("/api/v1")
	r.POST("/register", companyHandler.Register)
	r.POST("/login", companyHandler.Login)
	// r.POST("/create-bulk", ticketHandler.CreateTickets)
	// r.GET("/departments/:company_id", departmentHandler.GetDepartmentsByCompanyID)

	protected := r.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManger))

	department := protected.Group("/departments")
	department.POST("/add", departmentHandler.AddDepartments)
	// department.GET("/:company_id", departmentHandler.GetDepartmentsByCompanyID)

	Link := protected.Group("/links")
	Link.POST("/create", LinkHandler.CreateLink)
	Link.GET("/:company_id", LinkHandler.GetLinkByCompanyID)

	forms := protected.Group("/forms")
	forms.GET("/:company_id", formHandler.GetSubmitFormCompanyID)
	forms.GET("/:company_id/per-day", formHandler.GetSubmitFormPerDayByCompanyID)

	ticket := protected.Group("/tickets")
	ticket.POST("/create", ticketHandler.CreateTicket)
	ticket.POST("/create-bulk", ticketHandler.CreateTickets)
	ticket.GET("/:company_id", ticketHandler.GetTicketsByCompanyID)

	internal_protected := r.Group("")
	internal_protected.Use(middleware.HMACMiddleware(cfg.HMACSecret))
	internal_protected.POST("/create-bulk", ticketHandler.CreateTickets)
	internal_protected.GET("/departments/:company_id", departmentHandler.GetDepartmentsByCompanyID)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server running on %s", addr)
	log.Fatal(router.Run(addr))
}
