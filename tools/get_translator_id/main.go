// get_translator_id is a utility tool for extracting translator IDs from AnythingTranslate.
//
// This tool allows you to input a URL for a specific translator page and extracts
// the translator ID, which can then be used with the anythingtranslate library.
package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter translate URL (e.g., https://anythingtranslate.com/translators/pidgin-english-translator/) or 'exit' to quit: ")
		url, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}

		url = strings.TrimSpace(url)

		if strings.ToLower(url) == "exit" {
			fmt.Println("Exiting program.")
			break
		}

		if !strings.HasPrefix(url, "https://anythingtranslate.com/") {
			fmt.Fprintf(os.Stderr, "Invalid URL. URL must start with 'https://anythingtranslate.com/'\n")
			continue
		}

		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error making HTTP request: %v\n", err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading response body: %v\n", err)
			continue
		}

		re := regexp.MustCompile(`<input type="hidden" id="post-id" value="(\d+)">`)
		matches := re.FindStringSubmatch(string(body))

		if len(matches) < 2 {
			fmt.Fprintf(os.Stderr, "Could not find post-id in the HTML response\n")
			continue
		}

		postID := matches[1]
		fmt.Printf("Extracted post-id: %s\n\n", postID)
	}
}