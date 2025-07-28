# Contributing to AnythingTranslate Go Client

Thank you for considering contributing to the AnythingTranslate Go Client! This document provides guidelines and instructions for contributing.

## Code of Conduct

Please be respectful and considerate of others when contributing to this project. We aim to foster an inclusive and welcoming community.

## How to Contribute

### Reporting Bugs

If you find a bug, please create an issue with the following information:

1. A clear, descriptive title
2. A detailed description of the issue
3. Steps to reproduce the bug
4. Expected behavior
5. Actual behavior
6. Any relevant logs or error messages
7. Your Go version and operating system

### Suggesting Enhancements

If you have an idea for an enhancement, please create an issue with:

1. A clear, descriptive title
2. A detailed description of the enhancement
3. Any relevant examples or use cases

### Pull Requests

1. Fork the repository
2. Create a new branch for your changes
3. Make your changes
4. Run tests to ensure they pass
5. Submit a pull request

## Development Setup

1. Clone the repository
2. Install Go (version 1.18 or later recommended)
3. Run tests with `go test ./...`

## Coding Standards

- Follow standard Go coding conventions
- Use `gofmt` to format your code
- Write godoc-compatible comments for all exported types and functions
- Add tests for new functionality
- Ensure all tests pass before submitting a pull request

## Adding New Translators

If you want to add support for a new translator:

1. Use the tools in the `tools` directory to get the translator ID
2. Add a new constant to `translators.go`
3. Add documentation for the new translator
4. Add tests for the new translator

## License

By contributing to this project, you agree that your contributions will be licensed under the project's [MIT License](LICENSE).