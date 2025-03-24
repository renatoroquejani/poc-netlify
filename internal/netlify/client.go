package netlify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/kodestech/poc-netlify/internal/config"
	"github.com/netlify/open-api/go/models"
	"github.com/netlify/open-api/go/porcelain"
)

// Client encapsula a integração com a API da Netlify
type Client struct {
	netlify *porcelain.Netlify
	config  *config.Config
	auth    runtime.ClientAuthInfoWriter
}

// NewClient cria um novo cliente Netlify
func NewClient(cfg *config.Config) (*Client, error) {
	// Verificar se o token está configurado
	if cfg.NetlifyToken == "" {
		log.Printf("ERRO: Token da Netlify não configurado")
		return nil, fmt.Errorf("token da Netlify não configurado")
	}

	log.Printf("Criando cliente Netlify com token: %s...", cfg.NetlifyToken[:10])

	// Configurar cliente HTTP com autenticação
	transport := client.New("api.netlify.com", "/api/v1", []string{"https"})

	// Criar a autenticação
	auth := runtime.ClientAuthInfoWriterFunc(func(req runtime.ClientRequest, reg strfmt.Registry) error {
		req.SetHeaderParam("Authorization", "Bearer "+cfg.NetlifyToken)
		return nil
	})

	transport.DefaultAuthentication = auth

	// Criar cliente Netlify
	netlifyClient := porcelain.New(transport, strfmt.Default)

	return &Client{
		netlify: netlifyClient,
		config:  cfg,
		auth:    auth,
	}, nil
}

// VerifySiteById verifica se um site existe na Netlify pelo ID
func (c *Client) VerifySiteById(ctx context.Context, siteID string) (*models.Site, bool, error) {
	log.Printf("Verificando se o site com ID %s existe", siteID)

	// Criar um contexto com autenticação
	authCtx := c.createAuthContext(ctx)

	// Obter o site pelo ID
	site, err := c.netlify.GetSite(authCtx, siteID)
	if err != nil {
		if err.Error() == "site not found" {
			log.Printf("Site com ID %s não encontrado", siteID)
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("erro ao verificar site: %w", err)
	}

	if site == nil {
		log.Printf("Site com ID %s não encontrado", siteID)
		return nil, false, nil
	}

	log.Printf("Site encontrado: %s (ID: %s)", site.Name, site.ID)
	return site, true, nil
}

// VerifySite verifica se o site já existe na Netlify
func (c *Client) VerifySite(ctx context.Context) (*models.Site, bool, error) {
	log.Printf("Verificando se o site %s já existe", c.config.NetlifySubdomain)

	// Listar todos os sites do usuário
	sites, err := c.ListSites(ctx)
	if err != nil {
		return nil, false, fmt.Errorf("erro ao listar sites: %w", err)
	}

	// Verificar se algum site corresponde ao subdomínio desejado
	for _, site := range sites {
		// Verificar pelo domínio personalizado
		if site.CustomDomain == c.config.NetlifySubdomain {
			log.Printf("Site encontrado: %s (ID: %s)", site.Name, site.ID)
			return site, true, nil
		}

		// Verificar pelos aliases de domínio
		for _, domain := range site.DomainAliases {
			if domain == c.config.NetlifySubdomain {
				log.Printf("Site encontrado: %s (ID: %s)", site.Name, site.ID)
				return site, true, nil
			}
		}

		// Verificar pelo subdomínio padrão da Netlify
		if site.Name+".netlify.app" == c.config.NetlifySubdomain {
			log.Printf("Site encontrado: %s (ID: %s)", site.Name, site.ID)
			return site, true, nil
		}
	}

	log.Printf("Site %s não encontrado", c.config.NetlifySubdomain)
	return nil, false, nil
}

// CreateSite cria um novo site na Netlify
func (c *Client) CreateSite(ctx context.Context) (*models.Site, error) {
	log.Printf("Criando novo site: %s", c.config.NetlifySubdomain)

	// Configurar as opções do site
	siteParams := models.SiteSetup{
		Site: models.Site{
			Name:         c.config.Username,
			CustomDomain: c.config.NetlifySubdomain,
			DomainAliases: []string{
				c.config.NetlifySubdomain,
			},
		},
	}

	// Criar o site
	site, err := c.netlify.CreateSite(ctx, &siteParams, false)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar site: %w", err)
	}

	log.Printf("Site criado com sucesso: %s (ID: %s)", site.Name, site.ID)
	return site, nil
}

