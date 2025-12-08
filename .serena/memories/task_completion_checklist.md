# Task Completion Checklist

When completing a coding task in this project, follow these steps:

## 1. Code Generation
If templates were modified:
```bash
templ generate
```

## 2. Mock Generation
If interfaces were modified:
```bash
make gen-mocks
```

## 3. Linting and Formatting
Always run before committing:
```bash
make lint
```
This will:
- Tidy go.mod dependencies
- Run go vet
- Format code with golangci-lint
- Run static analysis linters

## 4. Testing
Run the test suite:
```bash
make test
```
Tests run with:
- 10 second timeout
- Race detection enabled
- Benchmark memory tracking

## 5. Build Verification
Ensure the project builds:
```bash
make build
```

## 6. Database Migrations
If database schema changed:
```bash
make migrate
```

## Notes
- All commands should pass without errors
- The project uses tabs for indentation (size 4)
- Mock files are auto-excluded from linting
- Tests should complete within 10 seconds