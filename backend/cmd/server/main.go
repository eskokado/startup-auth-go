package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	handlers "github.com/eskokado/startup-auth-go/backend/internal/handlers/auth"
	"github.com/eskokado/startup-auth-go/backend/internal/middleware"
	"github.com/eskokado/startup-auth-go/backend/internal/providers"
	provider "github.com/eskokado/startup-auth-go/backend/internal/providers"
	repository "github.com/eskokado/startup-auth-go/backend/internal/repositories"
	usecase "github.com/eskokado/startup-auth-go/backend/internal/usecase/auth"
	service "github.com/eskokado/startup-auth-go/backend/pkg/domain/services"
)

func main() {
	// Carregar variáveis de ambiente
	_ = godotenv.Load(".env")
	sender := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		parsePort(os.Getenv("SMTP_PORT")),
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
	)

	// 1. Configurar o banco de dados (SQLite para exemplo)
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&repository.GormUser{})

	// 2. Inicializar repositório
	userRepo := repository.NewGormUserRepository(db)

	// 3. Inicializar serviços
	emailService := service.NewEmailService(sender)

	// 4. Inicializar redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Endereço do Redis
		Password: "",               // Senha
		DB:       0,                // Banco padrão
	})

	// 4. Inicializar provedores
	cryptoProvider := provider.NewBcryptProvider(bcrypt.DefaultCost)
	tokenProvider := provider.NewJWTProvider("secret-key", 24*time.Hour)
	blacklistProvider := providers.NewRedisBlacklist(rdb)

	// 5. Inicializar casos de uso
	registerUseCase := usecase.NewRegisterUsecase(userRepo, cryptoProvider)
	loggerUseCase := usecase.NewLoginUsecase(
		userRepo, cryptoProvider, tokenProvider, blacklistProvider,
	)
	logoutUseCase := usecase.NewLogoutUsecase(blacklistProvider)
	requestPasswordResetUC := usecase.NewRequestPasswordReset(userRepo, emailService)
	resetPasswordUC := usecase.NewResetPassword(userRepo)
	updateNameUC := usecase.NewUpdateNameUseCase(userRepo)
	updatePasswordUC := usecase.NewUpdatePasswordUseCase(userRepo, cryptoProvider)

	// 6. Criar handlers HTTP
	registerHTTPHandler := handlers.NewRegisterHandler(registerUseCase, userRepo)
	loggerHTTPHandler := handlers.NewLoginHandler(loggerUseCase)
	logoutHTTPHandler := handlers.NewLogoutHandler(logoutUseCase)
	forgotPasswordHandler := handlers.NewForgotPasswordHandler(requestPasswordResetUC)
	resetPasswordHandler := handlers.NewResetPasswordHandler(resetPasswordUC)
	updateNameHandler := handlers.NewUpdateNameHandler(updateNameUC)
	updatePasswordHandler := handlers.NewUpdatePasswordHandler(updatePasswordUC)

	// 7. Configurar roteador Gin
	router := gin.Default()

	// 7.1 Configurar CORS (ANTES dos middlewares de autenticação)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 7.2 Criar middleware de autenticação (DEPOIS do CORS)
	authMiddleware := middleware.JWTAuthMiddleware(tokenProvider, blacklistProvider)

	// 8. Registrar rotas
	router.POST("/auth/register", registerHTTPHandler.Handle)
	router.POST("/auth/login", loggerHTTPHandler.Handle)
	router.DELETE("/auth/logout", authMiddleware, logoutHTTPHandler.Handle)
	router.POST("/auth/forgot-password", forgotPasswordHandler.Handle)
	router.POST("/auth/reset-password", resetPasswordHandler.Handle)
	router.PUT("/user/name/:userID", authMiddleware, updateNameHandler.Handle)
	router.PUT("/user/password/:userID", authMiddleware, updatePasswordHandler.Handle)

	// 9. Iniciar o servidor
	router.Run(":8080")
}

func parsePort(port string) int {
	var p int
	fmt.Sscanf(port, "%d", &p)
	return p
}
