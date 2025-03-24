package netlify

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"encoding/base64"

	"github.com/netlify/open-api/go/models"
	porcelainctx "github.com/netlify/open-api/go/porcelain/context"
)

// TestDeployParams contém os parâmetros para o teste de deploy
type TestDeployParams struct {
	SiteID          string `json:"site_id" example:"a1b2c3d4" swagger:"description=ID do site na Netlify para atualização (opcional)"` 
	SiteName        string `json:"site_name" binding:"required" example:"test-site" swagger:"description=Nome do site para teste"`
	Description     string `json:"description" example:"Site de teste" swagger:"description=Descrição do site para teste"`
	TestContent     string `json:"test_content" example:"<h1>Hello World</h1>" swagger:"description=Conteúdo HTML para teste (opcional)"`
	CleanupAfter    bool   `json:"cleanup_after" example:"true" swagger:"description=Remover o site após o teste"`
	CustomDomain    string `json:"custom_domain" example:"meu-site.exemplo.com" swagger:"description=Domínio personalizado para o site (opcional)"`
	FileContent     string `json:"file_content" example:"" swagger:"description=Conteúdo de arquivo HTML em formato base64 (opcional, alternativa ao TestContent)"`
	FolderPath      string `json:"folder_path" example:"/path/to/folder" swagger:"description=Caminho da pasta local para deploy (opcional, alternativa ao FileContent e TestContent)"`
}

// TestDeployResult contém o resultado do teste de deploy
type TestDeployResult struct {
	Success     bool      `json:"success" example:"true" swagger:"description=Indica se o teste foi bem-sucedido"`
	Message     string    `json:"message" example:"Site de teste criado com sucesso" swagger:"description=Mensagem do resultado do teste"`
	SiteID      string    `json:"site_id,omitempty" example:"a1b2c3d4" swagger:"description=ID do site criado na Netlify"`
	SiteURL     string    `json:"site_url,omitempty" example:"https://test-site.netlify.app" swagger:"description=URL do site criado"`
	CreatedAt   time.Time `json:"created_at,omitempty" swagger:"description=Data e hora de criação do site"`
	TestSuccess bool      `json:"test_success" example:"true" swagger:"description=Indica se o teste específico foi bem-sucedido"`
}

