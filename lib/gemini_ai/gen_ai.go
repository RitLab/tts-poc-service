package gemini_ai

import (
	"context"
	"fmt"
	"google.golang.org/genai"
	"tts-poc-service/config"
)

var embedModel = "text-embedding-004"

type GenAIMethod interface {
	SummarizeText(ctx context.Context, pdfBytes string) (string, error)
	GetTranscriptAudio(ctx context.Context, filepath string) (string, error)
	TextEmbedding(ctx context.Context, chunk string) ([]float32, error)
	GenerateFromContext(ctx context.Context, contextText, question string) (string, error)
}

type genAI struct {
	c *genai.Client
}

func NewGenAI(ctx context.Context) GenAIMethod {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.Config.General.GeminiAPIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return &genAI{}
	}
	return &genAI{c: client}
}

func (g *genAI) SummarizeText(ctx context.Context, pdfBytes string) (string, error) {
	parts := []*genai.Part{
		&genai.Part{
			Text: pdfBytes,
		},
		genai.NewPartFromText("Summarize this document"),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	result, err := g.c.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		contents,
		nil,
	)
	if err != nil {
		return "", err
	}
	return result.Text(), nil
}

func (g *genAI) GetTranscriptAudio(ctx context.Context, filepath string) (string, error) {
	uploadedFile, _ := g.c.Files.UploadFromPath(
		ctx,
		filepath,
		nil,
	)

	parts := []*genai.Part{
		genai.NewPartFromText("Generate a transcript of the speech."),
		genai.NewPartFromURI(uploadedFile.URI, uploadedFile.MIMEType),
	}
	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	result, err := g.c.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		contents,
		nil,
	)
	if err != nil {
		return "", err
	}
	return result.Text(), nil
}

func (g *genAI) TextEmbedding(ctx context.Context, chunk string) ([]float32, error) {
	parts := []*genai.Part{genai.NewPartFromText(chunk)}
	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	result, err := g.c.Models.EmbedContent(ctx, embedModel, contents, &genai.EmbedContentConfig{TaskType: "SEMANTIC_SIMILARITY"})
	if err != nil {
		return nil, err
	}
	return result.Embeddings[0].Values, nil
}

func (g *genAI) GenerateFromContext(ctx context.Context, contextText, question string) (string, error) {
	prompt := fmt.Sprintf("Answer the following question based on the provided context:\n\nContext:\n%s\n", contextText)
	contentContext := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(prompt, genai.RoleUser),
	}

	result, err := g.c.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text(question),
		contentContext,
	)
	if err != nil {
		return "", err
	}
	return result.Text(), nil
}
