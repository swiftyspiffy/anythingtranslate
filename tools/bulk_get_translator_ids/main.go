// bulk_get_translator_ids is a utility tool for extracting multiple translator IDs from AnythingTranslate.
//
// This tool fetches translator IDs for a predefined list of translators and formats them
// as Go constants that can be added to the translators.go file in the anythingtranslate library.
package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	// List of translator URL slugs to fetch IDs for
	// Uncomment additional translators as needed
	translators := []string{
		/*
			"pidgin-english-translator",
			"warhammer-40k-ork-translator",
			"lazy-ebonics-translator",
			"patois-translator",
			"mafioso-translator",
		*/
		"gen-alpha-slang-translator",
		"fancy-old-english-translator",
		"advanced-english-translator",
		"scientifically-inaccurate-and-funny-translator",
		"spanish-translator",
		"medieval-translator",
	}

	fmt.Println("Fetching translator IDs...")
	fmt.Println()

	for _, translator := range translators {
		url := fmt.Sprintf("https://anythingtranslate.com/translators/%s/", translator)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error making HTTP request for %s: %v\n", translator, err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			fmt.Printf("Error reading response body for %s: %v\n", translator, err)
			continue
		}

		re := regexp.MustCompile(`<input type="hidden" id="post-id" value="(\d+)">`)
		matches := re.FindStringSubmatch(string(body))

		if len(matches) < 2 {
			fmt.Printf("Could not find post-id in the HTML response for %s\n", translator)
			continue
		}

		postID := matches[1]
		constantName := formatTranslatorName(translator)

		fmt.Printf("Target_%s TranslatorTarget = \"%s\"\n", constantName, postID)
	}
}

// formatTranslatorName converts a hyphenated translator slug to a camel case constant name.
// For example, "fancy-old-english-translator" becomes "FancyOldEnglish".
func formatTranslatorName(translator string) string {
	// Remove the "-translator" suffix
	translator = strings.TrimSuffix(translator, "-translator")
	
	// Split by hyphens
	parts := strings.Split(translator, "-")

	// Capitalize each part
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}

	// Join the parts without separators
	return strings.Join(parts, "")
}