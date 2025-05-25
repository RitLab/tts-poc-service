package utils

import (
	"regexp"
	"strings"
)

// SplitIntoSentences Splits text into sentences based on punctuation.
func SplitIntoSentences(text string) []string {
	re := regexp.MustCompile(`(?m)([^.!?]+[.!?])`)
	matches := re.FindAllString(text, -1)

	var sentences []string
	for _, match := range matches {
		sentence := strings.TrimSpace(match)
		if sentence != "" {
			sentences = append(sentences, sentence)
		}
	}
	return sentences
}

// ChunkSentences safely chunks large text input for RAG using sentence-aware logic.
func ChunkSentences(sentences []string, maxSize, sentenceOverlap int) []string {
	var chunks []string
	var currentChunk []string
	currentLen := 0

	for i := 0; i < len(sentences); {
		sentence := sentences[i]
		if currentLen+len(sentence)+1 <= maxSize {
			currentChunk = append(currentChunk, sentence)
			currentLen += len(sentence) + 1 // +1 for space or punctuation
			i++
		} else {
			if len(currentChunk) > 0 {
				chunks = append(chunks, strings.Join(currentChunk, " "))
			}

			// Start new chunk with overlap
			if sentenceOverlap > 0 && len(currentChunk[len(currentChunk)-sentenceOverlap:]) >= int(float64(maxSize)*0.7) {
				currentChunk = []string{}
				currentLen = 0
			} else if sentenceOverlap > 0 && len(currentChunk) >= sentenceOverlap {
				overlapChunk := currentChunk[len(currentChunk)-sentenceOverlap:]
				currentChunk = append([]string{}, overlapChunk...)
				currentLen = len(strings.Join(currentChunk, " "))
			} else {
				currentChunk = []string{}
				currentLen = 0
			}
		}
	}

	// Add remaining chunk
	if len(currentChunk) > 0 {
		chunks = append(chunks, strings.Join(currentChunk, " "))
	}

	return chunks
}

func IsValidChunk(text string) bool {
	words := strings.Fields(text)
	if len(words) < 10 || len(text) < 50 {
		return false
	}
	return true
}
