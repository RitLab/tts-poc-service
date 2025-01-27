package query

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"os"
	"tts-poc-service/config"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/lib/storage"
	"tts-poc-service/pkg/common/constant"
	"tts-poc-service/pkg/common/decorator"
	"tts-poc-service/pkg/tts/domain"
)

type JoinMp3FilesQuery struct {
	Files []*multipart.FileHeader
}

type JoinMp3FilesResponse struct {
	Url string `json:"url"`
}

type JoinMp3FilesHandler decorator.QueryHandler[JoinMp3FilesQuery, JoinMp3FilesResponse]

type joinMp3FilesRepository struct {
	s3     storage.Storage
	logger *baselogger.Logger
}

func NewJoinMp3FilesRepository(s3 storage.Storage, log *baselogger.Logger) decorator.QueryHandler[JoinMp3FilesQuery, JoinMp3FilesResponse] {
	return decorator.ApplyQueryDecorators[JoinMp3FilesQuery, JoinMp3FilesResponse](
		joinMp3FilesRepository{s3: s3, logger: log},
		log)
}

func (g joinMp3FilesRepository) Handle(ctx context.Context, in JoinMp3FilesQuery) (JoinMp3FilesResponse, error) {
	filePaths := make([]string, len(in.Files))
	for idx, file := range in.Files {
		err := func() error {
			if err := domain.ValidateAudioFile(file); err != nil {
				return err
			}
			outputFile := fmt.Sprintf("%s/%s-%s", constant.AUDIO_FOLDER, uuid.NewString(), file.Filename)

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
			return JoinMp3FilesResponse{}, err
		}
	}
	outputFile := fmt.Sprintf("%s/output-%s.mp3", constant.AUDIO_FOLDER, uuid.NewString())

	// Open the output file
	out, err := os.Create(outputFile)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error creating file: %w", err))
		return JoinMp3FilesResponse{}, nil
	}
	defer out.Close()

	// Loop through the input files and concatenate their data
	for _, path := range filePaths {
		err = domain.AppendFile(out, path)
		if err != nil {
			g.logger.Hashcode(ctx).Error(fmt.Errorf("error appending file: %w", err))
			return JoinMp3FilesResponse{}, nil
		}
	}

	// put file to storage
	err = g.s3.PutObject(ctx, &storage.PutFileRequest{Path: outputFile})
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error put file: %w", err))
		return JoinMp3FilesResponse{}, err
	}

	filePaths = append(filePaths, outputFile)
	for _, path := range filePaths {
		err = os.Remove(path)
		if err != nil {
			g.logger.Hashcode(ctx).Error(fmt.Errorf("error removing file: %w", err))
		}
	}

	return JoinMp3FilesResponse{
		Url: fmt.Sprintf("%s://%s/%s/%s", config.Config.Storage.Method, config.Config.Storage.Endpoint, config.Config.Storage.BucketName, outputFile)}, nil
}
