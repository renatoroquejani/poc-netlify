package api

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kodestech/poc-netlify/internal/aws"
	"github.com/kodestech/poc-netlify/internal/config"
	"github.com/kodestech/poc-netlify/internal/netlify"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/kodestech/poc-netlify/docs"
	"github.com/joho/godotenv"
	"github.com/netlify/open-api/go/models"
)

// Server representa o servidor da API
type Server struct {
	config *config.Config
	router *gin.Engine
}

// DeployRequest representa os parâmetros para um deploy (mantido para compatibilidade)
type DeployRequest struct {
	Username     string `json:"username"`
	CustomDomain string `json:"custom_domain,omitempty"`
	S3Path       string `json:"s3_path"`
}

// DeployResponse representa a resposta de uma operação de deploy (mantido para compatibilidade)
type DeployResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	SiteURL      string `json:"site_url,omitempty"`
	Subdomain    string `json:"subdomain,omitempty"`
	CustomDomain string `json:"custom_domain,omitempty"`
	DeployID     string `json:"deploy_id,omitempty"`
}

// TestDeployRequest encapsula os parâmetros para teste de deploy na Netlify
type TestDeployRequest struct {
	SiteID          string `json:"site_id" example:"a1b2c3d4" swagger:"description=ID do site na Netlify para atualização (opcional)"`
	SiteName        string `json:"site_name" binding:"required" example:"test-site" swagger:"description=Nome do site para teste"`
	Description     string `json:"description" example:"Site de teste" swagger:"description=Descrição do site para teste"`
	TestContent     string `json:"test_content" example:"<h1>Hello World</h1>" swagger:"description=Conteúdo HTML para teste (opcional)"`
	CleanupAfter    bool   `json:"cleanup_after" example:"true" swagger:"description=Remover o site após o teste"`
	CustomDomain    string `json:"custom_domain" example:"meu-site.exemplo.com" swagger:"description=Domínio personalizado para o site (opcional)"`
	FileContent     string `json:"file_content" example:"" swagger:"description=Conteúdo de arquivo HTML em formato base64 (opcional, alternativa ao TestContent)"`
}

// TestDeployResponse representa a resposta de uma operação de teste de deploy
type TestDeployResponse struct {
	Success     bool      `json:"success" example:"true" swagger:"description=Indica se o teste foi bem-sucedido"`
	Message     string    `json:"message" example:"Site de teste criado com sucesso" swagger:"description=Mensagem do resultado do teste"`
	SiteID      string    `json:"site_id,omitempty" example:"a1b2c3d4" swagger:"description=ID do site criado na Netlify"`
	SiteURL     string    `json:"site_url,omitempty" example:"https://test-site.netlify.app" swagger:"description=URL do site criado"`
	CreatedAt   string    `json:"created_at,omitempty" example:"2023-10-25T15:30:45Z" swagger:"description=Data e hora de criação do site"`
	TestSuccess bool      `json:"test_success" example:"true" swagger:"description=Indica se o teste específico foi bem-sucedido"`
}

// DomainRequest representa uma requisição para modificar um domínio
type DomainRequest struct {
	SiteID string `json:"site_id" binding:"required" example:"e17e2166-d8ab-4cad-9916-a9a3fed7750d"`
	Domain string `json:"domain" binding:"omitempty" example:"example.com"`
}

// DomainResponse representa a resposta de operações de gerenciamento de domínios
type DomainResponse struct {
	Success bool   `json:"success" example:"true" swagger:"description=Indica se a operação foi bem-sucedida"`
	Message string `json:"message" example:"Domínio adicionado com sucesso" swagger:"description=Mensagem descritiva sobre o resultado da operação"`
	SiteID  string `json:"site_id" example:"e17e2166-d8ab-4cad-9916-a9a3fed7750d" swagger:"description=ID do site na Netlify"`
	Domain  string `json:"domain" example:"example.com" swagger:"description=Domínio personalizado gerenciado"`
}

