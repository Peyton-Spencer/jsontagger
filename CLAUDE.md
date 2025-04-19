# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands
- Build: `go build -o jsontagger`
- Run: `./jsontagger -file path/to/your/file.go [-snake|-camel]`
- Test all: `go test ./...`
- Test single: `go test -run TestProcessFile`
- Format code: `go fmt ./...`

## Code Style Guidelines
- Follow Go standard formatting (use `gofmt`)
- Use snake_case for file names
- Use camelCase for variable names
- Use PascalCase for exported functions, types, and variables
- Maintain consistent spacing (no trailing whitespace)
- Use meaningful variable and function names
- Handle errors explicitly with proper return values
- When processing files, validate inputs before modification
- Prefer early returns for error handling
- Write descriptive comments for public functions and types
- Use proper indentation (tabs, not spaces)
- Group related imports and separate standard library imports