// ConfigureDNS configura o DNS para o site
func (c *Client) ConfigureDNS(ctx context.Context, site *models.Site) error {
	log.Printf("Configurando DNS para o site %s", site.Name)

	// Configurar o domínio personalizado padrão (subdomínio)
	if err := c.configureCustomDomain(ctx, site, c.config.NetlifySubdomain); err != nil {
		return fmt.Errorf("erro ao configurar subdomínio padrão: %w", err)
	}

	// Se um domínio personalizado foi fornecido, configurá-lo também
	if c.config.CustomDomain != "" {
		log.Printf("Configurando domínio personalizado: %s", c.config.CustomDomain)
		if err := c.configureCustomDomain(ctx, site, c.config.CustomDomain); err != nil {
			return fmt.Errorf("erro ao configurar domínio personalizado: %w", err)
		}
	}

	return nil
}

// configureCustomDomain configura um domínio personalizado para o site
func (c *Client) configureCustomDomain(ctx context.Context, site *models.Site, domain string) error {
	// Verificar se o domínio já está configurado
	if site.CustomDomain == domain {
		log.Printf("Domínio %s já está configurado como domínio principal", domain)
		return nil
	}

	// Verificar se o domínio já está nos aliases
	domainFound := false
	for _, d := range site.DomainAliases {
		if d == domain {
			domainFound = true
			break
		}
	}

	if !domainFound {
		// Adicionar o domínio como alias
		site.DomainAliases = append(site.DomainAliases, domain)
		log.Printf("Adicionando %s como alias de domínio", domain)
	} else {
		log.Printf("Domínio %s já está configurado como alias", domain)
		return nil
	}

	// Atualizar o site
	siteSetup := &models.SiteSetup{
		Site: *site,
	}

	// Utilizar o contexto que foi passado, que já contém a autenticação
	updatedSite, err := c.netlify.UpdateSite(ctx, siteSetup)
	if err != nil {
		return fmt.Errorf("erro ao atualizar site com domínio personalizado: %w", err)
	}

	log.Printf("Domínio %s configurado com sucesso para o site %s", domain, updatedSite.Name)

	// Exibir mensagem informativa sobre DNS
	log.Printf("Para configurar o DNS para %s, adicione um registro CNAME apontando para o domínio Netlify", domain)

	return nil
}

// DeploySite realiza o deploy dos arquivos para o site
func (c *Client) DeploySite(ctx context.Context, site *models.Site, deployDir string) (*models.Deploy, error) {
	log.Printf("Iniciando deploy para o site %s a partir do diretório %s", site.Name, deployDir)

	// Configurar opções de deploy
	deployOptions := porcelain.DeployOptions{
		SiteID:  site.ID,
		Dir:     deployDir,
		IsDraft: false,
		Title:   fmt.Sprintf("Deploy automático para %s", c.config.NetlifySubdomain),
	}

	// Realizar o deploy
	deploy, err := c.netlify.DeploySite(ctx, deployOptions)
	if err != nil {
		return nil, fmt.Errorf("erro ao realizar deploy: %w", err)
	}

	log.Printf("Deploy iniciado com sucesso: ID %s", deploy.ID)
	return deploy, nil
}

