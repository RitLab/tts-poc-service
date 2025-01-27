package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"
	"tts-poc-service/config"
	"tts-poc-service/lib/baselogger"
	pkgError "tts-poc-service/pkg/common/error"
)

type Storage interface {
	PutMultipartObject(ctx context.Context, input *PutObjectRequest) error
	PutObject(ctx context.Context, input *PutFileRequest) error
	GetObject(ctx context.Context, bucketName, objectName string) ([]byte, error)
	PreSignedGetObject(ctx context.Context, input *GetObjectRequest) (*url.URL, error)
	DeleteObject(ctx context.Context, objectName string) error
}

type PutObjectRequest struct {
	Key        string
	File       multipart.File
	FileHeader *multipart.FileHeader
}

type PutFileRequest struct {
	Path string
}

type GetObjectRequest struct {
	Filename  string
	Duration  time.Duration
	ReqParams url.Values
}

type MinioHandler struct {
	log        *baselogger.Logger
	minio      *minio.Client
	bucketName string
}

func NewMinioHandler(log *baselogger.Logger) (mapStorages Storage) {
	minioClient, err := minio.New(config.Config.Storage.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Config.Storage.AccessKey, config.Config.Storage.SecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	return &MinioHandler{
		log:        log,
		minio:      minioClient,
		bucketName: config.Config.Storage.BucketName,
	}
}

func (m *MinioHandler) PutMultipartObject(ctx context.Context, input *PutObjectRequest) error {
	// get the file size an read the file content into a buffer
	size := input.FileHeader.Size
	buffer := make([]byte, size)
	input.File.Read(buffer)
	body := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	info, err := m.minio.PutObject(ctx, m.bucketName, input.Key, body, size, minio.PutObjectOptions{ContentType: fileType})
	if err != nil {
		m.log.Error(err)
		return err
	}
	m.log.Infof("Successfully uploaded %s of size %d\n", input.FileHeader.Filename, info.Size)
	return nil
}

func (m *MinioHandler) PutObject(ctx context.Context, input *PutFileRequest) error {
	file, err := os.Open(input.Path)
	if err != nil {
		m.log.Hashcode(ctx).Error(fmt.Errorf("error open file: %w", err))
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		m.log.Hashcode(ctx).Error(fmt.Errorf("failed to stat file: %w", err))
		return err
	}
	fileSize := fileInfo.Size()

	//buffer := make([]byte, fileSize)
	//_, err = file.Read(buffer)
	//if err != nil {
	//	m.log.Hashcode(ctx).Error(fmt.Errorf("failed to read file: %w", err))
	//	return err
	//}
	//
	//contentType := http.DetectContentType(buffer)

	info, err := m.minio.PutObject(ctx, m.bucketName, input.Path, file, fileSize, minio.PutObjectOptions{ContentType: "audio/mpeg"})
	if err != nil {
		m.log.Error(err)
		return err
	}
	m.log.Infof("Successfully uploaded %s of size %d\n", fileInfo.Name(), info.Size)
	return nil
}

func (m *MinioHandler) GetObject(ctx context.Context, bucketName, objectName string) ([]byte, error) {
	obj, err := m.minio.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		m.log.Error(err)
		return []byte{}, fmt.Errorf(pkgError.GENERAL_ERROR)
	}
	defer obj.Close()

	content := new(bytes.Buffer)
	if _, err = io.Copy(content, obj); err != nil {
		m.log.Error(err)
		return nil, fmt.Errorf(pkgError.GENERAL_ERROR)
	}
	return content.Bytes(), nil
}

func (m *MinioHandler) PreSignedGetObject(ctx context.Context, input *GetObjectRequest) (*url.URL, error) {
	url, err := m.minio.PresignedGetObject(ctx, m.bucketName, input.Filename, input.Duration, input.ReqParams)
	if err != nil {
		m.log.Error(err)
		return nil, err
	}
	m.log.Infof("Get object url: %s", url.String())
	return url, err
}

func (m *MinioHandler) DeleteObject(ctx context.Context, objectName string) error {
	err := m.minio.RemoveObject(ctx, m.bucketName, objectName, minio.RemoveObjectOptions{
		ForceDelete:      true,
		GovernanceBypass: true,
	})
	if err != nil {
		m.log.Error(err)
		return err
	}
	return nil
}
