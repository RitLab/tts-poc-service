package domain

import (
	"code.sajari.com/docconv/v2"
	"github.com/jdkato/prose/v2"
	"io"
	"strings"
)

type PdfReader struct {
	text string
}

func NewPdfReader(r io.Reader) (*PdfReader, error) {
	text, _, err := docconv.ConvertPDF(r)
	if err != nil {
		return nil, err
	}
	return &PdfReader{text: text}, nil
}

func (p *PdfReader) replaceSpace() string {
	var builder strings.Builder
	for _, char := range p.text {
		if char >= 32 || char == '\n' {
			builder.WriteRune(char)
		}
	}
	return builder.String()
}

func (p *PdfReader) tokenize(text string) (string, error) {
	newDoc, err := prose.NewDocument(text)
	if err != nil {
		return "", err
	}
	var sentences []string
	for _, sent := range newDoc.Sentences() {
		tokens := strings.Fields(sent.Text)
		if len(tokens) > 0 {
			sentences = append(sentences, strings.Join(tokens, " "))
		}
	}
	return strings.Join(sentences, " "), nil
}

func (p *PdfReader) CleanText() (string, error) {
	text := p.replaceSpace()
	return p.tokenize(text)
}

func (p *PdfReader) chunkText(text string, chunkSize int) []string {
	if chunkSize < 1 {
		chunkSize = 1000
	}
	words := strings.Fields(text)
	var chunks []string

	for i := 0; i < len(words); i += chunkSize {
		end := i + chunkSize
		if end > len(words) {
			end = len(words)
		}
		chunk := strings.Join(words[i:end], " ")
		chunks = append(chunks, chunk)
	}

	return chunks
}
