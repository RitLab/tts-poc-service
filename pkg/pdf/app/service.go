package app

import (
	"database/sql"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/lib/database"
	"tts-poc-service/lib/gemini_ai"
	"tts-poc-service/lib/storage"
	"tts-poc-service/pkg/pdf/adapters"
	"tts-poc-service/pkg/pdf/app/command"
	"tts-poc-service/pkg/pdf/app/query"
)

type PdfService struct {
	Queries  PdfQueries
	Commands PdfCommands
}

type PdfQueries struct {
	JoinPdfFilesHandler    query.JoinPdfFilesHandler
	SignPdfFileHandler     query.SignPdfFileHandler
	SummarizePdfHandler    query.SummarizePdfHandler
	GenerateContextHandler query.GenerateContextHandler
}

type PdfCommands struct {
	VerifyPdfFileHandler command.VerifyPdfFileHandler
	UploadContextPdf     command.UpdateContextPdfHandler
}

func NewPdfService(log *baselogger.Logger, s3 storage.Storage, db *sql.DB, ai gemini_ai.GenAIMethod, dbVector database.VectorDatabase) PdfService {
	pdfAdapter := adapters.NewPdfRepository(db, log)
	return PdfService{
		Queries: PdfQueries{
			JoinPdfFilesHandler:    query.NewJoinPdfFilesRepository(s3, log),
			SignPdfFileHandler:     query.NewSignPdfFilesRepository(pdfAdapter, pdfAdapter, s3, log),
			SummarizePdfHandler:    query.NewSummarizePdfRepository(ai, log),
			GenerateContextHandler: query.NewGenerateContextRepository(dbVector, ai, log),
		},
		Commands: PdfCommands{
			VerifyPdfFileHandler: command.NewVerifyPdfFilesRepository(pdfAdapter, log),
			UploadContextPdf:     command.NewUpdateContextPdfRepository(dbVector, ai, log),
		},
	}
}
