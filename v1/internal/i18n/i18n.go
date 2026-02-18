package i18n

import (
	"embed"
	"fmt"
	"sync"

	"gopkg.in/yaml.v3"
)

//go:embed translations/*.yaml
var translationFiles embed.FS

var (
	mu          sync.RWMutex
	loaded      *Translations
	loadedLang  string
)

// Load parses the translation file for the given language code and stores
// it as the active translations. Must be called once at startup before
// any call to T().
func Load(lang string) error {
	path := fmt.Sprintf("translations/%s.yaml", lang)
	data, err := translationFiles.ReadFile(path)
	if err != nil {
		return fmt.Errorf("language %q not available: %w", lang, err)
	}

	var t Translations
	if err := yaml.Unmarshal(data, &t); err != nil {
		return fmt.Errorf("parse %s translations: %w", lang, err)
	}

	if err := validate(&t, lang); err != nil {
		return err
	}

	mu.Lock()
	loaded = &t
	loadedLang = lang
	mu.Unlock()

	return nil
}

// T returns the active translations. Panics if Load has not been called.
func T() *Translations {
	mu.RLock()
	t := loaded
	mu.RUnlock()
	if t == nil {
		panic("i18n.T() called before i18n.Load()")
	}
	return t
}

// Lang returns the currently loaded language code.
func Lang() string {
	mu.RLock()
	defer mu.RUnlock()
	return loadedLang
}

// Available returns the list of available language codes by scanning
// the embedded translation files.
func Available() []string {
	entries, err := translationFiles.ReadDir("translations")
	if err != nil {
		return nil
	}
	var langs []string
	for _, e := range entries {
		name := e.Name()
		if len(name) > 5 && name[len(name)-5:] == ".yaml" {
			langs = append(langs, name[:len(name)-5])
		}
	}
	return langs
}

// TranslateBio looks up a key in the given translation map.
// Returns the translated string if found, or the original key as fallback.
func TranslateBio(m map[string]string, key string) string {
	if v, ok := m[key]; ok {
		return v
	}
	return key
}

// validate checks that critical fields are populated.
func validate(t *Translations, lang string) error {
	if t.Consciousness.SystemPrompt == "" {
		return fmt.Errorf("%s: missing consciousness.system_prompt", lang)
	}
	if t.Sense.Fallback == "" {
		return fmt.Errorf("%s: missing sense.fallback", lang)
	}
	if t.Reviewer.SystemPrompt == "" {
		return fmt.Errorf("%s: missing reviewer.system_prompt", lang)
	}
	if len(t.Sense.Keywords.Thermal) == 0 {
		return fmt.Errorf("%s: missing sense.keywords.thermal", lang)
	}
	if t.Defaults.SelfNarrative == "" {
		return fmt.Errorf("%s: missing defaults.self_narrative", lang)
	}
	return nil
}
