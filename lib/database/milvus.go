package database

import (
	"context"
	"fmt"
	"strings"
	"tts-poc-service/config"
	"tts-poc-service/lib/baselogger"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type Embedding struct {
	Id    int64
	Value [][]float32 `json:"value"`
	Chunk []string    `json:"chunk"`
}

type MilvusValue struct {
}

type VectorDatabase interface {
	CreateEmbeddedCollection(ctx context.Context, collectionName string) error
	StoreEmbedding(ctx context.Context, payload Embedding, collectionName string) error
	SearchEmbedding(ctx context.Context, collectionName string, queryEmbedding []float32, topK int) (string, error)
	CheckSimilarity(ctx context.Context, collectionName string, queryEmbedding []float32) error
}

type vectorClient struct {
	client.Client
}

func NewMilvusClient(ctx context.Context, log *baselogger.Logger) VectorDatabase {
	c, err := client.NewClient(ctx, client.Config{
		Address: config.Config.General.MilvusAddress,
	})
	if err != nil {
		log.Panic(err)
	}
	return &vectorClient{c}
}

func (v *vectorClient) CreateEmbeddedCollection(ctx context.Context, collectionName string) error {
	has, err := v.HasCollection(ctx, collectionName)
	if err != nil {
		return err
	}
	if has {
		return nil
	}

	// Define the schema for the collection.  The schema needs to align with your data.
	idField := entity.NewField().WithName("id").WithDescription("primary key id").
		WithDataType(entity.FieldTypeInt64).WithIsPrimaryKey(true).WithIsAutoID(true)
	embeddingField := entity.NewField().WithName("embedding").WithDescription("embedding vector").
		WithDataType(entity.FieldTypeFloatVector).WithDim(config.Config.General.EmbeddingDimension)
	textField := entity.NewField().WithName("text").WithDescription("original text").WithDataType(entity.FieldTypeVarChar).WithMaxLength(1000)
	schema := entity.NewSchema().WithName(collectionName).WithDynamicFieldEnabled(true).
		WithField(idField).WithField(embeddingField).WithField(textField)

	err = v.CreateCollection(ctx, schema, 2)
	if err != nil {
		return err
	}

	// Create an index for efficient similarity search.
	idx, err := entity.NewIndexIvfFlat(entity.COSINE, 1024)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	if err = v.CreateIndex(ctx, collectionName, "embedding", idx, false); err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	if err = v.LoadCollection(ctx, collectionName, false); err != nil {
		return fmt.Errorf("failed to load collection: %w", err)
	}
	return nil
}

func (v *vectorClient) StoreEmbedding(ctx context.Context, payload Embedding, collectionName string) error {
	if len(payload.Value) != len(payload.Chunk) {
		return fmt.Errorf("number of embeddings and text chunks do not match")
	}

	embeddingFields := make([][]float32, 0, len(payload.Value))
	textFields := make([]string, 0, len(payload.Chunk))

	for i := 0; i < len(payload.Value); i++ {
		embeddingFields = append(embeddingFields, payload.Value[i])
		textFields = append(textFields, payload.Chunk[i])
	}

	embeddingColumn := entity.NewColumnFloatVector("embedding", int(config.Config.General.EmbeddingDimension), embeddingFields)
	textColumn := entity.NewColumnVarChar("text", textFields)

	insertResult, err := v.Insert(ctx, collectionName, "", embeddingColumn, textColumn)
	if err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
	}

	if insertResult.Len() != len(payload.Value) {
		return fmt.Errorf("failed to insert all data: inserted %d, expected %d", insertResult.Len(), len(payload.Value))
	}

	return nil
}

func (v *vectorClient) SearchEmbedding(ctx context.Context, collectionName string, queryEmbedding []float32, topK int) (string, error) {
	searchParam, err := entity.NewIndexFlatSearchParam()
	if err != nil {
		return "", fmt.Errorf("failed to create search params: %w", err)
	}

	results, err := v.Search(
		ctx,
		collectionName,
		[]string{},
		"",
		[]string{"text"},
		[]entity.Vector{entity.FloatVector(queryEmbedding)},
		"embedding",   // Use the "embedding" field for the search.
		entity.COSINE, // Use cosine similarity.
		topK,
		searchParam,
	)
	if err != nil {
		return "", fmt.Errorf("failed to search: %w", err)
	}

	if len(results) == 0 || results[0].ResultCount == 0 {
		return "", nil
	}

	contextText := make([]string, 0, topK)
	for i := 0; i < topK; i++ {
		chunkText, _ := results[0].Fields.GetColumn("text").GetAsString(i)
		contextText = append(contextText, chunkText)
	}
	return strings.Join(contextText, " "), nil
}

func (v *vectorClient) CheckSimilarity(ctx context.Context, collectionName string, queryEmbedding []float32) error {
	searchParam, err := entity.NewIndexFlatSearchParam()
	if err != nil {
		return fmt.Errorf("failed to create search params: %w", err)
	}

	results, err := v.Search(
		ctx,
		collectionName,
		[]string{},
		"",
		[]string{"id"},
		[]entity.Vector{entity.FloatVector(queryEmbedding)},
		"embedding",   // Use the "embedding" field for the search.
		entity.COSINE, // Use cosine similarity.
		1,
		searchParam,
	)
	if err != nil {
		return fmt.Errorf("failed to search: %w", err)
	}

	if len(results) == 0 || results[0].Scores[0] > 0 {
		score := results[0].Scores[0]
		if score >= 0.98 {
			return fmt.Errorf("duplicate chunk found with ID: %s", results[0].IDs.FieldData().String())
		}
	}
	return nil
}
