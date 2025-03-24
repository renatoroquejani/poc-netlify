package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config contém todas as configurações necessárias para a aplicação
type Config struct {
	// Netlify
	NetlifyToken string
	BaseDomain   string

	// AWS
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
	S3BucketName       string
	S3Endpoint         string

	// API
	APIPort string

	// Aplicação
	Username         string
	CustomDomain     string
	S3Path           string
	NetlifySubdomain string
}

// LoadConfig carrega a configuração a partir de variáveis de ambiente e parâmetros
func LoadConfig() (*Config, error) {
	// Carregar variáveis de ambiente do arquivo .env, se existir
	_ = godotenv.Load()

	config := &Config{
		NetlifyToken:       os.Getenv("NETLIFY_TOKEN"),
		BaseDomain:         os.Getenv("BASE_DOMAIN"),
		AWSAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWSRegion:          os.Getenv("AWS_REGION"),
		S3BucketName:       os.Getenv("S3_BUCKET_NAME"),
		S3Endpoint:         os.Getenv("S3_ENDPOINT"),
		APIPort:            os.Getenv("API_PORT"),
	}

	// Definir valores padrão
	if config.APIPort == "" {
		config.APIPort = "8080"
	}

	// Definir valor padrão para o domínio base
	if config.BaseDomain == "" {
		config.BaseDomain = "sites.kodestech.com.br"
	}

	// Validar configurações obrigatórias
	if config.NetlifyToken == "" {
		return nil, fmt.Errorf("NETLIFY_TOKEN não definido")
	}

	if config.AWSAccessKeyID == "" || config.AWSSecretAccessKey == "" {
		return nil, fmt.Errorf("credenciais AWS não definidas")
	}

	// Se o endpoint S3 estiver definido (usando MinIO), não exigir região AWS
	if config.S3Endpoint == "" {
		if config.AWSRegion == "" {
			return nil, fmt.Errorf("AWS_REGION não definida")
		}
	}

	if config.S3BucketName == "" {
		return nil, fmt.Errorf("S3_BUCKET_NAME não definido")
	}

	return config, nil
}

// SetDeployParams configura os parâmetros específicos de deploy
func (c *Config) SetDeployParams(username, customDomain, s3Path string) error {
	if username == "" {
		return fmt.Errorf("nome de usuário não pode ser vazio")
	}

	c.Username = username
	c.CustomDomain = customDomain
	c.S3Path = s3Path

	// Gerar subdomínio da Netlify
	c.NetlifySubdomain = strings.ToLower(username) + "." + c.BaseDomain

	return nil
}
