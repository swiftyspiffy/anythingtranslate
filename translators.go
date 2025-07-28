// Package anythingtranslate provides a client for the AnythingTranslate service.
package anythingtranslate

// TranslatorTarget represents a specific translator available in the AnythingTranslate service.
// Each translator has a unique ID that is used to identify it in API requests.
type TranslatorTarget string

// Predefined translator targets available in the AnythingTranslate service.
// These constants can be used with the Translate method to specify which translator to use.
const (
	// Target_GenAlphaSlang translates text into Gen Alpha slang.
	Target_GenAlphaSlang TranslatorTarget = "17491"

	// Target_FancyOldEnglish translates text into fancy old English.
	Target_FancyOldEnglish TranslatorTarget = "29293"

	// Target_AdvancedEnglish translates text into advanced English with more sophisticated vocabulary.
	Target_AdvancedEnglish TranslatorTarget = "19036"

	// Target_ScientificallyInaccurateAndFunny translates text into scientifically inaccurate and humorous language.
	Target_ScientificallyInaccurateAndFunny TranslatorTarget = "20496"

	// Target_Spanish translates text into Spanish.
	Target_Spanish TranslatorTarget = "16120"

	// Target_Medieval translates text into medieval English.
	Target_Medieval TranslatorTarget = "29922"

	// Target_PidginEnglish translates text into Pidgin English.
	Target_PidginEnglish TranslatorTarget = "19412"

	// Target_Warhammer40kOrk translates text into the dialect of Orks from Warhammer 40k.
	Target_Warhammer40kOrk TranslatorTarget = "81700"

	// Target_LazyEbonics translates text into Ebonics with a lazy tone.
	Target_LazyEbonics TranslatorTarget = "43765"

	// Target_Patois translates text into Patois dialect.
	Target_Patois TranslatorTarget = "23034"

	// Target_Mafioso translates text into Mafioso-style language.
	Target_Mafioso TranslatorTarget = "105535"
)
