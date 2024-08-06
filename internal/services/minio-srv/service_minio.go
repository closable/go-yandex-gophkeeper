package miniosrv

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioService struct {
	Host       string
	ID         string
	Secret     string
	BucketName string
	SSL        bool
}

func NewMinioService(host, s3Bucket, s3AccessKey, s3SecretKey string) *MinioService {
	return &MinioService{
		Host:       host,
		ID:         s3AccessKey,
		Secret:     s3SecretKey,
		SSL:        false,
		BucketName: s3Bucket,
	}
}

func (m *MinioService) Upload(fileName, mark string) (int64, error) {
	s3Client, err := minio.New(m.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(m.ID, m.Secret, ""),
		Secure: m.SSL,
	})
	if err != nil {
		return 0, err
	}

	object, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer object.Close()
	objectStat, err := object.Stat()
	if err != nil {
		return 0, err
	}

	info, err := s3Client.PutObject(context.Background(), m.BucketName, mark, object, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return 0, err
	}
	fmt.Printf("Uploaded %s of size: %d", mark, info.Size)
	//fmt.Println(info.ETag, info.Key)

	return info.Size, nil
}

func (m *MinioService) Download(mark, filePath string) error {
	s3Client, err := minio.New(m.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(m.ID, m.Secret, ""),
		Secure: m.SSL,
	})
	if err != nil {
		return err
	}

	reader, err := s3Client.GetObject(context.Background(), m.BucketName, mark, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()

	localFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	stat, err := reader.Stat()
	if err != nil {
		return err
	}

	if _, err := io.CopyN(localFile, reader, stat.Size); err != nil {
		return err
	}
	fmt.Println("File downloaded successfully path ", filePath)
	return nil
}

func (m *MinioService) Delete(mark string) error {
	s3Client, err := minio.New(m.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(m.ID, m.Secret, ""),
		Secure: m.SSL,
	})
	if err != nil {
		return err
	}

	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}

	err = s3Client.RemoveObject(context.Background(), m.BucketName, mark, opts)
	if err != nil {
		return err
	}

	fmt.Printf("Delete success name %s", mark)
	return nil
}
