package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	handlers "github.com/eskokado/startup-auth-go/backend/internal/handlers/auth"
	crypto "github.com/eskokado/startup-auth-go/backend/internal/providers"
	repository "github.com/eskokado/startup-auth-go/backend/internal/repositories"
	usecase "github.com/eskokado/startup-auth-go/backend/internal/usecase/auth"
)

func main() {
	// Carregar variáveis de ambiente
	_ = godotenv.Load(".env")
	_ = gomail.NewDialer(
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

	// 4. Inicializar provedores
	cryptoProvider := crypto.NewBcryptProvider(bcrypt.DefaultCost) // Implementação do CryptoProvider

	// 5. Inicializar casos de uso
	registerUseCase := usecase.NewRegisterUsecase(userRepo, cryptoProvider)

	// 6. Criar handlers HTTP
	registerHTTPHandler := handlers.NewRegisterHandler(registerUseCase, userRepo)

	// 7. Configurar roteador Gin
	router := gin.Default()

	// 7.1 Criar middleware

	// 8. Registrar rotas
	router.POST("/auth/register", registerHTTPHandler.Handle)

	router.Run(":8080")
}

func parsePort(port string) int {
	var p int
	fmt.Sscanf(port, "%d", &p)
	return p
}
