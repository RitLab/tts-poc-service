package command

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"mime/multipart"
	"strings"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/pkg/common/decorator"
	pkgError "tts-poc-service/pkg/common/error"
	"tts-poc-service/pkg/pdf/domain"
)

type VerifyPdfFileQuery struct {
	File *multipart.FileHeader
	Key  string `form:"key" validate:"required"`
}

type VerifyPdfFileHandler decorator.CommandHandler[VerifyPdfFileQuery]

type verifyPdfFileRepository struct {
	repo   domain.QueryRepository
	logger *baselogger.Logger
}

func NewVerifyPdfFilesRepository(repo domain.QueryRepository, log *baselogger.Logger) decorator.CommandHandler[VerifyPdfFileQuery] {
	return decorator.ApplyCommandDecorators[VerifyPdfFileQuery](
		verifyPdfFileRepository{repo: repo, logger: log},
		log)
}

func (g verifyPdfFileRepository) Handle(ctx context.Context, in VerifyPdfFileQuery) (err error) {
	if err = domain.ValidatePdfFile(in.File); err != nil {
		return err
	}

	var blockKey *domain.Ed25519BlockKey
	block, err := g.repo.GetBlockById(ctx, in.Key)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error getting block by key: %w", err))
		return err
	}
	blockKey = domain.LoadKey(block)

	src, err := in.File.Open()
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error open file: %w", err))
		return err
	}
	defer src.Close()

	pdfReader, err := domain.NewPdfReader(src)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error read file: %w", err))
		return err
	}
	cleanedText, err := pdfReader.CleanText()
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error clean text: %w", err))
		return err
	}

	pdfCtx, err := api.ReadContext(src, model.NewDefaultConfiguration())
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error read context: %w", err))
		return err
	}
	hashContent := sha256.Sum256([]byte(cleanedText))

	if signature, found := pdfCtx.RootDict.Find("signature"); !found {
		return fmt.Errorf(pkgError.SIGNATURE_NOT_FOUND)
	} else {
		sig := strings.Replace(strings.Replace(signature.String(), "(", "", 1), ")", "", len(signature.String())-1)
		sigHex, _ := hex.DecodeString(sig)

		if isValid := blockKey.VerifySignature(sigHex, hashContent[:]); !isValid {
			return fmt.Errorf(pkgError.SIGNATURE_NOT_VALID)
		}
	}

	return nil
}
