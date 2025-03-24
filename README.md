# Netlify Integration API

API em Golang para integração com a plataforma Netlify, permitindo o deploy de arquivos estáticos armazenados em um bucket S3 e o gerenciamento de domínios personalizados.

## Funcionalidades

1. **Gerenciamento de Domínios**
   - Adicionar domínios personalizados a sites Netlify
   - Remover domínios personalizados
   - Definir domínio principal
   - Remover domínio principal

2. **Deploy de Sites**
   - Deploy de arquivos para a Netlify via upload direto 
   - Suporte a carregamento de arquivos HTML, CSS e JS
   - Suporte para especificar pastas locais

3. **Interface de Administração**
   - Interface web para gerenciamento
   - Listar sites existentes
   - Ver logs de operações

4. **Documentação API**
   - Documentação via Swagger UI
   - Testes de conexão com a Netlify

## Requisitos

- Go 1.21 ou superior
- Credenciais da Netlify (token de acesso)
- Credenciais da AWS (para acesso ao S3, quando usar deploy via S3)
- Domínio configurado na Netlify (opcional)

## Variáveis de Ambiente

```
# Credenciais da Netlify
NETLIFY_TOKEN=seu_token_de_acesso_netlify

# Credenciais da AWS
AWS_ACCESS_KEY_ID=sua_chave_de_acesso_aws
AWS_SECRET_ACCESS_KEY=sua_chave_secreta_aws
AWS_REGION=regiao_aws
S3_BUCKET_NAME=nome_do_bucket_s3

# Configurações da API
API_PORT=8080
GIN_MODE=debug  # Use 'release' em produção
BASE_DOMAIN=sites.seudominio.com.br
```

## Como Usar

### Executando Localmente

```bash
# Clonar o repositório
git clone https://github.com/renatoroquejani/poc-netlify.git
cd poc-netlify

# Configurar variáveis de ambiente
cp .env.example .env
# Edite o arquivo .env com suas credenciais

# Executar a aplicação
go run main.go
```

Acesse a interface web em `http://localhost:8080` e a documentação do Swagger em `http://localhost:8080/docs/swagger/index.html`

### API REST

#### Verificar Status

```
GET /api/status
```

Resposta:
```json
{
  "status": "online",
  "version": "1.0.0",
  "time": "2025-03-24T15:45:09-03:00"
}
```

#### Deploy de Site

```
POST /api/deploy/site
Content-Type: multipart/form-data

site_name: meu-site-teste
description: Site para testes
custom_domain: meu-site.exemplo.com
file: [arquivo HTML]
```

Resposta:
```json
{
  "success": true,
  "message": "Site criado com sucesso",
  "site_id": "12345abcde",
  "site_url": "https://meu-site-teste.netlify.app",
  "created_at": "2025-03-24T16:45:09-03:00",
  "test_success": true
}
```

#### Adicionar Domínio Personalizado

```
POST /api/domains/add
Content-Type: application/json

{
  "site_id": "12345abcde",
  "domain": "meu-site.exemplo.com"
}
```

Resposta:
```json
{
  "success": true,
  "message": "Domínio adicionado com sucesso",
  "site_id": "12345abcde",
  "domain": "meu-site.exemplo.com"
}
```

## Estrutura do Projeto

```
/docs           # Documentação Swagger
/internal       # Código interno da aplicação
  /api          # API web
  /netlify      # Integração com a Netlify
  /aws          # Integração com AWS S3
  /config       # Configurações da aplicação
/web            # Interface web
  /static       # Arquivos estáticos (HTML, CSS, JS)
main.go         # Ponto de entrada principal
```

## Documentação

A documentação completa da API está disponível através do Swagger UI em `/docs/swagger/index.html`.