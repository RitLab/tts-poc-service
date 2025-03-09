package app

import (
	"database/sql"
	"tts-poc-service/lib/baselogger"
	"tts-poc-service/lib/connection"
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
	JoinPdfFilesHandler query.JoinPdfFilesHandler
	SignPdfFileHandler  query.SignPdfFileHandler
	SummarizePdfHandler query.SummarizePdfHandler
}

type PdfCommands struct {
	VerifyPdfFileHandler command.VerifyPdfFileHandler
}

func NewPdfService(log *baselogger.Logger, s3 storage.Storage, db *sql.DB, httpCon connection.HttpConnectionInterface) PdfService {
	pdfAdapter := adapters.NewPdfRepository(db, log)
	openAIAdapter := adapters.NewOpenAIRepository(httpCon, log)
	return PdfService{
		Queries: PdfQueries{
			JoinPdfFilesHandler: query.NewJoinPdfFilesRepository(s3, log),
			SignPdfFileHandler:  query.NewSignPdfFilesRepository(pdfAdapter, pdfAdapter, s3, log),
			SummarizePdfHandler: query.NewSummarizePdfRepository(openAIAdapter, s3, log),
		},
		Commands: PdfCommands{
			VerifyPdfFileHandler: command.NewVerifyPdfFilesRepository(pdfAdapter, log),
		},
	}
}
