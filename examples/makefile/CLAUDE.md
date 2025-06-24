# Makefile Mock Generation Example

This example demonstrates using a **simple Makefile** for reliable mock generation.

## Quick Commands

The Makefile provides simple commands for all mock generation tasks:

```bash
# Generate all mocks
make generate-mocks

# Clean and regenerate
make clean-mocks generate-mocks

# Run tests
make test

# Complete workflow
make all
```

## AI Prompting

You can ask AI assistants to run these commands:

**Generate mocks:**
```
"Run make generate-mocks to generate all the mocks for this project"
```

**Full workflow:**
```
"Run make all to clean, generate mocks, and run tests"
```

**Check results:**
```
"Run make list-mocks to see what mock files were generated"
```

## Expected Results

The Makefile will:
1. Install Mockery if not present
2. Generate mocks in `mocks/` directory using .mockery.yaml configuration
3. Create MockUserRepository.go, MockEmailService.go, MockCacheService.go
4. Run tests to verify everything works

## Key Advantages

- **Dead simple** - `make generate-mocks` is all you need
- **Self-documenting** - `make help` shows all available commands
- **Reliable** - Uses proven Mockery tool with consistent configuration
- **Fast** - No server startup or protocol overhead

## Validation

```bash
# See all available commands
make help

# Generate mocks and run tests
make all

# Check what was generated
make list-mocks
```

This approach gives you all the reliability of Mockery with maximum simplicity.