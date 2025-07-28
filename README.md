# AnythingTranslate Go Client

[![Go Reference](https://pkg.go.dev/badge/github.com/swiftyspiffy/anythingtranslate.svg)](https://pkg.go.dev/github.com/swiftyspiffy/anythingtranslate)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go client library for the [AnythingTranslate](https://anythingtranslate.com/) service. This library allows you to programmatically translate text using various fun and creative translators.

## Features

- Simple, easy-to-use API
- Support for all AnythingTranslate translators
- Custom logging interface
- Automatic session management
- Comprehensive error handling

## Installation

```bash
go get github.com/swiftyspiffy/anythingtranslate
```

## Quick Start

```go
package main

import (
	"fmt"
	"log"

	"github.com/swiftyspiffy/anythingtranslate"
)

func main() {
	// Create a new client
	client := anythingtranslate.NewAnythingTranslateClient()

	// Translate text using a predefined translator
	translation, err := client.Translate("Hello, world!", anythingtranslate.Target_Warhammer40kOrk)
	if err != nil {
		log.Fatalf("Translation failed: %v", err)
	}

	fmt.Println(translation)
}
```

## Available Translators

The library provides constants for various translators:

- `Target_GenAlphaSlang` - Gen Alpha slang
- `Target_FancyOldEnglish` - Fancy old English
- `Target_AdvancedEnglish` - Advanced English with sophisticated vocabulary
- `Target_ScientificallyInaccurateAndFunny` - Scientifically inaccurate and humorous language
- `Target_Spanish` - Spanish
- `Target_Medieval` - Medieval English
- `Target_PidginEnglish` - Pidgin English
- `Target_Warhammer40kOrk` - Warhammer 40k Ork dialect
- `Target_LazyEbonics` - Ebonics with a lazy tone
- `Target_Patois` - Patois dialect
- `Target_Mafioso` - Mafioso-style language

## Custom Logging

You can provide your own logger implementation:

```go
type MyLogger struct{}

func (l MyLogger) Debug(msg string, fields ...anythingtranslate.Field) {
	// Custom debug logging implementation
}

func (l MyLogger) Info(msg string, fields ...anythingtranslate.Field) {
	// Custom info logging implementation
}

// Create a client with custom logger
client := anythingtranslate.NewAnythingTranslateClientWithLogger(MyLogger{})
```

## Using Custom Translator IDs

If you have a translator ID that's not included in the predefined constants:

```go
translation, err := client.TranslateWithCustomTarget("Hello, world!", "12345")
```

## Tools

The repository includes tools to help you find translator IDs:

- `tools/get_translator_id` - Get the ID for a single translator
- `tools/bulk_get_translator_ids` - Get IDs for multiple translators

To run these tools:

```bash
# Get a single translator ID
cd tools/get_translator_id
go run main.go

# Get multiple translator IDs
cd tools/bulk_get_translator_ids
go run main.go
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
