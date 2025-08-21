package domain

type GenerateBook struct {
	TotalPages    int    `json:"total_pages"`
	CharacterName string `json:"character_name"`
	Context       string `json:"context"`
}

type GeneratedAIPage struct {
	PageNumber   int    `json:"page_number"`
	Body         string `json:"body"`
	Illustration string `json:"illustration"`
}

type Page struct {
	PageNumber int    `json:"page_number"`
	Body       string `json:"body"`
}

type Book struct {
	Pages []GeneratedAIPage `json:"pages"`
}

type PageDatabase struct {
	Page
	ImagePath string `json:"image_path"`
	AudioPath string `json:"audio_path"`
}

type PageResponse struct {
	Page
	ImageUrl string `json:"image_url"`
	AudioUrl string `json:"audio_url"`
}

type BookResponse struct {
	Id    string         `json:"id"`
	Pages []PageResponse `json:"pages"`
}

var languageMapping = map[string]string{
	"English":    "en",
	"Indonesian": "id",
	"Japanese":   "ja",
	"German":     "de",
	"Korean":     "ko",
}

func LanguageMapping(inputLanguage string) string {
	if _, ok := languageMapping[inputLanguage]; ok {
		return languageMapping[inputLanguage]
	}
	return "en-US"
}

var voiceMapping = map[string]string{
	"Hero":     "Vindemiatrix",
	"Life":     "Kore",
	"Nature":   "Sadaltager",
	"Fantasy":  "Fenrir",
	"Space":    "Laomedeia",
	"Values":   "Gacrux",
	"Cultural": "Zephyr",
}

func VoiceMapping(inputVoice string) string {
	if _, ok := voiceMapping[inputVoice]; ok {
		return voiceMapping[inputVoice]
	}
	return "Kore"
}

var themeMapping = map[string]string{
	"Hero":     "Main character is a hero with kindness and empathy",
	"Life":     "Main character everyday life and growing up",
	"Nature":   "Main character love caring the environment and taking care animals",
	"Fantasy":  "Main character discover magical world with friendly monster",
	"Space":    "Main character discover outer space with high end technology",
	"Values":   "Main character is kindness, empathy, courage, bravery and love to help others",
	"Cultural": "Main character love to attend cultural festivals and tradition",
}

func ThemeMapping(inputTheme string) string {
	if _, ok := themeMapping[inputTheme]; ok {
		return themeMapping[inputTheme]
	}
	return "Main character everyday life and growing up"
}

type StoryBookCover struct {
	Id            string `json:"id"`
	CharacterName string `json:"character_name"`
	Title         string `json:"title"`
	TitleImageUrl string `json:"title_image_url"`
	Theme         string `json:"theme"`
}

type ListStoryBookResponse struct {
	StoryBooks []StoryBookCover `json:"story_books"`
}
