package query

import (
	"context"
	"fmt"
	"github.com/aarondl/null/v8"
	"os"
	"sync"
	"time"
	"tts-poc-service/database/models"
	"tts-poc-service/lib/htgo"
	"tts-poc-service/pkg/common/constant"

	"tts-poc-service/lib/baselogger"
	"tts-poc-service/lib/gemini_ai"
	"tts-poc-service/lib/storage"
	"tts-poc-service/pkg/common/decorator"
	"tts-poc-service/pkg/common/utils"
	"tts-poc-service/pkg/story/domain"
)

type GenerateStoryBookQuery struct {
	TotalPages    int    `json:"total_pages" validate:"required"`
	Language      string `json:"language" validate:"required"`
	CharacterName string `json:"character_name" validate:"required"`
	Theme         string `json:"theme" validate:"required"`
	Context       string `json:"context"`
}

type GenerateStoryBookResponse struct {
	domain.BookResponse
}

type GenerateStoryBookHandler decorator.QueryHandler[GenerateStoryBookQuery, GenerateStoryBookResponse]

type generateStoryBookRepository struct {
	ai     gemini_ai.GenAIMethod
	player htgo.Player
	repo   domain.CommandRepository
	s3     storage.Storage
	logger *baselogger.Logger
}

func NewGenerateStoryBookRepository(ai gemini_ai.GenAIMethod, player htgo.Player, repo domain.CommandRepository, s3 storage.Storage,
	log *baselogger.Logger) decorator.QueryHandler[GenerateStoryBookQuery, GenerateStoryBookResponse] {
	return decorator.ApplyQueryDecorators[GenerateStoryBookQuery, GenerateStoryBookResponse](
		generateStoryBookRepository{ai: ai, player: player, repo: repo, s3: s3, logger: log},
		log)
}

func (g generateStoryBookRepository) Handle(ctx context.Context, in GenerateStoryBookQuery) (GenerateStoryBookResponse, error) {
	language := domain.LanguageMapping(in.Language)
	theme := domain.ThemeMapping(in.Theme)
	generatedStoryBook, err := g.ai.GenerateStory(ctx, &domain.GenerateBook{
		TotalPages:    in.TotalPages,
		CharacterName: in.CharacterName,
		Context:       in.Context,
	}, in.Language, theme)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("generate story book failed: %w", err))
		return GenerateStoryBookResponse{}, err
	}

	bookResponse := make([]domain.PageResponse, len(generatedStoryBook.Pages))
	bookDatabase := make([]domain.PageDatabase, len(generatedStoryBook.Pages))
	var wg sync.WaitGroup
	id := utils.NewToken(10)
	for idx, book := range generatedStoryBook.Pages {
		wg.Add(1)
		go func(input domain.GeneratedAIPage) {
			defer wg.Done()
			eachImageFile, err := g.ai.GenerateImageStory(ctx, book, id)
			if err != nil {
				g.logger.Hashcode(ctx).Error(fmt.Errorf("generate story book id %s page %d failed: %w", id, input.PageNumber, err))
			}
			err = g.s3.PutObject(ctx, &storage.PutFileRequest{Path: eachImageFile, ContentType: "image/png"})
			if err != nil {
				g.logger.Hashcode(ctx).Error(fmt.Errorf("error put image file id %s page %d: %w", id, input.PageNumber, err))
			}
			imageUrl, err := g.s3.PreSignedGetObject(ctx, &storage.GetObjectRequest{
				Filename: eachImageFile,
				Duration: 30 * time.Minute,
			})
			if err != nil {
				g.logger.Hashcode(ctx).Error(fmt.Errorf("error get presigned url id %s page %d: %w", id, input.PageNumber, err))
			}

			eachAudioFile, err := g.player.SaveOneFile(book.Body, language, id, book.PageNumber)
			if err != nil {
				g.logger.Hashcode(ctx).Error(fmt.Errorf("generate audio book id %s page %d failed: %w", id, input.PageNumber, err))
			}
			err = g.s3.PutObject(ctx, &storage.PutFileRequest{Path: eachAudioFile, ContentType: "audio/wav"})
			if err != nil {
				g.logger.Hashcode(ctx).Error(fmt.Errorf("error put audio file id %s page %d: %w", id, input.PageNumber, err))
			}
			audioUrl, err := g.s3.PreSignedGetObject(ctx, &storage.GetObjectRequest{
				Filename: eachAudioFile,
				Duration: 30 * time.Minute,
			})
			if err != nil {
				g.logger.Hashcode(ctx).Error(fmt.Errorf("error get presigned url id %s page %d: %w", id, input.PageNumber, err))
			}
			bookDatabase[idx] = domain.PageDatabase{
				Page: domain.Page{
					PageNumber: book.PageNumber,
					Body:       book.Body,
				},
				ImagePath: eachImageFile,
				AudioPath: eachAudioFile,
			}
			bookResponse[idx] = domain.PageResponse{
				Page: domain.Page{
					PageNumber: book.PageNumber,
					Body:       book.Body,
				},
				ImageUrl: utils.IfThenElse(imageUrl == nil, "", imageUrl.String()),
				AudioUrl: utils.IfThenElse(audioUrl == nil, "", audioUrl.String()),
			}
		}(book)
	}
	wg.Wait()

	storyJson := null.JSON{}
	err = storyJson.Marshal(bookDatabase)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error marshal story database to json: %w", err))
		return GenerateStoryBookResponse{}, err
	}
	err = g.repo.InsertStoryBook(ctx, &models.StoryBook{
		ID:             id,
		Title:          bookDatabase[0].Body,
		CharacterName:  in.CharacterName,
		Language:       in.Language,
		Theme:          in.Theme,
		TitleImagePath: null.StringFrom(bookDatabase[0].ImagePath),
		Story:          storyJson,
	})

	imageFiles, err := os.ReadDir(constant.IMAGE_FOLDER)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error read image folder: %w", err))
	}

	// Loop through files and delete each
	for _, file := range imageFiles {
		filePath := constant.IMAGE_FOLDER + "/" + file.Name()
		err = os.Remove(filePath)
		if err != nil {
			g.logger.Hashcode(ctx).Error(fmt.Errorf("error remove image file: %w", err))
		}
	}

	audioFiles, err := os.ReadDir(constant.AUDIO_FOLDER)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("error read audio folder: %w", err))
	}

	// Loop through files and delete each
	for _, file := range audioFiles {
		filePath := constant.AUDIO_FOLDER + "/" + file.Name()
		err = os.Remove(filePath)
		if err != nil {
			g.logger.Hashcode(ctx).Error(fmt.Errorf("error remove audio file: %w", err))
		}
	}

	return GenerateStoryBookResponse{
		BookResponse: domain.BookResponse{
			Id:    id,
			Pages: bookResponse,
		},
	}, nil
}
