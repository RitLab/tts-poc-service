package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"tts-poc-service/config"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/lib/connection"
	pkgError "tts-poc-service/pkg/common/error"
	"tts-poc-service/pkg/pdf/domain"
)

type openAIRepository struct {
	http connection.HttpConnectionInterface
	log  *baselogger.Logger
}

func NewOpenAIRepository(http connection.HttpConnectionInterface, log *baselogger.Logger) *openAIRepository {
	return &openAIRepository{
		http: http,
		log:  log,
	}
}

func (o *openAIRepository) SummarizeText(ctx context.Context, in *domain.OpenAIRequest) (string, error) {
	apiKey := config.Config.General.OpenAIKey
	if apiKey == "" {
		return "", fmt.Errorf(pkgError.APIKEY_NOT_VALID)
	}

	jsonData, err := json.Marshal(in)
	if err != nil {
		return "", err
	}

	headers := make(map[string]string)
	headers["Authorization"] = fmt.Sprintf("Bearer %s", apiKey)

	status, response, err := o.http.PostWithContext(ctx, config.Config.General.OpenAIEndpoint, headers, jsonData)
	if err != nil {
		return "", err
	}

	if status != http.StatusOK && status != http.StatusCreated {
		o.log.Hashcode(ctx).Error(fmt.Errorf("error response from OpenAI service: %d %s", status, string(response)))
		return "", fmt.Errorf(pkgError.API_CALL_ERROR)
	}

	out := domain.OpenAIResponse{}
	json.Unmarshal(response, &out)
	if len(out.Choices) > 0 {
		return out.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf(pkgError.API_CALL_ERROR)
}
