# Code Style and Conventions

## General Style
- **Indentation**: Tabs with size 4 (as per .editorconfig)
- **Line Endings**: LF (Unix-style)
- **Charset**: UTF-8
- **Trailing Whitespace**: Must be trimmed
- **Final Newline**: Required in all files

## YAML/TOML Files
- **Indentation**: 2 spaces (not tabs)

## Go-Specific Conventions
- **Go Version**: 1.25
- **Formatting**: Uses `gofmt`, `gofumpt`, and `goimports`
- **Import Grouping**: Via `gci` with sections:
  - standard library
  - default (third-party)
  - localmodule (project imports)

## Linting Rules
Enabled linters:
- `forbidigo`: Prevents use of forbidden identifiers
- `lll`: Line length limit
- `prealloc`: Detects slice declarations that could preallocate
- `predeclared`: Finds shadowing of predeclared identifiers
- `staticcheck`: Static analysis checks

Disabled linters:
- `errcheck`: Error checking is disabled

Mock files in `mocks/` directory are excluded from most linting rules.

## File Generation
- Templates use `templ` and must be regenerated via `templ generate` before building
- Mocks are generated using `moq` tool