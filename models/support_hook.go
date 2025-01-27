package models

import (
	"context"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func addSupportIdBeforeHook(ctx context.Context, exec boil.ContextExecutor, s *Support) error {
	s.ID = uuid.NewString()
	return nil
}

func init() {
	AddSupportHook(boil.BeforeInsertHook, addSupportIdBeforeHook)
}