// NewServer cria um novo servidor da API
func NewServer(cfg *config.Config) *Server {
	// Configurar o modo do Gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Criar router Gin
	router := gin.Default()

	// Configurar CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Adicionar middleware de registro
	router.Use(func(c *gin.Context) {
		// Registrar request
		start := time.Now()
		c.Next()
		// Registrar response
		log.Printf("[API] %s %s %d %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(start),
		)
		c.Next()
	})

	// Configurar Swagger
	router.GET("/docs/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// URL do Swagger UI: http://localhost:8080/docs/swagger/index.html
	server := &Server{
		config: cfg,
		router: router,
	}

	// Configurar rotas
	server.setupRoutes()

	return server
}

// setupRoutes configura as rotas da API
func (s *Server) setupRoutes() {
	// Grupo de rotas para a API
	apiGroup := s.router.Group("/api")
	{
		// Rota de status
		// @Summary Verificar status da API
		// @Description Retorna o status atual da API, versão e timestamp
		// @Tags status
		// @Produce json
		// @Success 200 {object} map[string]interface{}
		// @Router /api/status [get]
		apiGroup.GET("/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "online",
				"version": "1.0.0",
				"time":    time.Now().Format(time.RFC3339),
			})
		})

		// As rotas de deploy foram removidas conforme solicitado

		// Rota para deploy de site
		// @Summary Cria ou atualiza sites na Netlify
		// @Description Cria um novo site na Netlify ou atualiza um existente quando o ID é fornecido
		// @Tags deploy
		// @Accept multipart/form-data
		// @Produce json
		// @Param site_id formData string false "ID do site na Netlify para atualização (opcional)"
		// @Param site_name formData string true "Nome do site para teste"
		// @Param description formData string false "Descrição do site para teste"
		// @Param cleanup_after formData bool false "Remover o site após o teste"
		// @Param custom_domain formData string false "Domínio personalizado para o site (opcional)"
		// @Param file formData file false "Arquivo HTML para deploy (opcional)"
		// @Param folder_path formData string false "Caminho da pasta local para deploy (opcional)"
		// @Success 200 {object} TestDeployResponse
		// @Failure 400 {object} TestDeployResponse
		// @Failure 500 {object} TestDeployResponse
		// @Router /api/deploy/site [post]
		apiGroup.POST("/deploy/site", s.handleTestDeploy)

		// Rota para adicionar domínio personalizado
		// @Summary Adiciona um domínio personalizado a um site
		// @Description Adiciona um domínio personalizado como alias para um site existente na Netlify
		// @Tags domains
		// @Accept json
		// @Produce json
		// @Param request body DomainRequest true "Dados do domínio a ser adicionado"
		// @Success 200 {object} DomainResponse
		// @Failure 400 {object} DomainResponse
		// @Failure 500 {object} DomainResponse
		// @Router /api/domains/add [post]
		apiGroup.POST("/domains/add", s.handleAddDomain)

		// Rota para remover domínio personalizado
		// @Summary Remove um domínio personalizado de um site
		// @Description Remove um domínio personalizado dos aliases de um site existente na Netlify
		// @Tags domains
		// @Accept json
		// @Produce json
		// @Param request body DomainRequest true "Dados do domínio a ser removido"
		// @Success 200 {object} DomainResponse
		// @Failure 400 {object} DomainResponse
		// @Failure 500 {object} DomainResponse
		// @Router /api/domains/remove [post]
		apiGroup.POST("/domains/remove", s.handleRemoveDomain)

		// Rota para definir domínio como principal
		// @Summary Define um domínio como o domínio principal
		// @Description Define um domínio personalizado como o domínio principal de um site existente na Netlify
		// @Tags domains
		// @Accept json
		// @Produce json
		// @Param request body DomainRequest true "Dados do domínio a ser definido como principal"
		// @Success 200 {object} DomainResponse
		// @Failure 400 {object} DomainResponse
		// @Failure 500 {object} DomainResponse
		// @Router /api/domains/set-default [post]
		apiGroup.POST("/domains/set-default", s.handleSetDefaultDomain)

		// Rota para remover domínio principal
		// @Summary Remove o domínio principal de um site
		// @Description Remove o domínio principal de um site existente na Netlify, mantendo os aliases
		// @Tags domains
		// @Accept json
		// @Produce json
		// @Param request body DomainRequest true "ID do site a ter o domínio principal removido"
		// @Success 200 {object} DomainResponse
		// @Failure 400 {object} DomainResponse
		// @Failure 500 {object} DomainResponse
		// @Router /api/domains/remove-primary [post]
		apiGroup.POST("/domains/remove-primary", s.handleRemovePrimaryDomain)

		// Rota para testar conexão com a API da Netlify
		// @Summary Testa a conexão com a API da Netlify
		// @Description Testa a conexão com a API da Netlify e exibe informações sobre o token
		// @Tags netlify
		// @Accept json
		// @Produce json
		// @Success 200 {object} map[string]interface{}
		// @Failure 500 {object} map[string]interface{}
		// @Router /api/test/netlify/connection [get]
		apiGroup.GET("/test/netlify/connection", s.handleTestNetlifyConnection)

		// Rota para exibir logs da aplicação
		// @Summary Exibe os logs da aplicação
		// @Description Retorna os logs da aplicação em formato texto
		// @Tags logs
		// @Produce text
		// @Success 200 {string} string "Logs da aplicação"
		// @Router /api/test/logs [get]
		apiGroup.GET("/test/logs", s.handleViewLogs)

		// Rota para listar sites da Netlify
		// @Summary Lista todos os sites do usuário na Netlify
		// @Description Retorna a lista de todos os sites do usuário na Netlify
		// @Tags netlify
		// @Accept json
		// @Produce json
		// @Success 200 {object} map[string]interface{}
		// @Failure 500 {object} map[string]interface{}
		// @Router /api/sites [get]
		apiGroup.GET("/sites", s.handleListSites)
	}

	// Servir arquivos estáticos para a interface web
	s.router.Static("/static", "./web/static")
	
	// Redirecionar raiz para o index.html
	// @Summary Redirecionar para a interface web
	// @Description Redireciona a rota raiz para a interface web estática
	// @Tags interface
	// @Produce html
	// @Success 301 {string} string "Redirecionamento"
	// @Router / [get]
	s.router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/static")
	})
}

