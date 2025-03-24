package main

// @title Netlify Deploy API
// @version 1.0
// @description API para facilitar o deploy de arquivos estáticos na Netlify a partir de um bucket S3
// @host localhost:8080
// @BasePath /
// @schemes http https

import (
	"log"
	"os"

	"github.com/kodestech/poc-netlify/internal/api"
	"github.com/kodestech/poc-netlify/internal/config"
	"io"
)

func main() {
	// Configurar o logger para escrever em um arquivo
	logFile, err := os.OpenFile("netlify-deploy.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Erro ao abrir arquivo de log: %v", err)
	}
	defer logFile.Close()
	
	// Configurar o logger para escrever no terminal e no arquivo
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	
	log.Println("Iniciando aplicação Netlify Deploy")
	
	// Carregar configurações
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}
	
	// Criar servidor API
	s := api.NewServer(cfg)

	// Iniciar servidor
	log.Printf("Iniciando servidor...")
	if err := s.Start(); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
