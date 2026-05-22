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

	aiService := service.NewAIService(
		cfg.AIBackendUrl,
		cfg.HMACSecret,
	)

	ticketRepository := repository.NewTicketRepository(db)
	ticketService := service.NewTicketService(ticketRepository)
	ticketHandler := handler.NewTicketHandler(ticketService, cfg)

	formCron := cron.NewFormCron(formService, companyService, aiService)
	formCron.Start()
	defer formCron.Stop()

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	// Swagger
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})
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
	department.GET("", departmentHandler.GetDepartmentsByCompanyIDAuth)
	department.POST("", departmentHandler.CreateDepartment)
	department.PATCH("/:id", departmentHandler.UpdateDepartmentStatus)

	forms := protected.Group("/forms")
	forms.GET("", formHandler.GetSubmitFormCompanyID)
	forms.GET("/per-day", formHandler.GetSubmitFormPerDayByCompanyID)

	ticket := protected.Group("/tickets")
	ticket.GET("", ticketHandler.GetTicketsByCompanyID)
	ticket.POST("", ticketHandler.CreateTicket)
	ticket.POST("/bulk", ticketHandler.CreateTickets)
	ticket.PATCH("/:id", ticketHandler.UpdateTicketStatus)

	internal_protected := r.Group("internal")
	internal_protected.Use(middleware.HMACMiddleware(cfg.HMACSecret))
	internal_protected.POST("/tickets/bulk", ticketHandler.CreateTickets)
	internal_protected.GET("/departments/:company_id", departmentHandler.GetDepartmentsByCompanyID)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server running on %s", addr)
	log.Fatal(router.Run(addr))
}