// handleDeploy processa uma requisição de deploy (rota desativada)
func (s *Server) handleDeploy(c *gin.Context) {
	var req DeployRequest

	// Fazer log da requisição recebida
	log.Printf("[API] Recebida requisição de deploy: %+v", c.Request)
	
	// Fazer bind dos dados da requisição
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[API] Erro ao processar JSON da requisição: %v", err)
		c.JSON(http.StatusBadRequest, DeployResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao processar requisição: %v", err),
		})
		return
	}

	log.Printf("[API] Dados de deploy validados: username=%s, customDomain=%s, s3Path=%s", req.Username, req.CustomDomain, req.S3Path)

	// Configurar o cliente da Netlify
	netlifyClient, err := netlify.NewClient(s.config)
	if err != nil {
		log.Printf("[API] Erro ao criar cliente Netlify: %v", err)
		c.JSON(http.StatusInternalServerError, DeployResponse{
			Success: false,
			Message: fmt.Sprintf("Erro interno do servidor: %v", err),
		})
		return
	}
	
	// Configurar parâmetros de deploy no objeto config
	if err := s.config.SetDeployParams(req.Username, req.CustomDomain, req.S3Path); err != nil {
		log.Printf("[API] Erro ao configurar parâmetros de deploy: %v", err)
		c.JSON(http.StatusBadRequest, DeployResponse{
			Success: false,
			Message: fmt.Sprintf("Erro nos parâmetros de deploy: %v", err),
		})
		return
	}

	// Verificar se o site já existe
	site, exists, err := netlifyClient.VerifySite(c.Request.Context())
	if err != nil {
		log.Printf("[API] Erro ao verificar site na Netlify: %v", err)
		c.JSON(http.StatusInternalServerError, DeployResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao verificar site na Netlify: %v", err),
		})
		return
	}

	// Se o site não existir, criá-lo
	if !exists {
		log.Printf("[API] Site não encontrado. Criando novo site para %s", req.Username)
		site, err = netlifyClient.CreateSite(c.Request.Context())
		if err != nil {
			log.Printf("[API] Erro ao criar site na Netlify: %v", err)
			c.JSON(http.StatusInternalServerError, DeployResponse{
				Success: false,
				Message: fmt.Sprintf("Erro ao criar site na Netlify: %v", err),
			})
			return
		}
		log.Printf("[API] Site criado com sucesso: %s", site.CustomDomain)
	} else {
		log.Printf("[API] Site existente encontrado: %s", site.CustomDomain)
	}

	// Configurar cliente S3
	_, err = aws.NewS3Client(s.config)
	if err != nil {
		log.Printf("[API] Erro ao criar cliente S3: %v", err)
		c.JSON(http.StatusInternalServerError, DeployResponse{
			Success: false,
			Message: fmt.Sprintf("Erro interno do servidor: %v", err),
		})
		return
	}
	
	// TODO: Implementar o deploy efetivo dos arquivos
	// Por enquanto, apenas retornamos sucesso
	log.Printf("[API] Deploy iniciado com sucesso para %s.%s", req.Username, s.config.BaseDomain)

	// Retornar resposta de sucesso
	c.JSON(http.StatusAccepted, DeployResponse{
		Success:      true,
		Message:      "Deploy iniciado com sucesso",
		Subdomain:    fmt.Sprintf("%s.%s", req.Username, s.config.BaseDomain),
		CustomDomain: req.CustomDomain,
		// Deploy ID e URL do site seriam definidos aqui em uma implementação real
	})
}

// handleTestDeploy processa uma requisição de deploy de site na Netlify
// @Summary Cria ou atualiza sites na Netlify
// @Description Cria um novo site na Netlify ou atualiza um existente quando o ID é fornecido
// @Tags deploy
// @Accept multipart/form-data
// @Produce json
// @Param site_id formData string false "ID do site na Netlify para atualização (opcional)"
// @Param site_name formData string true "Nome do site para teste"
// @Param description formData string false "Descrição do site para teste"
// @Param cleanup_after formData bool false "Remover o site após o teste"
// @Param custom_domain formData string false "Domínio personalizado para o site (opcional)"
// @Param file formData file false "Arquivo HTML para deploy (opcional)"
// @Param folder_path formData string false "Caminho da pasta local para deploy (opcional)"
// @Success 200 {object} TestDeployResponse
// @Failure 400 {object} TestDeployResponse
// @Failure 500 {object} TestDeployResponse
// @Router /api/deploy/site [post]
func (s *Server) handleTestDeploy(c *gin.Context) {
	// Processar upload de arquivo
	log.Printf("[TEST-API] Recebida requisição para teste de deploy: %s", c.Request.URL.Path)
	s.handleMultipartTestDeploy(c)
}

// handleMultipartTestDeploy processa uma requisição de teste de deploy com upload de arquivo
func (s *Server) handleMultipartTestDeploy(c *gin.Context) {
	// Obter parâmetros do formulário
	siteID := c.PostForm("site_id")
	siteName := c.PostForm("site_name")
	if siteName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Nome do site é obrigatório",
		})
		return
	}

	description := c.PostForm("description")
	cleanupAfterStr := c.PostForm("cleanup_after")
	customDomain := c.PostForm("custom_domain")
	testContent := c.PostForm("test_content")

	// Converter valores booleanos
	cleanupAfter := false
	if cleanupAfterStr != "" {
		cleanupAfter, _ = strconv.ParseBool(cleanupAfterStr)
	}

	folderPath := c.PostForm("folder_path")

	// Verificar se o caminho da pasta foi fornecido e se existe
	if folderPath != "" {
		// Verificar se a pasta existe
		if _, err := os.Stat(folderPath); os.IsNotExist(err) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fmt.Sprintf("Pasta não encontrada: %s", folderPath),
			})
			return
		}
	}

	// Obter o arquivo enviado
	file, err := c.FormFile("file")
	if err != nil && folderPath == "" && testContent == "" {
		// Se não foi fornecido arquivo, pasta ou conteúdo de teste, continuamos com um conteúdo padrão
		log.Printf("[TEST-API] Nenhum arquivo, pasta ou conteúdo fornecido para deploy, usando conteúdo padrão")
	}

	// Variável para armazenar o conteúdo do arquivo em base64
	var fileContentBase64 string

	// Se um arquivo foi enviado, processar
	if err == nil {
		log.Printf("[TEST-API] Arquivo recebido: %s (%d bytes)", file.Filename, file.Size)

		// Abrir o arquivo
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": fmt.Sprintf("Erro ao abrir arquivo: %v", err),
			})
			return
		}
		defer src.Close()

		// Ler o conteúdo do arquivo
		fileContent, err := io.ReadAll(src)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": fmt.Sprintf("Erro ao ler arquivo: %v", err),
			})
			return
		}

		// Converter para base64
		fileContentBase64 = base64.StdEncoding.EncodeToString(fileContent)
		log.Printf("[TEST-API] Arquivo convertido para base64 (%d bytes)", len(fileContentBase64))
	}

	// Criar cliente Netlify
	netlifyClient, err := netlify.NewClient(s.config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Erro ao criar cliente Netlify: %v", err),
		})
		return
	}

	// Executar o teste de deploy
	result, err := netlifyClient.ExecuteTestDeploy(context.Background(), netlify.TestDeployParams{
		SiteID:          siteID,
		SiteName:        siteName,
		Description:     description,
		TestContent:     testContent,
		CleanupAfter:    cleanupAfter,
		CustomDomain:    customDomain,
		FileContent:     fileContentBase64,
		FolderPath:      folderPath,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Erro ao executar teste de deploy: %v", err),
		})
		return
	}

	// Retornar resultado
	c.JSON(http.StatusOK, result)
}

