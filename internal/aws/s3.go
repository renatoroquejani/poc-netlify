package aws

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/kodestech/poc-netlify/internal/config"
)

// S3Client encapsula a integração com o AWS S3
type S3Client struct {
	client *s3.Client
	config *config.Config
}

// NewS3Client cria um novo cliente S3
func NewS3Client(cfg *config.Config) (*S3Client, error) {
	var s3Client *s3.Client

	// Verificar se estamos usando um endpoint personalizado (MinIO)
	if cfg.S3Endpoint != "" {
		log.Printf("Usando endpoint S3 personalizado: %s", cfg.S3Endpoint)
		
		// Configurar cliente para MinIO
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:               cfg.S3Endpoint,
				SigningRegion:     "us-east-1", // Região padrão para MinIO
				HostnameImmutable: true,
			}, nil
		})

		// Carregar configurações com endpoint personalizado
		awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
			awsconfig.WithEndpointResolverWithOptions(customResolver),
			awsconfig.WithCredentialsProvider(aws.CredentialsProviderFunc(
				func(ctx context.Context) (aws.Credentials, error) {
					return aws.Credentials{
						AccessKeyID:     cfg.AWSAccessKeyID,
						SecretAccessKey: cfg.AWSSecretAccessKey,
					}, nil
				},
			)),
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao carregar configurações para MinIO: %w", err)
		}

		// Criar cliente S3 com configurações personalizadas
		s3Client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
			o.UsePathStyle = true // Importante para MinIO
		})
	} else {
		// Configurar cliente para AWS S3 padrão
		awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
			awsconfig.WithRegion(cfg.AWSRegion),
			awsconfig.WithCredentialsProvider(aws.CredentialsProviderFunc(
				func(ctx context.Context) (aws.Credentials, error) {
					return aws.Credentials{
						AccessKeyID:     cfg.AWSAccessKeyID,
						SecretAccessKey: cfg.AWSSecretAccessKey,
					}, nil
				},
			)),
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao carregar configurações da AWS: %w", err)
		}

		// Criar cliente S3 padrão
		s3Client = s3.NewFromConfig(awsCfg)
	}

	return &S3Client{
		client: s3Client,
		config: cfg,
	}, nil
}

// DownloadFiles baixa os arquivos do bucket S3 para um diretório local
func (c *S3Client) DownloadFiles(ctx context.Context, localDir string) error {
	log.Printf("Baixando arquivos do bucket %s, caminho %s para %s", c.config.S3BucketName, c.config.S3Path, localDir)

	// Criar diretório local se não existir
	if err := os.MkdirAll(localDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório local: %w", err)
	}

	// Listar objetos no bucket
	prefix := c.config.S3Path
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	paginator := s3.NewListObjectsV2Paginator(c.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(c.config.S3BucketName),
		Prefix: aws.String(prefix),
	})

	fileCount := 0
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("erro ao listar objetos do S3: %w", err)
		}

		for _, obj := range page.Contents {
			// Pular diretórios (objetos que terminam com "/")
			key := *obj.Key
			if strings.HasSuffix(key, "/") {
				continue
			}

			// Determinar o caminho relativo do arquivo
			relPath := key
			if prefix != "" {
				relPath = strings.TrimPrefix(key, prefix)
			}

			// Caminho completo do arquivo local
			localPath := filepath.Join(localDir, relPath)

			// Criar diretórios pai se necessário
			if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
				return fmt.Errorf("erro ao criar diretório para %s: %w", localPath, err)
			}

			// Baixar o arquivo
			if err := c.downloadFile(ctx, key, localPath); err != nil {
				return fmt.Errorf("erro ao baixar arquivo %s: %w", key, err)
			}

			fileCount++
			log.Printf("Arquivo baixado: %s", key)
		}
	}

	log.Printf("Download concluído: %d arquivos baixados para %s", fileCount, localDir)
	return nil
}

// downloadFile baixa um único arquivo do S3
func (c *S3Client) downloadFile(ctx context.Context, s3Key, localPath string) error {
	// Obter o objeto do S3
	resp, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.config.S3BucketName),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		return fmt.Errorf("erro ao obter objeto do S3: %w", err)
	}
	defer resp.Body.Close()

	// Criar arquivo local
	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo local: %w", err)
	}
	defer file.Close()

	// Copiar conteúdo do S3 para o arquivo local
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao copiar conteúdo do arquivo: %w", err)
	}

	return nil
}