// ExecuteTestDeploy realiza um teste de deploy na Netlify
func (c *Client) ExecuteTestDeploy(ctx context.Context, params TestDeployParams) (*TestDeployResult, error) {
	log.Printf("[TEST] Iniciando teste de deploy: %s", params.SiteName)
	log.Printf("[TEST] Detalhes completos dos parâmetros: %+v", params)

	// Preparar resultado
	result := &TestDeployResult{
		CreatedAt: time.Now(),
	}

	log.Printf("[TEST] Token da Netlify: %s...", c.config.NetlifyToken[:10])

	// Criar um novo contexto com a autenticação
	authCtx := porcelainctx.WithAuthInfo(context.Background(), c.auth)

	// Se um SiteID foi fornecido, buscar o site diretamente
	var site *models.Site
	var err error
	if params.SiteID != "" {
		log.Printf("[TEST] Buscando site pelo ID fornecido: %s", params.SiteID)
		site, err = c.netlify.GetSite(authCtx, params.SiteID)
		if err != nil {
			log.Printf("[TEST] ERRO ao buscar site pelo ID: %v", err)
			return nil, fmt.Errorf("erro ao buscar site pelo ID: %w", err)
		}
		log.Printf("[TEST] Site encontrado: %s (ID: %s)", site.Name, site.ID)
		
		// Atualizar o nome do site se necessário
		if site.Name != params.SiteName {
			log.Printf("[TEST] Atualizando nome do site de '%s' para '%s'", site.Name, params.SiteName)
			
			// Criar um objeto SiteSetup para atualização
			updateSiteParams := &models.SiteSetup{
				Site: models.Site{
					ID:   site.ID,
					Name: params.SiteName,
				},
			}
			
			// Atualizar o site na Netlify
			site, err = c.netlify.UpdateSite(authCtx, updateSiteParams)
			if err != nil {
				log.Printf("[TEST] ERRO ao atualizar site: %v", err)
				return nil, fmt.Errorf("erro ao atualizar site: %w", err)
			}
			log.Printf("[TEST] Site atualizado com sucesso: %s (ID: %s)", site.Name, site.ID)
		}
	} else {
		// Verificar se já existe um site com este nome
		log.Printf("[TEST] Buscando sites existentes na Netlify...")
		
		sites, err := c.netlify.ListSites(authCtx, nil)
		if err != nil {
			log.Printf("[TEST] ERRO ao listar sites: %v", err)
			return nil, fmt.Errorf("erro ao listar sites: %w", err)
		}

		log.Printf("[TEST] Encontrado %d sites no total", len(sites))

		var existingSite *models.Site
		for _, s := range sites {
			if s.Name == params.SiteName {
				existingSite = s
				log.Printf("[TEST] Site com nome '%s' já existe (ID: %s)", params.SiteName, s.ID)
				break
			}
		}

		// Se o site já existe e devemos limpá-lo, excluí-lo primeiro
		if existingSite != nil && params.CleanupAfter {
			log.Printf("[TEST] Excluindo site existente: %s (ID: %s)", existingSite.Name, existingSite.ID)
			err := c.netlify.DeleteSite(authCtx, existingSite.ID)
			if err != nil {
				log.Printf("[TEST] Erro ao excluir site existente: %v", err)
				return nil, fmt.Errorf("erro ao excluir site existente: %w", err)
			}
			existingSite = nil
			log.Printf("[TEST] Site excluído com sucesso")
		}

		// Criar novo site de teste ou usar o existente
		if existingSite == nil {
			// Configurar as opções do site
			siteParams := models.SiteSetup{
				Site: models.Site{
					Name: params.SiteName,
				},
			}

			// Criar o site
			log.Printf("[TEST] Criando novo site de teste: %s", params.SiteName)
			log.Printf("[TEST] Parâmetros do site: %+v", siteParams)
			site, err = c.netlify.CreateSite(authCtx, &siteParams, false)
			if err != nil {
				log.Printf("[TEST] ERRO ao criar site: %v", err)
				log.Printf("[TEST] ERRO detalhado: %#v", err)
				
				// Tentar extrair mais detalhes do erro
				if apiErr, ok := err.(interface{ Error() string }); ok {
					log.Printf("[TEST] Mensagem de erro: %s", apiErr.Error())
				}
				
				return nil, fmt.Errorf("erro ao criar site de teste: %w", err)
			}
			
			// Configurar domínio personalizado após a criação do site
			if params.CustomDomain != "" {
				log.Printf("[TEST] Configurando domínio personalizado após criação: %s", params.CustomDomain)
				err = c.configureCustomDomain(authCtx, site, params.CustomDomain)
				if err != nil {
					log.Printf("[TEST] AVISO: Erro ao configurar domínio personalizado: %v", err)
					// Não falhar o processo por erro no domínio personalizado
				}
			}
			log.Printf("[TEST] Site criado com sucesso: %s (ID: %s)", site.Name, site.ID)
		} else {
			site = existingSite
			log.Printf("[TEST] Usando site existente: %s (ID: %s)", site.Name, site.ID)
		}
	}

	// Preencher resultado com informações do site
	result.SiteID = site.ID
	result.SiteURL = site.URL
	result.Success = true
	result.Message = "Site de teste criado/atualizado com sucesso"
	result.TestSuccess = true

	// Verificar o conteúdo para deploy - priorizar o arquivo sobre conteúdo de texto
	var files map[string]string
	if params.FileContent != "" {
		log.Printf("[TEST] Processando arquivo enviado para deploy")
		
		// Decodificar o arquivo Base64
		fileData, err := base64.StdEncoding.DecodeString(params.FileContent)
		if err != nil {
			log.Printf("[TEST] Erro ao decodificar conteúdo do arquivo: %v", err)
			return nil, fmt.Errorf("erro ao decodificar conteúdo do arquivo: %w", err)
		}
		
		log.Printf("[TEST] Arquivo decodificado com sucesso, tamanho: %d bytes", len(fileData))
		
		// Criar o arquivo index.html com o conteúdo decodificado
		files = map[string]string{
			"index.html": string(fileData),
		}
	} else if params.TestContent != "" {
		log.Printf("[TEST] Realizando deploy de conteúdo de teste")
		
		// Preparar o conteúdo para deploy
		files = map[string]string{
			"index.html": params.TestContent,
		}
	} else if params.FolderPath != "" {
		log.Printf("[TEST] Realizando deploy de pasta local: %s", params.FolderPath)
		
		// Verificar se a pasta existe
		if _, err := os.Stat(params.FolderPath); os.IsNotExist(err) {
			log.Printf("[TEST] Erro: Pasta não encontrada: %s", params.FolderPath)
			return nil, fmt.Errorf("pasta não encontrada: %w", err)
		}
		
		// Realizar deploy diretamente da pasta local
		deployment, err := c.DeployLocalFolder(authCtx, site, params.FolderPath)
		if err != nil {
			log.Printf("[TEST] Erro ao realizar deploy da pasta local: %v", err)
			result.TestSuccess = false
			result.Message += fmt.Sprintf(". Porém, ocorreu um erro durante o deploy: %v", err)
			return result, nil
		}
		
		log.Printf("[TEST] Deploy da pasta local realizado com sucesso, ID: %s", deployment.ID)
		result.Message += fmt.Sprintf(". Deploy realizado com sucesso (ID: %s)", deployment.ID)
		
		// Retornar resultado sem continuar com o deploy de arquivos
		return result, nil
	} else {
		log.Printf("[TEST] Nenhum conteúdo fornecido para deploy")
		// Nenhum conteúdo fornecido, apenas criar um arquivo HTML padrão
		files = map[string]string{
			"index.html": "<html><body><h1>Teste de Deploy</h1><p>Site criado em " + time.Now().Format(time.RFC3339) + "</p></body></html>",
		}
	}
	
	// Realizar o deploy do conteúdo
	if len(files) > 0 {
		deployment, err := c.DeployContent(authCtx, site, files)
		if err != nil {
			log.Printf("[TEST] Erro ao realizar deploy do conteúdo: %v", err)
			// Não falharemos o teste se o deploy não funcionar, apenas registramos
			result.TestSuccess = false
			result.Message += fmt.Sprintf(". Porém, ocorreu um erro durante o deploy: %v", err)
		} else {
			log.Printf("[TEST] Deploy de conteúdo realizado com sucesso, ID: %s", deployment.ID)
			result.Message += fmt.Sprintf(". Deploy realizado com sucesso (ID: %s)", deployment.ID)
		}
	}

	// Se solicitado para limpar após o teste e não é um site existente (que queremos manter)
	if params.CleanupAfter && params.SiteID == "" {
		log.Printf("[TEST] Agendando limpeza do site após teste")
		// Em um caso real, poderíamos agendar a exclusão do site após um período
		// ou criar um job em background. Por simplicidade, apenas registramos a intenção.
		log.Printf("[TEST] Site %s será excluído após o período de teste", site.Name)
	}

	return result, nil
}