// handleAddDomain adiciona um domínio personalizado a um site
// @Summary Adiciona um domínio personalizado a um site
// @Description Adiciona um domínio personalizado como alias para um site existente na Netlify
// @Tags domains
// @Accept json
// @Produce json
// @Param request body DomainRequest true "Dados do domínio a ser adicionado"
// @Success 200 {object} DomainResponse
// @Failure 400 {object} DomainResponse
// @Failure 500 {object} DomainResponse
// @Router /api/domains/add [post]
func (s *Server) handleAddDomain(c *gin.Context) {
	log.Printf("[handleAddDomain] Recebendo requisição para adicionar domínio")
	
	var req DomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[handleAddDomain] Erro ao processar JSON: %v", err)
		c.JSON(http.StatusBadRequest, DomainResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao processar requisição: %v", err),
		})
		return
	}

	log.Printf("[handleAddDomain] Requisição recebida: SiteID=%s, Domain=%s", req.SiteID, req.Domain)
	
	// Validar campos obrigatórios
	if req.SiteID == "" || req.Domain == "" {
		log.Printf("[handleAddDomain] Campos obrigatórios não informados")
		c.JSON(http.StatusBadRequest, DomainResponse{
			Success: false,
			Message: "ID do site e domínio são obrigatórios",
			SiteID:  req.SiteID,
			Domain:  req.Domain,
		})
		return
	}

	// Criar cliente Netlify
	log.Printf("[handleAddDomain] Criando cliente Netlify")
	netlifyClient, err := netlify.NewClient(s.config)
	if err != nil {
		log.Printf("[handleAddDomain] Erro ao criar cliente Netlify: %v", err)
		c.JSON(http.StatusInternalServerError, DomainResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao criar cliente Netlify: %v", err),
			SiteID:  req.SiteID,
			Domain:  req.Domain,
		})
		return
	}

	// Obter informações do site primeiro para verificar se tem domínio principal
	sites, err := netlifyClient.ListSites(context.Background())
	if err != nil {
		log.Printf("[handleAddDomain] Erro ao listar sites: %v", err)
		c.JSON(http.StatusInternalServerError, DomainResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao verificar sites: %v", err),
			SiteID:  req.SiteID,
			Domain:  req.Domain,
		})
		return
	}

	// Encontrar o site com o ID especificado
	var siteFound bool
	var customDomain string
	for _, site := range sites {
		if site.ID == req.SiteID {
			siteFound = true
			customDomain = site.CustomDomain
			break
		}
	}

	// Verificar se o site existe
	if !siteFound {
		log.Printf("[handleAddDomain] Site com ID %s não encontrado", req.SiteID)
		c.JSON(http.StatusBadRequest, DomainResponse{
			Success: false,
			Message: fmt.Sprintf("Site com ID %s não encontrado", req.SiteID),
			SiteID:  req.SiteID,
			Domain:  req.Domain,
		})
		return
	}

	// Verificar se o site já tem domínio principal configurado
	if customDomain != "" {
		log.Printf("[handleAddDomain] Site %s já tem domínio principal configurado: %s", req.SiteID, customDomain)
		c.JSON(http.StatusBadRequest, DomainResponse{
			Success: false,
			Message: fmt.Sprintf("O site já tem um domínio principal configurado: %s. Remova o domínio principal antes de adicionar um alias.", customDomain),
			SiteID:  req.SiteID,
			Domain:  req.Domain,
		})
		return
	}

	// Adicionar domínio
	log.Printf("[handleAddDomain] Adicionando domínio %s ao site %s", req.Domain, req.SiteID)
	
	// Criar contexto com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	err = netlifyClient.AddCustomDomain(ctx, req.SiteID, req.Domain)
	if err != nil {
		log.Printf("[handleAddDomain] Erro ao adicionar domínio: %v", err)
		
		// Verificar mensagens de erro específicas
		errMsg := err.Error()
		if strings.Contains(errMsg, "domínio principal configurado") || 
		   strings.Contains(errMsg, "domínio principal") || 
		   strings.Contains(errMsg, "update domain aliases while primary") {
		   
			c.JSON(http.StatusBadRequest, DomainResponse{
				Success: false,
				Message: "O site já tem um domínio principal configurado. Remova o domínio principal antes de adicionar um alias.",
				SiteID:  req.SiteID,
				Domain:  req.Domain,
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, DomainResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao adicionar domínio: %v", err),
			SiteID:  req.SiteID,
			Domain:  req.Domain,
		})
		return
	}

	log.Printf("[handleAddDomain] Domínio adicionado com sucesso")
	c.JSON(http.StatusOK, DomainResponse{
		Success: true,
		Message: "Domínio adicionado com sucesso",
		SiteID:  req.SiteID,
		Domain:  req.Domain,
	})
}

