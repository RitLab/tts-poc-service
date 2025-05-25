package hook_models

import (
	"context"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"tts-poc-service/models"
)

func addSupportIdBeforeHook(ctx context.Context, exec boil.ContextExecutor, s *models.Support) error {
	s.ID = uuid.NewString()
	return nil
}

func init() {
	models.AddSupportHook(boil.BeforeInsertHook, addSupportIdBeforeHook)
}
