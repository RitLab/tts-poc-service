package htgo

import (
	"fmt"
	htgotts "github.com/Ritlab/htgo-tts"
	handlers "github.com/Ritlab/htgo-tts/handlers"
	"github.com/google/uuid"
	"strings"
	"tts-poc-service/pkg/common/constant"
	"tts-poc-service/pkg/common/utils"
)

type PlayerInterface interface {
	Save(text, lang string) ([]string, error)
	SaveOneFile(text, lang, id string, pageNumber int) (string, error)
	Play(text, lang string) error
}

type Player struct {
}

func (p *Player) Save(text, lang string) ([]string, error) {
	speech := htgotts.Speech{
		Folder:   constant.AUDIO_FOLDER,
		Language: lang,
		Handler:  &handlers.Native{},
	}
	out := make([]string, 0)
	if len(text) > 100 {
		splitSentences := strings.Split(text, ".")
		totalIter := make([]string, 0)
		for _, sentence := range splitSentences {
			if len(sentence) > 100 {
				splitText := strings.Split(sentence, " ")

				totalText := 0
				eachIter := make([]string, 0)
				for i := range splitText {
					if i == len(splitText)-1 {
						eachIter = append(eachIter, splitText[i])
						newIter := strings.Join(eachIter, " ")
						totalIter = append(totalIter, newIter)
					} else if totalText+len(splitText[i]) < 100 {
						eachIter = append(eachIter, splitText[i])
						totalText += len(splitText[i])
					} else {
						newIter := strings.Join(eachIter, " ")
						totalIter = append(totalIter, newIter)

						eachIter = make([]string, 0)
						eachIter = append(eachIter, splitText[i])
						totalText = len(splitText[i])
					}
				}
			} else {
				totalIter = append(totalIter, sentence)
			}
		}

		for i := range totalIter {
			file, err := speech.CreateSpeechFile(totalIter[i], uuid.NewString())
			if err != nil {
				return nil, err
			}
			out = append(out, file)
		}
	} else {
		file, err := speech.CreateSpeechFile(text, uuid.NewString())
		if err != nil {
			return nil, err
		}
		out = append(out, file)
	}

	return out, nil
}

func (p *Player) SaveOneFile(text, lang, id string, pageNumber int) (string, error) {
	speech := htgotts.Speech{
		Folder:   constant.AUDIO_FOLDER,
		Language: lang,
		Handler:  &handlers.Native{},
	}
	textWithoutMarkdown := utils.StripMarkdown(text)
	file, err := speech.CreateSpeechFile(textWithoutMarkdown, fmt.Sprintf("audio-%s-%d", id, pageNumber))
	if err != nil {
		return "", err
	}

	return file, nil
}

func (p *Player) Play(text, lang string) error {
	speech := htgotts.Speech{
		Folder:   constant.AUDIO_FOLDER,
		Language: lang,
		Handler:  &handlers.Native{},
	}

	if len(text) > 100 {
		splitSentences := strings.Split(text, ".")
		totalIter := make([]string, 0)
		for _, sentence := range splitSentences {
			if len(sentence) > 100 {
				splitText := strings.Split(text, " ")

				totalText := 0
				eachIter := make([]string, 0)
				for i := range splitText {
					if i == len(splitText)-1 {
						eachIter = append(eachIter, splitText[i])
						newIter := strings.Join(eachIter, " ")
						totalIter = append(totalIter, newIter)
					} else if totalText+len(splitText[i]) < 100 {
						eachIter = append(eachIter, splitText[i])
						totalText += len(splitText[i])
					} else {
						newIter := strings.Join(eachIter, " ")
						totalIter = append(totalIter, newIter)

						eachIter = make([]string, 0)
						eachIter = append(eachIter, splitText[i])
						totalText = len(splitText[i])
					}
				}
			} else {
				totalIter = append(totalIter, sentence)
			}
		}

		for i := range totalIter {
			err := speech.Speak(totalIter[i])
			if err != nil {
				return err
			}
		}
	} else {
		err := speech.Speak(text)
		if err != nil {
			return err
		}
	}

	return nil
}