// handleRemoveDomain remove um domínio personalizado de um site
// @Summary Remove um domínio personalizado de um site
// @Description Remove um domínio personalizado dos aliases de um site existente na Netlify
// @Tags domains
// @Accept json
// @Produce json
// @Param request body DomainRequest true "Dados do domínio a ser removido"
// @Success 200 {object} DomainResponse
// @Failure 400 {object} DomainResponse
// @Failure 500 {object} DomainResponse
// @Router /api/domains/remove [post]
func (s *Server) handleRemoveDomain(c *gin.Context) {
	log.Printf("[handleRemoveDomain] Recebendo requisição para remover domínio")
	
	var req DomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[handleRemoveDomain] Erro ao processar JSON: %v", err)
		c.JSON(http.StatusBadRequest, DomainResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao processar requisição: %v", err),
		})
		return
	}

	log.Printf("[handleRemoveDomain] Requisição recebida: SiteID=%s, Domain=%s", req.SiteID, req.Domain)
	
	// Validar campos obrigatórios
	if req.SiteID == "" || req.Domain == "" {
		log.Printf("[handleRemoveDomain] Campos obrigatórios não informados")
		c.JSON(http.StatusBadRequest, DomainResponse{
			Success: false,
			Message: "ID do site e domínio são obrigatórios",
			SiteID:  req.SiteID,
			Domain:  req.Domain,
		})
		return
	}

	// Criar cliente Netlify
	log.Printf("[handleRemoveDomain] Criando cliente Netlify")
	netlifyClient, err := netlify.NewClient(s.config)
	if err != nil {
		log.Printf("[handleRemoveDomain] Erro ao criar cliente Netlify: %v", err)
		c.JSON(http.StatusInternalServerError, DomainResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao criar cliente Netlify: %v", err),
			SiteID:  req.SiteID,
			Domain:  req.Domain,
		})
		return
	}

	// Remover domínio
	log.Printf("[handleRemoveDomain] Removendo domínio %s do site %s", req.Domain, req.SiteID)
	
	// Criar contexto com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	err = netlifyClient.RemoveCustomDomain(ctx, req.SiteID, req.Domain)
	if err != nil {
		log.Printf("[handleRemoveDomain] Erro ao remover domínio: %v", err)
		c.JSON(http.StatusInternalServerError, DomainResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao remover domínio: %v", err),
			SiteID:  req.SiteID,
			Domain:  req.Domain,
		})
		return
	}

	log.Printf("[handleRemoveDomain] Domínio removido com sucesso")
	c.JSON(http.StatusOK, DomainResponse{
		Success: true,
		Message: "Domínio removido com sucesso",
		SiteID:  req.SiteID,
		Domain:  req.Domain,
	})
}

// handleSetDefaultDomain define um domínio como o domínio principal de um site
// @Summary Define um domínio como o domínio principal
// @Description Define um domínio personalizado como o domínio principal de um site existente na Netlify
// @Tags domains
// @Accept json
// @Produce json
// @Param request body DomainRequest true "Dados do domínio a ser definido como principal"
// @Success 200 {object} DomainResponse
// @Failure 400 {object} DomainResponse
// @Failure 500 {object} DomainResponse
// @Router /api/domains/set-default [post]
func (s *Server) handleSetDefaultDomain(c *gin.Context) {
	log.Printf("[handleSetDefaultDomain] Recebendo requisição para definir domínio padrão")
	
	var req DomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[handleSetDefaultDomain] Erro ao processar JSON: %v", err)
		c.JSON(http.StatusBadRequest, DomainResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao processar requisição: %v", err),
		})
		return
	}

	log.Printf("[handleSetDefaultDomain] Requisição recebida: SiteID=%s, Domain=%s", req.SiteID, req.Domain)
	
	// Validar campos obrigatórios
	if req.SiteID == "" || req.Domain == "" {
		log.Printf("[handleSetDefaultDomain] Campos obrigatórios não informados")
		c.JSON(http.StatusBadRequest, DomainResponse{
			Success: false,
			Message: "ID do site e domínio são obrigatórios",
			SiteID:  req.SiteID,
			Domain:  req.Domain,
		})
		return
	}

	// Criar cliente Netlify
	log.Printf("[handleSetDefaultDomain] Criando cliente Netlify")
	netlifyClient, err := netlify.NewClient(s.config)
	if err != nil {
		log.Printf("[handleSetDefaultDomain] Erro ao criar cliente Netlify: %v", err)
		c.JSON(http.StatusInternalServerError, DomainResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao criar cliente Netlify: %v", err),
			SiteID:  req.SiteID,
			Domain:  req.Domain,
		})
		return
	}

	// Definir domínio padrão
	log.Printf("[handleSetDefaultDomain] Definindo domínio %s como padrão para o site %s", req.Domain, req.SiteID)
	
	// Criar contexto com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	err = netlifyClient.SetDefaultDomain(ctx, req.SiteID, req.Domain, "")
	if err != nil {
		log.Printf("[handleSetDefaultDomain] Erro ao definir domínio padrão: %v", err)
		c.JSON(http.StatusInternalServerError, DomainResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao definir domínio padrão: %v", err),
			SiteID:  req.SiteID,
			Domain:  req.Domain,
		})
		return
	}

	log.Printf("[handleSetDefaultDomain] Domínio definido como padrão com sucesso")
	c.JSON(http.StatusOK, DomainResponse{
		Success: true,
		Message: "Domínio definido como padrão com sucesso",
		SiteID:  req.SiteID,
		Domain:  req.Domain,
	})
}

