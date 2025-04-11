package shared

import "sync"

// LanguageManager is a singleton that holds the immutable list of languages.
type LanguageManager struct {
	languages []Language
}

var (
	instance  *LanguageManager
	once      sync.Once
	languages = []Language{
		{id: 1, culture: "en-US", name: "English", rtl: false},
		{id: 2, culture: "tr-TR", name: "Türkçe", rtl: false},
		{id: 3, culture: "de-DE", name: "Deutsch", rtl: false},
		{id: 4, culture: "fr-FR", name: "Français", rtl: false},
		{id: 5, culture: "es-ES", name: "Español", rtl: false},
		{id: 6, culture: "pt-PT", name: "Português", rtl: false},
		{id: 7, culture: "it-IT", name: "Italiano", rtl: false},
		{id: 8, culture: "el-GR", name: "Ελληνικά", rtl: false},
		{id: 9, culture: "ru-RU", name: "Русский", rtl: false},
		{id: 10, culture: "ja-JP", name: "日本語", rtl: false},
		{id: 11, culture: "zh-CN", name: "中文（简体", rtl: false},
		{id: 12, culture: "ar-SA", name: "العربية", rtl: false},
	}
)

// GetLanguageManager returns the singleton instance of LanguageManager.
// It initializes the list only once.
func GetLanguageManager() *LanguageManager {
	once.Do(func() {
		instance = &LanguageManager{languages: languages}
	})
	return instance
}

// GetLanguages returns a copy of the languages list to ensure immutability.
func (lm *LanguageManager) GetLanguages() []Language {
	copyList := make([]Language, len(lm.languages))
	copy(copyList, lm.languages)
	return copyList
}

// GetLanguageByCulture returns a pointer to the Language matching the given culture code.
// If not found, it returns en-US as default language.
func (lm *LanguageManager) GetLanguageByCulture(culture string) *Language {
	for _, lang := range lm.languages {
		if lang.culture == culture {
			// Returning the address of a loop variable is problematic because it gets reused.
			// To avoid this, we create a new variable for each match.
			l := lang
			return &l
		}
	}
	return &languages[0]
}
