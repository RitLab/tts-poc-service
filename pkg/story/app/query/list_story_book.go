package query

import (
	"context"
	"fmt"
	"sync"
	"tts-poc-service/database/models"

	"tts-poc-service/lib/baselogger"
	"tts-poc-service/lib/storage"
	"tts-poc-service/pkg/common/decorator"
	"tts-poc-service/pkg/story/domain"
)

type ListStoryBookQuery struct {
}

type ListStoryBookResponse struct {
	domain.ListStoryBookResponse
}

type ListStoryBookHandler decorator.QueryHandler[ListStoryBookQuery, ListStoryBookResponse]

type listStoryBookRepository struct {
	repo   domain.QueryRepository
	s3     storage.Storage
	logger *baselogger.Logger
}

func NewListStoryBookRepository(repo domain.QueryRepository, s3 storage.Storage,
	log *baselogger.Logger) decorator.QueryHandler[ListStoryBookQuery, ListStoryBookResponse] {
	return decorator.ApplyQueryDecorators[ListStoryBookQuery, ListStoryBookResponse](
		listStoryBookRepository{repo: repo, s3: s3, logger: log},
		log)
}

func (g listStoryBookRepository) Handle(ctx context.Context, _ ListStoryBookQuery) (ListStoryBookResponse, error) {
	listStoryBook, err := g.repo.ListStoryBooks(ctx)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("list story book failed: %w", err))
		return ListStoryBookResponse{}, err
	}

	bookResponse := make([]domain.StoryBookCover, len(listStoryBook))
	var wg sync.WaitGroup
	for idx, book := range listStoryBook {
		wg.Add(1)
		go func(input *models.StoryBook) {
			defer wg.Done()
			imageUrl, err := g.s3.PreSignedGetObject(ctx, &storage.GetObjectRequest{
				Filename: book.TitleImagePath.String,
				Duration: 30,
			})
			if err != nil {
				g.logger.Hashcode(ctx).Error(fmt.Errorf("error get presigned url id %s: %w", book.ID, err))
			}
			bookResponse[idx] = domain.StoryBookCover{
				Id:            book.ID,
				CharacterName: book.CharacterName,
				Title:         book.Title,
				Theme:         book.Theme,
				TitleImageUrl: imageUrl.String(),
			}
		}(book)
	}
	wg.Wait()

	return ListStoryBookResponse{
		domain.ListStoryBookResponse{
			StoryBooks: bookResponse,
		},
	}, nil
}