// handleRemovePrimaryDomain remove o domínio principal de um site
// @Summary Remove o domínio principal de um site
// @Description Remove o domínio principal de um site existente na Netlify, mantendo os aliases
// @Tags domains
// @Accept json
// @Produce json
// @Param request body DomainRequest true "ID do site a ter o domínio principal removido"
// @Success 200 {object} DomainResponse
// @Failure 400 {object} DomainResponse
// @Failure 500 {object} DomainResponse
// @Router /api/domains/remove-primary [post]
func (s *Server) handleRemovePrimaryDomain(c *gin.Context) {
	log.Printf("[handleRemovePrimaryDomain] Recebendo requisição para remover domínio principal")
	
	var req DomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[handleRemovePrimaryDomain] Erro ao processar JSON: %v", err)
		c.JSON(http.StatusBadRequest, DomainResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao processar requisição: %v", err),
		})
		return
	}

	log.Printf("[handleRemovePrimaryDomain] Requisição recebida: SiteID=%s", req.SiteID)
	
	// Validar campos obrigatórios
	if req.SiteID == "" {
		log.Printf("[handleRemovePrimaryDomain] Campo obrigatório não informado")
		c.JSON(http.StatusBadRequest, DomainResponse{
			Success: false,
			Message: "ID do site é obrigatório",
			SiteID:  req.SiteID,
		})
		return
	}

	// Criar cliente Netlify
	log.Printf("[handleRemovePrimaryDomain] Criando cliente Netlify")
	netlifyClient, err := netlify.NewClient(s.config)
	if err != nil {
		log.Printf("[handleRemovePrimaryDomain] Erro ao criar cliente Netlify: %v", err)
		c.JSON(http.StatusInternalServerError, DomainResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao criar cliente Netlify: %v", err),
			SiteID:  req.SiteID,
		})
		return
	}

	// Remover domínio principal
	log.Printf("[handleRemovePrimaryDomain] Removendo domínio principal do site %s", req.SiteID)
	
	// Criar contexto com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	err = netlifyClient.RemovePrimaryDomain(ctx, req.SiteID)
	if err != nil {
		log.Printf("[handleRemovePrimaryDomain] Erro ao remover domínio principal: %v", err)
		c.JSON(http.StatusInternalServerError, DomainResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao remover domínio principal: %v", err),
			SiteID:  req.SiteID,
		})
		return
	}

	log.Printf("[handleRemovePrimaryDomain] Domínio principal removido com sucesso")
	c.JSON(http.StatusOK, DomainResponse{
		Success: true,
		Message: "Domínio principal removido com sucesso",
		SiteID:  req.SiteID,
	})
}

// handleTestNetlifyConnection testa a conexão com a API da Netlify
// @Summary Testa a conexão com a API da Netlify
// @Description Testa a conexão com a API da Netlify e exibe informações sobre o token
// @Tags netlify
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/test/netlify/connection [get]
func (s *Server) handleTestNetlifyConnection(c *gin.Context) {
	log.Printf("[handleTestNetlifyConnection] Testando conexão com a API da Netlify")
	
	// Verificar variáveis de ambiente carregadas
	log.Printf("[handleTestNetlifyConnection] Variáveis de ambiente carregadas:")
	log.Printf("NETLIFY_TOKEN presente: %v", s.config.NetlifyToken != "")
	log.Printf("BASE_DOMAIN: %s", s.config.BaseDomain)
	
	// Exibir valor original das variáveis de ambiente
	originalToken := os.Getenv("NETLIFY_TOKEN")
	maskTokenValue := "<vazio>"
	if originalToken != "" {
		if len(originalToken) > 10 {
			maskTokenValue = originalToken[:10] + "..."
		} else {
			maskTokenValue = "<token muito curto>"
		}
	}
	
	log.Printf("[handleTestNetlifyConnection] Token original do ambiente: %s", maskTokenValue)
	
	// Recarregar .env para depuração
	log.Printf("[handleTestNetlifyConnection] Tentando recarregar arquivo .env")
	if err := godotenv.Load(); err != nil {
		log.Printf("[handleTestNetlifyConnection] Erro ao carregar arquivo .env: %v", err)
	} else {
		log.Printf("[handleTestNetlifyConnection] Arquivo .env carregado com sucesso")
		tokenFromEnv := os.Getenv("NETLIFY_TOKEN")
		envTokenMasked := "<vazio>"
		if tokenFromEnv != "" {
			if len(tokenFromEnv) > 10 {
				envTokenMasked = tokenFromEnv[:10] + "..."
			} else {
				envTokenMasked = "<token muito curto>"
			}
		}
		log.Printf("[handleTestNetlifyConnection] Token do .env: %s", envTokenMasked)
	}
	
	// Verificar todas as variáveis de ambiente
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "NETLIFY_") {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				key := parts[0]
				value := parts[1]
				maskedValue := "<vazio>"
				if value != "" {
					if len(value) > 10 {
						maskedValue = value[:10] + "..."
					} else {
						maskedValue = "<valor muito curto>"
					}
				}
				log.Printf("[handleTestNetlifyConnection] Variável encontrada: %s = %s", key, maskedValue)
			}
		}
	}
	
	// Retornar informações sobre as variáveis de ambiente
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Informações sobre as variáveis de ambiente",
		"environment": gin.H{
			"netlify_token": gin.H{
				"config_value": s.config.NetlifyToken != "",
				"env_value": originalToken != "",
				"env_reloaded": os.Getenv("NETLIFY_TOKEN") != "",
			},
			"base_domain": s.config.BaseDomain,
		},
	})
}

// handleViewLogs exibe os logs da aplicação
func (s *Server) handleViewLogs(c *gin.Context) {
	log.Printf("[handleViewLogs] Exibindo logs da aplicação")
	
	// Verificar se o arquivo de logs existe
	logFilePath := "netlify-deploy.log"
	logContent, err := os.ReadFile(logFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Erro ao ler arquivo de logs: %v", err),
		})
		return
	}
	
	// Converter o conteúdo do arquivo para string
	logText := string(logContent)
	
	// Retornar os logs formatados como JSON
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"logs": logText,
	})
}

