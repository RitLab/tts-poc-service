package query

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/lib/storage"
	"tts-poc-service/pkg/common/decorator"
	"tts-poc-service/pkg/story/domain"
)

type GetStoryBookQuery struct {
	Id string `json:"id" validate:"required"`
}

type GetStoryBookResponse struct {
	domain.BookResponse
}

type GetStoryBookHandler decorator.QueryHandler[GetStoryBookQuery, GetStoryBookResponse]

type getStoryBookRepository struct {
	repo   domain.QueryRepository
	s3     storage.Storage
	logger *baselogger.Logger
}

func NewGetStoryBookRepository(repo domain.QueryRepository, s3 storage.Storage,
	log *baselogger.Logger) decorator.QueryHandler[GetStoryBookQuery, GetStoryBookResponse] {
	return decorator.ApplyQueryDecorators[GetStoryBookQuery, GetStoryBookResponse](
		getStoryBookRepository{repo: repo, s3: s3, logger: log},
		log)
}

func (g getStoryBookRepository) Handle(ctx context.Context, in GetStoryBookQuery) (GetStoryBookResponse, error) {
	getStoryBook, err := g.repo.GetStoryBookById(ctx, in.Id)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("get story book failed: %w", err))
		return GetStoryBookResponse{}, err
	}
	storyPages := make([]domain.PageDatabase, 0)
	err = json.Unmarshal(getStoryBook.Story.JSON, &storyPages)
	if err != nil {
		g.logger.Hashcode(ctx).Error(fmt.Errorf("unmarshal story book failed: %w", err))
		return GetStoryBookResponse{}, err
	}

	bookResponse := make([]domain.PageResponse, len(storyPages))
	var wg sync.WaitGroup
	for idx, book := range storyPages {
		wg.Add(1)
		go func(input domain.PageDatabase) {
			defer wg.Done()
			imageUrl, err := g.s3.PreSignedGetObject(ctx, &storage.GetObjectRequest{
				Filename: book.ImagePath,
				Duration: 30,
			})
			if err != nil {
				g.logger.Hashcode(ctx).Error(fmt.Errorf("error get presigned url id: %s page number %d: %w", in.Id, book.PageNumber, err))
			}
			audioUrl, err := g.s3.PreSignedGetObject(ctx, &storage.GetObjectRequest{
				Filename: book.AudioPath,
				Duration: 30,
			})
			if err != nil {
				g.logger.Hashcode(ctx).Error(fmt.Errorf("error get presigned url id: %s page number %d: %w", in.Id, book.PageNumber, err))
			}
			bookResponse[idx] = domain.PageResponse{
				Page: domain.Page{
					PageNumber: book.PageNumber,
					Body:       book.Body,
				},
				ImageUrl: imageUrl.String(),
				AudioUrl: audioUrl.String(),
			}
		}(book)
	}
	wg.Wait()

	return GetStoryBookResponse{
		domain.BookResponse{
			Id:    getStoryBook.ID,
			Pages: bookResponse,
		},
	}, nil
}
