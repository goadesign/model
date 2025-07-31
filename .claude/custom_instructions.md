# Custom Instructions for Model Project

## Build and Testing Guidelines

- **ONLY run `make` after making code changes** - this builds the project and runs tests
- **DO NOT run npm commands** - the frontend is built as part of the make process
- **DO NOT attempt to start development servers** - the user will do it
- **DO NOT run webpack, yarn, or other frontend build tools directly** - everything is handled by make

## Testing Workflow

1. Make code changes
2. Run `make` to build and test
3. Stop here

This project uses a Go-based build system that handles all compilation and bundling automatically.