// handleListSites lista todos os sites do usuário na Netlify
func (s *Server) handleListSites(c *gin.Context) {
	log.Printf("[API] Recebida requisição para listar sites")
	
	// Configurar o cliente da Netlify
	netlifyClient, err := netlify.NewClient(s.config)
	if err != nil {
		log.Printf("[API] Erro ao criar cliente Netlify: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Erro interno do servidor: %v", err),
		})
		return
	}
	
	// Listar sites
	sites, err := netlifyClient.ListSites(context.Background())
	if err != nil {
		log.Printf("[API] Erro ao listar sites: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Erro ao listar sites: %v", err),
		})
		return
	}
	
	// Preparar resposta
	type SiteInfo struct {
		ID            string   `json:"id"`
		Name          string   `json:"name"`
		URL           string   `json:"url"`
		SslURL        string   `json:"ssl_url"`
		CustomDomain  string   `json:"custom_domain,omitempty"`
		DomainAliases []string `json:"domain_aliases,omitempty"`
		CreatedAt     string    `json:"created_at"`
	}
	
	var sitesList []SiteInfo
	
	for _, site := range sites {
		createdAt := ""
		
		// Tratar a data de criação de forma segura
		// O tipo exato pode variar dependendo da implementação da API
		if reflect.ValueOf(site.CreatedAt).IsValid() && !reflect.ValueOf(site.CreatedAt).IsZero() {
			createdAt = fmt.Sprintf("%v", site.CreatedAt)
		}
		
		siteInfo := SiteInfo{
			ID:            site.ID,
			Name:          site.Name,
			URL:           site.URL,
			SslURL:        site.SslURL,
			CustomDomain:  site.CustomDomain,
			DomainAliases: site.DomainAliases,
			CreatedAt:     createdAt,
		}
		
		sitesList = append(sitesList, siteInfo)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("Encontrados %d sites", len(sites)),
		"sites": sitesList,
	})
}

// maskToken oculta parte do token para exibição segura
func maskToken(token string) string {
	if token == "" {
		return "<token vazio>"
	}
	
	if len(token) <= 10 {
		return "<token muito curto>"
	}
	
	return token[:10] + "..."
}

// processDeploy realiza o processo de deploy em segundo plano
func (s *Server) processDeploy(req DeployRequest) {
	ctx := context.Background()
	log.Printf("Iniciando deploy para usuário: %s, caminho S3: %s", req.Username, req.S3Path)

	// Inicializar cliente Netlify
	netlifyClient, err := netlify.NewClient(s.config)
	if err != nil {
		log.Printf("Erro ao inicializar cliente Netlify: %v", err)
		return
	}

	// Verificar se o site já existe
	site, exists, err := netlifyClient.VerifySite(ctx)
	if err != nil {
		log.Printf("Erro ao verificar site: %v", err)
		return
	}

	// Criar site se não existir
	if !exists {
		site, err = netlifyClient.CreateSite(ctx)
		if err != nil {
			log.Printf("Erro ao criar site: %v", err)
			return
		}
	}

	// Configurar DNS
	if err := netlifyClient.ConfigureDNS(ctx, site); err != nil {
		log.Printf("Aviso: erro ao configurar DNS: %v", err)
	}

	// Inicializar cliente S3
	s3Client, err := aws.NewS3Client(s.config)
	if err != nil {
		log.Printf("Erro ao inicializar cliente S3: %v", err)
		return
	}

	// Criar diretório temporário para os arquivos
	tempDir, err := os.MkdirTemp("", "netlify-deploy-*")
	if err != nil {
		log.Printf("Erro ao criar diretório temporário: %v", err)
		return
	}
	defer os.RemoveAll(tempDir) // Limpar diretório temporário ao finalizar

	// Baixar arquivos do S3
	if err := s3Client.DownloadFiles(ctx, tempDir); err != nil {
		log.Printf("Erro ao baixar arquivos do S3: %v", err)
		return
	}

	// Realizar deploy dos arquivos
	deploy, err := netlifyClient.DeploySite(ctx, site, tempDir)
	if err != nil {
		log.Printf("Erro ao iniciar deploy: %v", err)
		return
	}

	log.Printf("Deploy iniciado com sucesso: ID %s", deploy.ID)

	// Aguardar conclusão do deploy
	finalDeploy, err := netlifyClient.WaitForDeploy(ctx, deploy.ID)
	if err != nil {
		log.Printf("Erro ao aguardar deploy: %v", err)
	} else {
		log.Printf("Deploy concluído com sucesso! Site disponível em: %s", finalDeploy.URL)
		log.Printf("URL do subdomínio: https://%s", s.config.NetlifySubdomain)
		if s.config.CustomDomain != "" {
			log.Printf("URL do domínio personalizado: https://%s", s.config.CustomDomain)
		}
	}
}