// DeployContent realiza o deploy de conteúdo para o site
func (c *Client) DeployContent(ctx context.Context, site *models.Site, files map[string]string) (*models.Deploy, error) {
	log.Printf("Iniciando deploy de conteúdo para o site %s", site.Name)

	// Criar um diretório temporário
	tmpDir, err := c.createTempFilesFromContent(files)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar arquivos temporários: %w", err)
	}

	// Configurar opções de deploy
	deployOptions := porcelain.DeployOptions{
		SiteID:  site.ID,
		Dir:     tmpDir,
		IsDraft: false,
		Title:   fmt.Sprintf("Deploy de conteúdo para %s", site.Name),
	}

	// Realizar o deploy
	deploy, err := c.netlify.DeploySite(ctx, deployOptions)
	if err != nil {
		return nil, fmt.Errorf("erro ao realizar deploy de conteúdo: %w", err)
	}

	log.Printf("Deploy de conteúdo iniciado com sucesso: ID %s", deploy.ID)
	return deploy, nil
}

func (c *Client) createTempFilesFromContent(files map[string]string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "netlify-deploy")
	if err != nil {
		return "", err
	}

	for filename, content := range files {
		filePath := filepath.Join(tmpDir, filename)
		err = os.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			return "", err
		}
	}

	return tmpDir, nil
}

// DeployLocalFolder realiza o deploy de uma pasta local para o site
func (c *Client) DeployLocalFolder(ctx context.Context, site *models.Site, folderPath string) (*models.Deploy, error) {
	log.Printf("Iniciando deploy da pasta local %s para o site %s", folderPath, site.Name)

	// Configurar opções de deploy
	deployOptions := porcelain.DeployOptions{
		SiteID:  site.ID,
		Dir:     folderPath,
		IsDraft: false,
		Title:   fmt.Sprintf("Deploy da pasta %s para %s", folderPath, site.Name),
	}

	// Realizar o deploy
	deploy, err := c.netlify.DeploySite(ctx, deployOptions)
	if err != nil {
		return nil, fmt.Errorf("erro ao realizar deploy da pasta local: %w", err)
	}

	log.Printf("Deploy da pasta local iniciado com sucesso: ID %s", deploy.ID)
	return deploy, nil
}

// CreateOrGetSite cria um novo site ou retorna um existente na Netlify
func (c *Client) CreateOrGetSite(ctx context.Context, siteName, description string) (*models.Site, error) {
	log.Printf("Verificando se o site %s já existe", siteName)

	// Listar todos os sites do usuário
	sites, err := c.ListSites(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar sites: %w", err)
	}

	// Verificar se algum site corresponde ao nome desejado
	for _, site := range sites {
		if site.Name == siteName {
			log.Printf("Site encontrado: %s (ID: %s)", site.Name, site.ID)
			return site, nil
		}
	}

	// Configurar as opções do site
	siteParams := models.SiteSetup{
		Site: models.Site{
			Name: siteName,
		},
	}

	// Criar o site
	log.Printf("Criando novo site: %s", siteName)
	site, err := c.netlify.CreateSite(ctx, &siteParams, false)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar site: %w", err)
	}

	log.Printf("Site criado com sucesso: %s (ID: %s)", site.Name, site.ID)
	return site, nil
}

// WaitForDeploy aguarda a conclusão do deploy
func (c *Client) WaitForDeploy(ctx context.Context, deployID string) (*models.Deploy, error) {
	log.Printf("Aguardando conclusão do deploy %s", deployID)

	// Aguardar até que o deploy esteja pronto
	finished := false
	var deploy *models.Deploy
	var err error

	for !finished {
		log.Printf("Verificando status do deploy %s...", deployID)

		// Obter status do deploy
		deploy, err = c.netlify.GetDeploy(ctx, deployID)
		if err != nil {
			return nil, fmt.Errorf("erro ao verificar status do deploy: %w", err)
		}

		// Verificar se o deploy foi concluído
		if deploy.State == "ready" {
			finished = true
			log.Printf("Deploy concluído com sucesso: %s", deploy.URL)
		} else if deploy.State == "error" {
			return nil, fmt.Errorf("erro no deploy: %s", deploy.ErrorMessage)
		} else {
			log.Printf("Deploy ainda em progresso (estado: %s). Aguardando 2 segundos...", deploy.State)
			time.Sleep(2 * time.Second)
		}
	}

	return deploy, nil
}

