package query

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"io"
	"mime/multipart"
	"os"
	"tts-poc-service/config"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/lib/storage"
	"tts-poc-service/pkg/common/constant"
	"tts-poc-service/pkg/common/decorator"
	pkgError "tts-poc-service/pkg/common/error"
	"tts-poc-service/pkg/pdf/domain"
)

type JoinPdfFilesQuery struct {
	Files []*multipart.FileHeader
}

type JoinPdfFilesResponse struct {
	Url string `json:"url"`
}

type JoinPdfFilesHandler decorator.QueryHandler[JoinPdfFilesQuery, JoinPdfFilesResponse]

type joinPdfFilesRepository struct {
	s3     storage.Storage
	logger *baselogger.Logger
}

func NewJoinPdfFilesRepository(s3 storage.Storage, log *baselogger.Logger) decorator.QueryHandler[JoinPdfFilesQuery, JoinPdfFilesResponse] {
	return decorator.ApplyQueryDecorators[JoinPdfFilesQuery, JoinPdfFilesResponse](
		joinPdfFilesRepository{s3: s3, logger: log},
		log)
}

func (g joinPdfFilesRepository) Handle(ctx context.Context, in JoinPdfFilesQuery) (JoinPdfFilesResponse, error) {
	if len(in.Files) < 2 {
		return JoinPdfFilesResponse{}, fmt.Errorf(pkgError.FILE_SHOULD_MORE_THAN_TWO)
	}
	filePaths := make([]string, len(in.Files))
	for idx, file := range in.Files {
		err := func() error {
			if err := domain.ValidatePdfFile(file); err != nil {
				return err
			}
			outputFile := fmt.Sprintf("%s/%s-%s", constant.PDF_FOLDER, uuid.NewString(), file.Filename)

			src, err := file.Open()
			if err != nil {
				g.logger.Hashcode(ctx).Error(fmt.Errorf("error open file: %w", err))
				return err
			}
			defer src.Close()

			dst, err := os.Create(outputFile)
			if err != nil {
				g.logger.Hashcode(ctx).Error(fmt.Errorf("error create file: %w", err))
				return err
			}
			defer dst.Close()

			if _, err = io.Copy(dst, src); err != nil {
				g.logger.Hashcode(ctx).Error(fmt.Errorf("error copy to a new file: %w", err))
				return err
			}
			filePaths[idx] = outputFile
			return nil
		}()
		if err != nil {
			return JoinPdfFilesResponse{}, err
		}
	}
	outputFile := fmt.Sprintf("%s/output-%s.pdf", constant.PDF_FOLDER, uuid.NewString())

	err := api.MergeCreateFile(filePaths, outputFile, false, nil)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error merge file: %w", err))
		return JoinPdfFilesResponse{}, err
	}

	// put file to storage
	err = g.s3.PutObject(ctx, &storage.PutFileRequest{Path: outputFile, ContentType: "application/pdf"})
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error put file: %w", err))
		return JoinPdfFilesResponse{}, err
	}

	filePaths = append(filePaths, outputFile)
	for _, path := range filePaths {
		err = os.Remove(path)
		if err != nil {
			g.logger.Hashcode(ctx).Error(fmt.Errorf("error removing file: %w", err))
		}
	}

	return JoinPdfFilesResponse{
		Url: fmt.Sprintf("%s://%s/%s/%s", config.Config.Storage.Method, config.Config.Storage.Endpoint, config.Config.Storage.BucketName, outputFile)}, nil
}
