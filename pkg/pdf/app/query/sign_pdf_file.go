package query

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"mime/multipart"
	"os"
	"rsc.io/pdf"
	"strings"
	"tts-poc-service/config"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/lib/storage"
	"tts-poc-service/models"
	"tts-poc-service/pkg/common/constant"
	"tts-poc-service/pkg/common/decorator"
	pkgError "tts-poc-service/pkg/common/error"
	"tts-poc-service/pkg/pdf/domain"
)

type SignPdfFileQuery struct {
	File *multipart.FileHeader
	Key  string `form:"key" validate:"required"`
}

type SignPdfFileResponse struct {
	Url string `json:"url"`
}

type SignPdfFileHandler decorator.QueryHandler[SignPdfFileQuery, SignPdfFileResponse]

type signPdfFileRepository struct {
	repo     domain.CommandRepository
	readRepo domain.QueryRepository
	s3       storage.Storage
	logger   *baselogger.Logger
}

func NewSignPdfFilesRepository(repo domain.CommandRepository, readRepo domain.QueryRepository, s3 storage.Storage,
	log *baselogger.Logger) decorator.QueryHandler[SignPdfFileQuery, SignPdfFileResponse] {
	return decorator.ApplyQueryDecorators[SignPdfFileQuery, SignPdfFileResponse](
		signPdfFileRepository{repo: repo, readRepo: readRepo, s3: s3, logger: log},
		log)
}

func (g signPdfFileRepository) Handle(ctx context.Context, in SignPdfFileQuery) (SignPdfFileResponse, error) {
	if err := domain.ValidatePdfFile(in.File); err != nil {
		return SignPdfFileResponse{}, err
	}

	var blockKey *domain.Ed25519BlockKey
	block, err := g.readRepo.GetBlockById(ctx, in.Key)
	if err == nil {
		blockKey = domain.LoadKey(block)
	} else if err.Error() == pkgError.PRIVATE_KEY_NOT_FOUND {
		key := domain.Ed25519Key{Key: in.Key}
		blockKey, err = key.GeneratePublicAndPrivateKey()
		if err != nil {
			return SignPdfFileResponse{}, err
		}
		err = g.repo.InsertBlockKey(ctx, &models.Block{
			ID:         in.Key,
			PrivateKey: hex.EncodeToString(blockKey.PrivateKey),
			PublicKey:  hex.EncodeToString(blockKey.PublicKey),
		})
		if err != nil {
			return SignPdfFileResponse{}, err
		}
	} else {
		return SignPdfFileResponse{}, err
	}

	src, err := in.File.Open()
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error open file: %w", err))
		return SignPdfFileResponse{}, err
	}
	defer src.Close()

	doc, err := pdf.NewReader(src, in.File.Size)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error open file: %w", err))
		return SignPdfFileResponse{}, err
	}

	pdfCtx, err := api.ReadContext(src, model.NewDefaultConfiguration())
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error read context: %w", err))
		return SignPdfFileResponse{}, err
	}

	content := make([]string, 0)
	numPages := doc.NumPage()
	for i := 1; i <= numPages; i++ {
		page := doc.Page(i)
		if page.V.IsNull() {
			continue
		}
		for _, v := range page.Content().Text {
			content = append(content, v.S)
		}
	}
	hashContent := sha256.Sum256([]byte(strings.Join(content, "")))
	signature := blockKey.SignData(hashContent[:])
	sig := hex.EncodeToString(signature)
	pdfCtx.RootDict.InsertString("signature", sig)

	outputFile := fmt.Sprintf("%s/output-%s.pdf", constant.PDF_FOLDER, uuid.NewString())

	err = api.WriteContextFile(pdfCtx, outputFile)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error merge file: %w", err))
		return SignPdfFileResponse{}, err
	}

	// put file to storage
	err = g.s3.PutObject(ctx, &storage.PutFileRequest{Path: outputFile, ContentType: "application/pdf"})
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error put file: %w", err))
		return SignPdfFileResponse{}, err
	}

	err = os.Remove(outputFile)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error removing file: %w", err))
	}

	return SignPdfFileResponse{
		Url: fmt.Sprintf("%s://%s/%s/%s", config.Config.Storage.Method, config.Config.Storage.ExternalEndpoint, config.Config.Storage.BucketName, outputFile)}, nil
}