// ListSites lista todos os sites do usuário
func (c *Client) ListSites(ctx context.Context) ([]*models.Site, error) {
	log.Printf("Listando sites do usuário...")

	// Criar um contexto com autenticação
	authCtx := c.createAuthContext(ctx)

	// Listar sites
	sites, err := c.netlify.ListSites(authCtx, nil)
	if err != nil {
		log.Printf("Erro ao listar sites: %v", err)
		return nil, fmt.Errorf("erro ao listar sites: %w", err)
	}

	log.Printf("Encontrados %d sites", len(sites))
	return sites, nil
}

// AddCustomDomain adiciona um domínio personalizado a um site.
// Se o site não possuir domínio principal, o domínio recebido será definido como principal (com validação TXT).
// Se já houver domínio principal, o domínio recebido será adicionado como alias, desde que não exista.
func (c *Client) AddCustomDomain(ctx context.Context, siteID, domain string) error {
	log.Printf("Iniciando adição de domínio %s para o site %s", domain, siteID)

	if siteID == "" {
		return fmt.Errorf("ID do site não pode ser vazio")
	}
	if domain == "" {
		return fmt.Errorf("domínio não pode ser vazio")
	}

	authCtx := c.createAuthContext(ctx)
	site, err := c.netlify.GetSite(authCtx, siteID)
	if err != nil {
		log.Printf("Erro ao obter site %s: %v", siteID, err)
		return fmt.Errorf("erro ao obter site: %w", err)
	}
	if site == nil {
		return fmt.Errorf("site não encontrado")
	}
	log.Printf("Site obtido com sucesso: %s", site.Name)

	// Se o site não possui domínio principal, definir o domínio recebido como principal
	if site.CustomDomain == "" {
		log.Printf("Site %s não possui domínio principal. Configurando %s como domínio principal.", site.Name, domain)
		// O valor "12345abc" deve ser o mesmo configurado no registro TXT no DNS
		if err := c.SetDefaultDomain(ctx, siteID, domain, "12345abc"); err != nil {
			return fmt.Errorf("erro ao definir domínio como principal: %w", err)
		}
		return nil
	}

	// Se já existe um domínio personalizado, verificar se o domínio recebido já está cadastrado
	if site.CustomDomain == domain {
		return fmt.Errorf("domínio %s já existe como domínio principal", domain)
	}
	for _, d := range site.DomainAliases {
		if d == domain {
			return fmt.Errorf("domínio %s já existe como alias", domain)
		}
	}

	// Adicionar o domínio como alias
	site.DomainAliases = append(site.DomainAliases, domain)
	log.Printf("Adicionando %s como alias de domínio para o site %s", domain, site.Name)

	siteSetup := &models.SiteSetup{
		Site: *site,
	}
	log.Printf("Atualizando site com novo alias...")
	updatedSite, err := c.netlify.UpdateSite(authCtx, siteSetup)
	if err != nil {
		log.Printf("Erro ao atualizar site com novo domínio: %v", err)
		return fmt.Errorf("erro ao adicionar domínio personalizado: %w", err)
	}

	log.Printf("Domínio %s adicionado com sucesso como alias para o site %s", domain, updatedSite.Name)
	log.Printf("Para configurar o DNS para %s, adicione um registro CNAME apontando para o domínio Netlify", domain)
	return nil
}

