package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
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

	router := gin.Default()

	router.Run(":8080")
}

func parsePort(port string) int {
	var p int
	fmt.Sscanf(port, "%d", &p)
	return p
}
