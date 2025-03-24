# Netlify Integration PoC

Prova de Conceito (PoC) em Golang para integração com a plataforma Netlify, permitindo o deploy de arquivos estáticos armazenados em um bucket S3.

## Funcionalidades

1. **Verificação e Criação de Sites na Netlify**
   - Verifica se um site já existe para um subdomínio específico
   - Cria um novo site se necessário
   - Formato do subdomínio: `<usuario>.sites.kodestech.com.br`

2. **Configuração de DNS**
   - Configura DNS para o subdomínio padrão
   - Permite apontar domínios personalizados para o site

3. **Deploy de Arquivos Estáticos**
   - Conecta-se ao bucket S3 para obter os arquivos
   - Realiza o deploy dos arquivos na Netlify
   - Verifica o sucesso do deploy

4. **Interface Web e API**
   - Interface web amigável para iniciar deploys
   - API REST para integração com outros sistemas

5. **Ambiente de Teste da Netlify**
   - Rota dedicada para testes rápidos da API Netlify
   - Interface gráfica para criação de sites de teste
   - Documentação via Swagger

## Requisitos

- Go 1.21 ou superior
- Credenciais da Netlify (token de acesso)
- Credenciais da AWS (para acesso ao S3)
- Domínio principal configurado na Netlify: `sites.kodestech.com.br`

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
```

## Como Usar

### Executando Localmente

```bash
# Clonar o repositório
git clone https://github.com/kodestech/poc-netlify.git
cd poc-netlify

# Configurar variáveis de ambiente
cp .env.example .env
# Edite o arquivo .env com suas credenciais

# Executar a aplicação
go run main.go
```

Acesse a interface web em `http://localhost:8080`

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
  "time": "2025-03-18T15:45:09-03:00"
}
```

#### Iniciar Deploy

```
POST /api/deploy
Content-Type: application/json

{
  "username": "renato",
  "custom_domain": "renato.com.br",
  "s3_path": "sites/renato"
}
```

Resposta:
```json
{
  "success": true,
  "message": "Deploy iniciado com sucesso",
  "subdomain": "renato.sites.kodestech.com.br"
}
```

#### Testar Criação de Sites

```
POST /api/test/netlify
Content-Type: application/json

{
  "site_name": "teste-site",
  "description": "Site para fins de teste",
  "test_content": "<h1>Hello World</h1>",
  "cleanup_after": true,
  "use_custom_domain": false
}
```

Resposta:
```json
{
  "success": true,
  "message": "Site de teste criado com sucesso",
  "site_id": "12345abcde",
  "site_url": "https://teste-site.netlify.app",
  "created_at": "2025-03-19T16:45:09-03:00",
  "test_success": true
}
```

## Estrutura do Projeto

```
/cmd            # Ponto de entrada da aplicação (legado)
/internal       # Código interno da aplicação
  /api          # API web
  /netlify      # Integração com a Netlify
  /aws          # Integração com AWS S3
  /config       # Configurações da aplicação
/web            # Interface web
  /static       # Arquivos estáticos (HTML, CSS, JS)
/pkg            # Pacotes reutilizáveis
main.go         # Ponto de entrada principal
