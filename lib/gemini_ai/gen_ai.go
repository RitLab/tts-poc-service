package gemini_ai

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/genai"
	"os"
	"tts-poc-service/config"
	"tts-poc-service/pkg/common/constant"
	"tts-poc-service/pkg/common/utils"
	"tts-poc-service/pkg/story/domain"
)

var embedModel = "text-embedding-004"

type GenAIMethod interface {
	SummarizeText(ctx context.Context, pdfBytes string) (string, error)
	GetTranscriptAudio(ctx context.Context, filepath string) (string, error)
	SummarizeYoutubeUrl(ctx context.Context, youtubeUrl string) (string, error)
	TranscriptYoutubeUrl(ctx context.Context, youtubeUrl string) (string, error)
	TextEmbedding(ctx context.Context, chunk string) ([]float32, error)
	GenerateFromContext(ctx context.Context, contextText, question string) (string, error)
	GenerateStory(ctx context.Context, story *domain.GenerateBook, language, theme string) (*domain.Book, error)
	GenerateImageStory(ctx context.Context, story domain.GeneratedAIPage, id string) (string, error)
	GenerateAudioStory(ctx context.Context, story domain.GeneratedAIPage, id, languageCode, voice string) (string, error)
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

func (g *genAI) SummarizeYoutubeUrl(ctx context.Context, youtubeUrl string) (string, error) {
	parts := []*genai.Part{
		&genai.Part{
			FileData: &genai.FileData{FileURI: youtubeUrl, MIMEType: "video/*"},
		},
		genai.NewPartFromText("Summarize the video"),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	result, err := g.c.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-preview-05-20",
		contents,
		nil,
	)
	if err != nil {
		return "", err
	}
	return result.Text(), nil
}

func (g *genAI) TranscriptYoutubeUrl(ctx context.Context, youtubeUrl string) (string, error) {
	parts := []*genai.Part{
		&genai.Part{
			FileData: &genai.FileData{FileURI: youtubeUrl, MIMEType: "video/*"},
		},
		genai.NewPartFromText("Transcribe the video"),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	result, err := g.c.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-preview-05-20",
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

func (g *genAI) GenerateStory(ctx context.Context, story *domain.GenerateBook, language, theme string) (*domain.Book, error) {
	jsonFormat := `{
		  "pages": [
			{
			  "page_number": 1,
			  "body": "string",
			  "illustration": "string"
			},
			{
			  "page_number": 1,
			  "body": "string",
			  "illustration": "string"
			}
		  ]
		}`

	prompt := fmt.Sprintf("Help to create a book story with the name of main character is %s with maximum %d pages.\nTheme of this story is about %s.\nCreate story in %s language with every pages have 1 sentences less than 100 character and the first page is the title of the story.\nDo not attach any text on the image.\n\nGive the result in this json format %s",
		story.CharacterName, story.TotalPages, theme, language, jsonFormat)
	contentContext := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(prompt, genai.RoleUser),
	}

	result, err := g.c.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text(story.Context),
		contentContext,
	)
	if err != nil {
		return nil, err
	}

	textResponse := result.Text()
	jsonTextResponse := utils.GetJsonResponseOnly(textResponse)

	var book *domain.Book
	err = json.Unmarshal([]byte(jsonTextResponse), &book)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (g *genAI) GenerateImageStory(ctx context.Context, story domain.GeneratedAIPage, id string) (string, error) {
	prompt := fmt.Sprintf("generate image with this illustration: %s with children cartoon style", story.Illustration)
	imageConfig := &genai.GenerateContentConfig{
		ResponseModalities: []string{"TEXT", "IMAGE"},
	}
	result, err := g.c.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash-preview-image-generation",
		genai.Text(prompt),
		imageConfig,
	)
	if err != nil {
		return "", err
	}

	fname := fmt.Sprintf("%s/image-%s-%d.png", constant.IMAGE_FOLDER, id, story.PageNumber)
	for _, part := range result.Candidates[0].Content.Parts {
		if part.InlineData == nil {
			continue
		}
		err = os.WriteFile(fname, part.InlineData.Data, 777)
		if err != nil {
			continue
		}
		break
	}
	return fname, nil
}

func (g *genAI) GenerateAudioStory(ctx context.Context, story domain.GeneratedAIPage, id, languageCode, voice string) (string, error) {
	prompt := fmt.Sprintf("Generate audio based on given text")
	contentContext := &genai.GenerateContentConfig{
		ResponseModalities: []string{"AUDIO"},
		SystemInstruction:  genai.NewContentFromText(prompt, genai.RoleUser),
		SpeechConfig: &genai.SpeechConfig{
			VoiceConfig:  &genai.VoiceConfig{PrebuiltVoiceConfig: &genai.PrebuiltVoiceConfig{VoiceName: voice}},
			LanguageCode: languageCode,
		},
	}
	result, err := g.c.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-preview-tts",
		genai.Text(story.Body),
		contentContext,
	)
	if err != nil {
		return "", err
	}
	fname := fmt.Sprintf("%s/audio-%s-%d.wav", constant.AUDIO_FOLDER, id, story.PageNumber)
	for _, part := range result.Candidates[0].Content.Parts {
		if part.InlineData == nil {
			continue
		}
		err = os.WriteFile(fname, part.InlineData.Data, 777)
		if err != nil {
			continue
		}
		break
	}
	return fname, nil
}
