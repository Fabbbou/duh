# Contributing to Duh

Thank you for your interest in contributing to Duh! 🎉

## 🚀 Quick Start

### Prerequisites
- Go 1.25.4+ 
- Git
- a Go editor (I like VSCode + Go extensions, but just keep your fav editor)

### Development Setup

1. **Fork and clone the repository (use the github button)**

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Run the tests:**
    **Unit tests:**
    They should be all green, but better check the CI to confirm what you should expect
    Run the go test command that run all tests
    ```bash
    go test ./...
    ```
    
    You may use your own editor setup to run the tests (I do), but keep in mind to double check with the command as well, cause it's how the CI do.
    
    **Manual Tests:**
    You could also use Manual
    I often use directly `go run cmd/cli/main.go <cli args>` command to run the cli when i wanna manually test
    
    ```bash
    go build -o duh cmd/cli/main.go
    ./duh --help
    ```
## 📁 Project Structure
I followed DDD conventions (example [here](https://github.com/sklinkert/go-ddd))
```
duh/
├── cmd/cli/               # Main CLI entry point
├── internal/
│   ├── application/cli/   # CLI commands and handlers
│   ├── domain/           # Business logic
│   │   ├── entity/       # Data models
│   │   ├── repository/   # Data access interfaces
│   │   └── service/      # Business services
│   └── infrastructure/   # External dependencies (file system, etc.)
├── tests/
│   ├── e2e/             # End-to-end tests
│   └── integration/     # Integration tests
└── .github/workflows/    # CI/CD pipelines
```

## 🛠️ Development Workflow

### Adding New Features

1. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes:**
   - Add new CLI commands in `internal/application/cli/`
   - Add business logic in `internal/domain/service/`
   - Add tests for new functionality

3. **Create a Pull Request from your fork repo**

### Adding New CLI Commands

1. **Create command file:** `internal/application/cli/your_command.go`
2. **Follow existing patterns:** See `alias_cli.go` or `exports_cli.go`
3. **Add to root CLI:** Update `internal/application/cli/cli.go`
4. **Add tests:** Create corresponding test file
5. **Add E2E tests:** Update `tests/e2e/nominal_test.go`

### Code Style

- **Go formatting:** Use `go fmt` (automatic in most editors)
- **Linting:** Run `go vet` before committing
- **Testing:** Write tests for new functionality
- **Error handling:** Use explicit error checking, avoid panics
- **CLI patterns:** Follow existing cobra command patterns

## 🧪 Testing

### Unit Tests
```bash
go test ./internal/...
```

### Integration Tests
```bash
go test ./tests/integration/...
```

### End-to-End Tests
```bash
go test ./tests/e2e/ -v
```

### Test Coverage
```bash
go test -cover ./...
```

## 📝 Pull Request Guidelines

### Before Submitting
- [ ] Tests pass (`go test ./...`)
- [ ] Code is formatted (`go fmt`)
- [ ] No linting errors (`go vet`)
- [ ] New features have tests
- [ ] Documentation is updated if needed

### PR Description
- **Clear title:** Describe what the PR does
- **Problem:** What issue does this solve?
- **Solution:** How does this PR solve it?
- **Testing:** How was this tested?

### Example PR Title
```
feat: add self-update command for automatic binary updates
fix: resolve Windows installation path issues
docs: update README with new installation methods
```

## 🐛 Reporting Issues

### Bug Reports
Use the **Bug Report** template and include:
- Duh version (`duh --version`)
- Operating system
- Steps to reproduce
- Expected vs actual behavior
- Error messages or logs

### Feature Requests
Use the **Feature Request** template and include:
- Clear description of the feature
- Use case or problem it solves
- Proposed CLI interface (if applicable)

## 💡 Development Tips

### Local Testing
```bash
# Build and test locally
go build -o duh cmd/cli/main.go
./duh alias set test "echo hello"
./duh inject

# Test installation script
cat install.sh | sh
```

### Debugging
```bash
# Run with verbose output
./duh alias set test "echo hello" -v

# Check file locations
ls ~/.local/share/duh/
```

### Adding Dependencies
```bash
go get github.com/new/dependency
go mod tidy
```

## 🎯 Good First Issues

Looking for your first contribution? Check issues labeled:
- `good first issue`
- `help wanted`
- `documentation`

Common starter tasks:
- Improve error messages
- Add command aliases
- Enhance shell completion
- Write additional tests
- Improve documentation

## 📞 Getting Help

- **Issues:** For bugs and feature requests
- **Discussions:** For questions and ideas
- **Email:** For security issues

## 📄 License

By contributing to Duh, you agree that your contributions will be licensed under the same license as the project.

---

Thank you for contributing! 🚀