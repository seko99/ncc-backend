package s3storage

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	cfg    *config.Config
	log    logger.Logger
	ctx    context.Context
	client *minio.Client
}

func (ths *Storage) PutFile(localFileName string, remoteFileName string) error {

	uinfo, err := ths.client.FPutObject(ths.ctx, ths.cfg.S3.BucketName, remoteFileName, localFileName, minio.PutObjectOptions{ContentType: "text/plain"})
	if err != nil {
		return fmt.Errorf("can't put file to storage: %w", err)
	}

	if uinfo.Size == 0 {
		return fmt.Errorf("zero len file")
	}

	return nil
}

func (ths *Storage) CreateBucket() (bool, error) {
	err := ths.client.MakeBucket(ths.ctx, ths.cfg.S3.BucketName, minio.MakeBucketOptions{})
	if err != nil {
		if exists, _ := ths.client.BucketExists(ths.ctx, ths.cfg.S3.BucketName); exists {
			return true, nil
		}
		return false, fmt.Errorf("can't make bucket: %w", err)
	}
	return false, nil
}

func (ths *Storage) Connect() error {
	endpoint := fmt.Sprintf("%s:%d", ths.cfg.S3.Host, ths.cfg.S3.Port)
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(ths.cfg.S3.AccessKeyID, ths.cfg.S3.SecretAccessKey, ""),
		Secure: ths.cfg.S3.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("cna't connect to S3 storage: %w", err)
	}

	ths.client = client

	return nil
}

func NewS3(cfg *config.Config, log logger.Logger) *Storage {
	return &Storage{
		cfg: cfg,
		log: log,
		ctx: context.Background(),
	}
}