// handleDeployFromS3 processa uma requisição de deploy a partir do S3 (rota desativada)
func (s *Server) handleDeployFromS3(c *gin.Context) {
	log.Printf("[API] Recebida requisição para deploy a partir do S3: %s", c.Request.URL.Path)
	
	// Obter parâmetros do formulário
	siteID := c.PostForm("site_id")
	siteName := c.PostForm("site_name")
	if siteName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Nome do site é obrigatório",
		})
		return
	}

	s3Path := c.PostForm("s3_path")
	if s3Path == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Caminho no bucket S3 é obrigatório",
		})
		return
	}

	customDomain := c.PostForm("custom_domain")
	// Variável description obtida mas não utilizada, comentada para evitar erro
	// description := c.PostForm("description")

	log.Printf("[API] Parâmetros de deploy do S3: siteID=%s, siteName=%s, s3Path=%s, customDomain=%s", 
		siteID, siteName, s3Path, customDomain)

	// Configurar o cliente da Netlify
	netlifyClient, err := netlify.NewClient(s.config)
	if err != nil {
		log.Printf("[API] Erro ao criar cliente Netlify: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Erro interno do servidor: %v", err),
		})
		return
	}
	
	// Configurar parâmetros de deploy no objeto config
	// Usamos o siteName como username para manter a consistência
	if err := s.config.SetDeployParams(siteName, customDomain, s3Path); err != nil {
		log.Printf("[API] Erro ao configurar parâmetros de deploy: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("Erro nos parâmetros de deploy: %v", err),
		})
		return
	}

	// Configurar cliente S3
	s3Client, err := aws.NewS3Client(s.config)
	if err != nil {
		log.Printf("[API] Erro ao criar cliente S3: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Erro ao criar cliente S3: %v", err),
		})
		return
	}

	// Criar diretório temporário para os arquivos
	tempDir, err := os.MkdirTemp("", "netlify-s3-deploy-*")
	if err != nil {
		log.Printf("[API] Erro ao criar diretório temporário: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Erro ao criar diretório temporário: %v", err),
		})
		return
	}
	// Adicionar defer para limpar o diretório temporário ao final
	defer os.RemoveAll(tempDir)

	log.Printf("[API] Baixando arquivos do S3 para diretório temporário: %s", tempDir)
	
	// Baixar arquivos do S3
	if err := s3Client.DownloadFiles(c.Request.Context(), tempDir); err != nil {
		log.Printf("[API] Erro ao baixar arquivos do S3: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Erro ao baixar arquivos do S3: %v", err),
		})
		return
	}

	// Verificar se existem arquivos no diretório
	files, err := os.ReadDir(tempDir)
	if err != nil {
		log.Printf("[API] Erro ao ler diretório temporário: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Erro ao verificar arquivos baixados: %v", err),
		})
		return
	}

	if len(files) == 0 {
		log.Printf("[API] Nenhum arquivo encontrado no caminho S3 especificado: %s", s3Path)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("Nenhum arquivo encontrado no caminho S3 especificado: %s", s3Path),
		})
		return
	}

	log.Printf("[API] %d arquivos/diretórios baixados do S3", len(files))

	// Processar o deploy na Netlify
	var site *models.Site
	var exists bool
	
	// Se o ID do site foi fornecido, usamos ele
	if siteID != "" {
		// Verificar se o site existe com o ID fornecido
		site, exists, err = netlifyClient.VerifySiteById(c.Request.Context(), siteID)
		if err != nil {
			log.Printf("[API] Erro ao verificar site na Netlify com ID %s: %v", siteID, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": fmt.Sprintf("Erro ao verificar site na Netlify: %v", err),
			})
			return
		}
		
		if !exists {
			log.Printf("[API] Site com ID %s não encontrado", siteID)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fmt.Sprintf("Site com ID %s não encontrado", siteID),
			})
			return
		}
	} else {
		// Caso contrário, verificamos se existe um site com o nome fornecido
		site, exists, err = netlifyClient.VerifySite(c.Request.Context())
		if err != nil {
			log.Printf("[API] Erro ao verificar site na Netlify: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": fmt.Sprintf("Erro ao verificar site na Netlify: %v", err),
			})
			return
		}

		// Se o site não existir, criá-lo
		if !exists {
			log.Printf("[API] Site não encontrado. Criando novo site: %s", siteName)
			site, err = netlifyClient.CreateSite(c.Request.Context())
			if err != nil {
				log.Printf("[API] Erro ao criar site na Netlify: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": fmt.Sprintf("Erro ao criar site na Netlify: %v", err),
				})
				return
			}
			log.Printf("[API] Site criado com sucesso: %s (ID: %s)", site.Name, site.ID)
		} else {
			log.Printf("[API] Site existente encontrado: %s (ID: %s)", site.Name, site.ID)
		}
	}

	// Realizar deploy dos arquivos baixados do S3
	log.Printf("[API] Iniciando deploy dos arquivos para o site %s (ID: %s)", site.Name, site.ID)
	deploy, err := netlifyClient.DeploySite(c.Request.Context(), site, tempDir)
	if err != nil {
		log.Printf("[API] Erro ao iniciar deploy: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Erro ao iniciar deploy: %v", err),
		})
		return
	}

	log.Printf("[API] Deploy iniciado com sucesso: ID %s", deploy.ID)

	// Aguardar a conclusão do deploy se o cliente ainda estiver conectado
	finalDeploy, err := netlifyClient.WaitForDeploy(c.Request.Context(), deploy.ID)
	if err != nil {
		log.Printf("[API] Erro ao aguardar deploy: %v", err)
		// Mesmo com erro, retornamos sucesso parcial pois o deploy foi iniciado
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": fmt.Sprintf("Deploy iniciado, mas não foi possível aguardar a conclusão: %v", err),
			"site_id": site.ID,
			"site_url": site.URL,
			"deploy_id": deploy.ID,
			"test_success": false,
		})
		return
	}

	// Deploy concluído com sucesso
	log.Printf("[API] Deploy concluído com sucesso! Site disponível em: %s", finalDeploy.URL)
	
	// Retornar resposta de sucesso
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Deploy concluído com sucesso",
		"site_id": site.ID,
		"site_url": finalDeploy.URL,
		"deploy_id": deploy.ID,
		"created_at": time.Now().Format(time.RFC3339),
		"test_success": true,
	})
}

// Start inicia o servidor da API
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.config.APIPort)
	log.Printf("Iniciando servidor na porta %s", s.config.APIPort)
	return s.router.Run(addr)
}
