---
applyTo: '**/*.go, **/*.sh, **/*.md'
---
Provide project context and coding guidelines that AI should follow when generating code, answering questions, or reviewing changes.

# Duh Project Instructions

## Project Overview

**Duh** is a simple and effective dotfiles manager written in Go. It helps users manage shell aliases, environment exports, and shell functions through a CLI interface.

## Techs used
Go programming language (v1.25.4), Cobra for CLI

## Architecture

The project follows Clean Architecture principles with three main layers:


### Domain Layer (`internal/domain/`)
- **Entities** (`entity/`): Core business objects (Repository, Script, Function, etc.)
- **Repositories** (`repository/`): Interface definitions for data access
- **Services** (`service/`): Business logic and use cases
- **Utils** (`utils/`): Domain-specific utilities

### Application Layer (`internal/application/`)
- **CLI** (`cli/`): Cobra-based command line interface
- **Contexts** (`contexts/`): Application context initialization

### Infrastructure Layer (`internal/infrastructure/`)
- **Filesystem** (`filesystem/`): File system operations and repositories
- **Git** (`gitt/`): Git operations
- **TOML** (`tomll/`): Configuration file handling

## Key Components

### CLI Commands Structure
All CLI commands follow this pattern:
- Located in `internal/application/cli/`
- Named as `*_subcommand.go`
- Use Cobra framework
- Follow the builder pattern: `Build*Subcommand(cliService service.CliService) *cobra.Command`

### Testing Strategy

#### Unit Tests (`*_test.go` files)
- Located alongside source files
- Use `testify/assert` for assertions
    - always prefer `assert.NoError(t, err)` over `if err != nil { t.Fatalf(...) }`
- Mock repositories using interfaces
- Test individual components in isolation

#### Integration Tests (`tests/integration/`)
- Test interactions between components
- Use real file system operations in temp directories

#### E2E Tests (`tests/e2e/`)
- Test complete user workflows
- Use actual CLI execution
- Override XDG paths for isolation

## Code Standards

### Naming Conventions
- Use Go standard naming (camelCase for private, PascalCase for public)
- CLI commands use kebab-case with aliases
- File names use snake_case
- Test files end with `_test.go`

### Error Handling
- Return errors from functions, don't panic
- Use `cmd.PrintErrf()` for CLI error output
- Wrap errors with context when appropriate

### Testing Requirements
When adding new functionality:
1. **Always add unit tests** in `*_test.go` files
2. **Add E2E tests** for new CLI commands
3. **Update mock repositories** when adding new data access patterns
4. **Test both success and error paths**
5. **Test command aliases** for CLI commands

### CLI Command Patterns
```go
func Build*Subcommand(cliService service.CliService) *cobra.Command {
    cmd := &cobra.Command{
        Use:     "command [subcommand]",
        Aliases: []string{"alias1", "alias2"},
        Short:   "Short description",
    }
    
    subCmd := &cobra.Command{
        Use:   "subcommand",
        Short: "Subcommand description",
        Run: func(cmd *cobra.Command, args []string) {
            // Implementation
        },
    }
    
    cmd.AddCommand(subCmd)
    return cmd
}
```

### Flag Handling
- Use `cmd.Flags().Bool*()` for boolean flags
- Support both long (`--flag`) and short (`-f`) versions
- Add helpful descriptions for all flags
- Handle flag parsing errors appropriately

## Dependencies

### Core Dependencies
- **Cobra**: CLI framework
- **TOML**: Configuration file parsing
- **XDG**: Cross-platform directory paths
- **Testify**: Testing assertions

### Testing Dependencies
- **Testify**: Assertions and test utilities
- **Go standard library**: `testing` package

## File Structure Guidelines

```
internal/
├── application/
│   └── cli/                 # CLI command implementations
├── domain/
│   ├── entity/             # Core business objects
│   ├── repository/         # Repository interfaces
│   ├── service/           # Business logic
│   └── utils/             # Domain utilities
└── infrastructure/
    ├── filesystem/        # File system operations
    ├── gitt/             # Git operations
    └── tomll/            # TOML handling

tests/
├── e2e/                   # End-to-end tests
└── integration/           # Integration tests

cmd/
└── cli/                   # Main CLI entry point
```

## Common Development Tasks

### Adding a New Flag to Existing Command
1. Add flag definition using `cmd.Flags().Type()`
2. Read flag value using `cmd.Flags().GetType()`
3. Add conditional logic based on flag value
4. Update tests to cover new flag behavior
5. Add E2E tests for new flag

### Adding Mock Data for Tests
1. Update repository mocks with realistic data
2. Include edge cases (empty data, warnings, errors)
3. Ensure data is consistent across related tests
4. Update both unit and E2E test data

### Working with Shell Functions
- Scripts contain multiple functions
- Functions have names and documentation
- Scripts can have warnings
- Distinguish between "activated" (enabled) and "all" scripts
- Support both `GetActivatedFunctions()` and `GetAllFunctions()`

## Important Notes

- **Always test CLI commands** with both success and error scenarios
- **Mock external dependencies** (file system, git) in unit tests
- **Use temp directories** for tests that create files
- **Test command aliases** - many commands have multiple valid names
- **Handle special characters** in alias/export values properly
- **Maintain backwards compatibility** when adding new features
- **Follow the existing patterns** for consistency