// RemoveCustomDomain remove um domínio personalizado de um site
func (c *Client) RemoveCustomDomain(ctx context.Context, siteID, domain string) error {
	log.Printf("Iniciando remoção de domínio %s do site %s", domain, siteID)

	// Verificar se o siteID é válido
	if siteID == "" {
		return fmt.Errorf("ID do site não pode ser vazio")
	}

	// Verificar se o domínio é válido
	if domain == "" {
		return fmt.Errorf("domínio não pode ser vazio")
	}

	// Criar um contexto com autenticação
	authCtx := c.createAuthContext(ctx)

	// Obter o site atual
	site, err := c.netlify.GetSite(authCtx, siteID)
	if err != nil {
		return fmt.Errorf("erro ao obter site: %w", err)
	}

	if site == nil {
		return fmt.Errorf("site não encontrado")
	}

	log.Printf("Site obtido com sucesso: %s", site.Name)

	// Verificar se o domínio está nos aliases
	var domainIndex = -1
	for i, d := range site.DomainAliases {
		if d == domain {
			domainIndex = i
			break
		}
	}

	if domainIndex >= 0 {
		// Remover o domínio dos aliases
		site.DomainAliases = append(site.DomainAliases[:domainIndex], site.DomainAliases[domainIndex+1:]...)
		log.Printf("Removendo %s dos aliases de domínio para o site %s", domain, site.Name)
	} else {
		log.Printf("Domínio %s não encontrado nos aliases do site %s", domain, site.Name)
		return nil
	}

	// Atualizar o site
	siteSetup := &models.SiteSetup{
		Site: *site,
	}

	updatedSite, err := c.netlify.UpdateSite(authCtx, siteSetup)
	if err != nil {
		return fmt.Errorf("erro ao remover domínio personalizado: %w", err)
	}

	log.Printf("Domínio %s removido com sucesso do site %s", domain, updatedSite.Name)

	return nil
}

