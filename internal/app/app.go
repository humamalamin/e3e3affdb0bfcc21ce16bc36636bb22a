package app

import (
	"context"
	"latihan-portal-news/config"
	"latihan-portal-news/internal/adapter/cloudflare"
	"latihan-portal-news/internal/adapter/handler"
	"latihan-portal-news/internal/adapter/repository"
	"latihan-portal-news/internal/core/service"
	authLib "latihan-portal-news/lib/auth"
	middlewareLib "latihan-portal-news/lib/middleware"
	"latihan-portal-news/lib/pagination"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func RunServer() {
	cfg := config.NewConfig()
	db, err := cfg.ConnectionPostgres()
	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}

	// Create Temporary Directory
	err = os.MkdirAll("./temp/content", 0755)
	if err != nil {
		log.Fatal("error creating temporary directory: ", err.Error())
	}

	// Cloudflare R2
	cfgR2 := cfg.LoadAwsConfig()
	s3Client := s3.NewFromConfig(cfgR2)
	clouflareR2 := cloudflare.NewCloudflareR2Repository(s3Client, &cfg)

	// JWT and Middleware
	jwtAuth := authLib.NewJwt(&cfg)
	middlewareAuth := middlewareLib.NewMiddleware(&cfg)

	// Pagination Lib
	pageHelper := pagination.NewPagination()

	// Repository Module App
	authRepo := repository.NewAuthRepository(db.DB)
	categoryRepo := repository.NewCategoryRepository(db.DB)
	contentRepo := repository.NewContentRepository(db.DB)
	userRepo := repository.NewUserRepository(db.DB)

	// Service Module App
	authService := service.NewAuthService(authRepo, &cfg, jwtAuth)
	categoryService := service.NewCategoryService(categoryRepo)
	contentService := service.NewContentService(contentRepo, &cfg, clouflareR2)
	userService := service.NewUserService(userRepo)

	// Handler Module App
	authHandler := handler.NewAuthhandler(authService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	contentHandler := handler.NewContentHandler(contentService, pageHelper)
	userHandler := handler.NewUserHandler(userService)

	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip}  ${status} - ${latency} ${method} ${path}\n",
	}))

	if os.Getenv("APP_ENV") != "production" {
		confiSwg := swagger.Config{
			BasePath: "/api",
			FilePath: "./docs/swagger.json",
			Path:     "docs",
			Title:    "Swagger API News Portal",
		}
		app.Use(swagger.New(confiSwg))
	}

	apiGroup := app.Group("/api")
	apiGroup.Post("/login", authHandler.Login)

	adminApp := apiGroup.Group("/admin")
	adminApp.Use(middlewareAuth.CheckToken())

	// Category Admin
	categoryApp := adminApp.Group("/categories")
	categoryApp.Get("/", categoryHandler.GetCategories)
	categoryApp.Post("/", categoryHandler.CreateCategory)
	categoryApp.Put("/:categoryID", categoryHandler.EditCategoryByID)
	categoryApp.Delete("/:categoryID", categoryHandler.DeleteCategoryByID)
	categoryApp.Get("/:categoryID", categoryHandler.GetCategoryByID)

	// Content Admin
	contentApp := adminApp.Group("/contents")
	contentApp.Get("/", contentHandler.GetContents)
	contentApp.Post("/", contentHandler.CreateContent)
	contentApp.Put("/:contentID", contentHandler.UpdateContent)
	contentApp.Delete("/:contentID", contentHandler.DeleteContentByID)
	contentApp.Get("/:contentID", contentHandler.GetContentByID)
	contentApp.Post("/upload-image", contentHandler.UploadImageCloudFlareR2)

	// Profile
	userApp := adminApp.Group("/users")
	userApp.Get("/profile", userHandler.GetUserByID)
	userApp.Put("/update-password", userHandler.UpdatePassword)

	// FE
	feApp := apiGroup.Group("/fe")
	feApp.Get("/contents", contentHandler.GetContentByQuery)

	go func() {
		if cfg.App.AppPort == "" {
			cfg.App.AppEnv = os.Getenv("APP_PORT")
		}

		err := app.Listen(":" + cfg.App.AppPort)
		if err != nil {
			log.Fatal("error starting server: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	// Block until a signal is received.
	<-quit

	log.Println("server shutdown of 5 second.")
	// gracefully shutdown the server, waiting max 5 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app.ShutdownWithContext(ctx)
}