// SetDefaultDomain define um domínio como o domínio principal de um site.
// Remove o domínio dos aliases (caso exista) para evitar duplicação e envia o valor do TXT para verificação.
func (c *Client) SetDefaultDomain(ctx context.Context, siteID, domain, recordTxtValue string) error {
	log.Printf("Iniciando configuração de domínio %s como padrão para o site %s", domain, siteID)

	if siteID == "" {
		return fmt.Errorf("ID do site não pode ser vazio")
	}
	if domain == "" {
		return fmt.Errorf("domínio não pode ser vazio")
	}

	authCtx := c.createAuthContext(ctx)
	site, err := c.netlify.GetSite(authCtx, siteID)
	if err != nil {
		return fmt.Errorf("erro ao obter site: %w", err)
	}
	if site == nil {
		return fmt.Errorf("site não encontrado")
	}

	// Remover o domínio dos aliases, se existir, para evitar duplicação
	newAliases := []string{}
	for _, alias := range site.DomainAliases {
		if alias != domain {
			newAliases = append(newAliases, alias)
		}
	}
	site.DomainAliases = newAliases

	// Definir o domínio principal
	site.CustomDomain = domain
	log.Printf("Definindo %s como domínio principal para o site %s", domain, site.Name)

	// Monta o payload com o campo "record_txt_value" para verificação
	payload := map[string]interface{}{
		"name":             site.Name,
		"custom_domain":    site.CustomDomain,
		"domain_aliases":   site.DomainAliases,
		"record_txt_value": recordTxtValue,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("erro ao serializar payload: %w", err)
	}

	url := fmt.Sprintf("https://api.netlify.com/api/v1/sites/%s", siteID)
	req, err := http.NewRequest("PATCH", url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("erro ao criar request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.NetlifyToken)

	clientHTTP := &http.Client{Timeout: 10 * time.Second}
	resp, err := clientHTTP.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro ao definir domínio principal: status %d, resposta: %s", resp.StatusCode, string(body))
	}

	log.Printf("Domínio %s definido como principal para o site %s com verificação TXT", domain, site.Name)
	return nil
}

// SwitchDefaultDomain troca o domínio principal de um site.
// Converte o domínio atualmente principal em alias e promove um domínio existente entre os aliases a principal.
func (c *Client) SwitchDefaultDomain(ctx context.Context, siteID, newPrincipalDomain, recordTxtValue string) error {
	log.Printf("Iniciando a troca do domínio principal para %s no site %s", newPrincipalDomain, siteID)

	if siteID == "" {
		return fmt.Errorf("ID do site não pode ser vazio")
	}
	if newPrincipalDomain == "" {
		return fmt.Errorf("domínio não pode ser vazio")
	}

	authCtx := c.createAuthContext(ctx)
	site, err := c.netlify.GetSite(authCtx, siteID)
	if err != nil {
		return fmt.Errorf("erro ao obter site: %w", err)
	}
	if site == nil {
		return fmt.Errorf("site não encontrado")
	}

	// Se o novo domínio já for o principal, retorna erro
	if site.CustomDomain == newPrincipalDomain {
		return fmt.Errorf("o domínio %s já é o principal", newPrincipalDomain)
	}

	// Verificar se o novo domínio existe nos aliases e removê-lo
	found := false
	newAliases := []string{}
	for _, alias := range site.DomainAliases {
		if alias == newPrincipalDomain {
			found = true
		} else {
			newAliases = append(newAliases, alias)
		}
	}
	if !found {
		return fmt.Errorf("o domínio %s não foi encontrado entre os aliases", newPrincipalDomain)
	}

	// Adicionar o domínio atualmente principal aos aliases (se existir e não duplicar)
	oldPrincipal := site.CustomDomain
	if oldPrincipal != "" {
		exists := false
		for _, alias := range newAliases {
			if alias == oldPrincipal {
				exists = true
				break
			}
		}
		if !exists {
			newAliases = append(newAliases, oldPrincipal)
		}
	}

	// Definir o novo domínio como principal e atualizar os aliases
	site.CustomDomain = newPrincipalDomain
	site.DomainAliases = newAliases

	payload := map[string]interface{}{
		"name":             site.Name,
		"custom_domain":    site.CustomDomain,
		"domain_aliases":   site.DomainAliases,
		"record_txt_value": recordTxtValue,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("erro ao serializar payload: %w", err)
	}

	url := fmt.Sprintf("https://api.netlify.com/api/v1/sites/%s", siteID)
	req, err := http.NewRequest("PATCH", url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("erro ao criar request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.NetlifyToken)

	clientHTTP := &http.Client{Timeout: 10 * time.Second}
	resp, err := clientHTTP.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro ao definir domínio principal: status %d, resposta: %s", resp.StatusCode, string(body))
	}

	log.Printf("Domínio %s agora é o principal para o site %s", newPrincipalDomain, site.Name)
	return nil
}

// RemovePrimaryDomain remove o domínio principal de um site
func (c *Client) RemovePrimaryDomain(ctx context.Context, siteID string) error {
	log.Printf("Iniciando remoção do domínio principal do site %s", siteID)

	// Verificar se o siteID é válido
	if siteID == "" {
		return fmt.Errorf("ID do site não pode ser vazio")
	}

	// Criar um contexto com autenticação
	authCtx := c.createAuthContext(ctx)

	// Obter o site atual
	site, err := c.netlify.GetSite(authCtx, siteID)
	if err != nil {
		return fmt.Errorf("erro ao obter site: %w", err)
	}

	if site == nil {
		return fmt.Errorf("site não encontrado")
	}

	log.Printf("Site obtido com sucesso: %s", site.Name)

	// Verificar se existe um domínio principal
	if site.CustomDomain == "" {
		log.Printf("Site %s não tem domínio principal configurado", site.Name)
		return nil
	}

	// Armazenar o valor atual para o log
	oldDomain := site.CustomDomain

	// Limpar o domínio principal
	site.CustomDomain = ""
	log.Printf("Removendo %s como domínio principal do site %s", oldDomain, site.Name)

	// Atualizar o site
	siteSetup := &models.SiteSetup{
		Site: *site,
	}

	updatedSite, err := c.netlify.UpdateSite(authCtx, siteSetup)
	if err != nil {
		return fmt.Errorf("erro ao remover domínio principal: %w", err)
	}

	log.Printf("Domínio principal removido com sucesso do site %s", updatedSite.Name)

	return nil
}

// GetAuth retorna a autenticação do cliente Netlify
func (c *Client) GetAuth() runtime.ClientAuthInfoWriter {
	return c.auth
}

// createAuthContext cria um contexto com informações de autenticação
func (c *Client) createAuthContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "netlify.auth_info", c.auth)
}
