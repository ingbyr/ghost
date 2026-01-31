# Ghost Host Manager - Agent Development Guide

This guide helps agentic coding agents understand the Ghost Host Manager codebase structure, build processes, and coding conventions.

## Project Overview

Ghost is a cross-platform desktop application built with Wails v2 that manages system host files. It allows users to create, organize, and apply local and remote host configurations.

**Architecture:**
- **Backend:** Go (v1.23) with Wails v2 framework
- **Frontend:** Vue 3 + Element Plus UI library + Vite
- **Build:** Wails handles cross-compilation and packaging

## Build & Development Commands

### Backend (Go)
```bash
# Build the application
wails build

# Run in development mode
wails dev

# Build for production with specific platform
wails build -platform windows/amd64
wails build -platform darwin/amd64
wails build -platform linux/amd64

# Generate Wails bindings (if needed)
wails generate module
```

### Frontend (Vue.js)
```bash
cd frontend

# Install dependencies
npm install

# Development server (runs on port 8173)
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

### Full Application
```bash
# Install frontend dependencies and build
wails build

# Development mode with hot reload
wails dev
```

## Testing

**Note:** This codebase currently has no automated tests. When adding tests:

### Go Tests
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests in specific package
go test ./models
go test ./application

# Run single test
go test -run TestSpecificFunction ./package
```

### Frontend Tests (if added)
```bash
cd frontend
npm test              # Run all tests
npm run test:unit     # Run unit tests only
npm run test:coverage # Run with coverage
```

## Code Style & Conventions

### Go Backend

#### Naming Conventions
- **PascalCase** for exported types, functions, constants: `HostGroup`, `NewHostApp()`
- **camelCase** for unexported items: `configStorage`, `hostManager`
- **ALL_CAPS** for constants: `GhostSectionStart`, `ConfigFile`
- **Package names:** lowercase, single word: `models`, `application`, `system`

#### Import Organization
```go
import (
    // Standard library
    "context"
    "fmt"
    "time"
    
    // Third-party libraries
    "github.com/google/uuid"
    "github.com/wailsapp/wails/v2"
    
    // Local packages
    "ghost/models"
    "ghost/storage"
)
```

#### Error Handling
- Always wrap errors with context: `fmt.Errorf("failed to load config: %w", err)`
- Use descriptive error messages with operation context
- Log errors appropriately for debugging
- Return errors from functions, don't panic

#### Function Documentation
```go
// NewHostApp creates a new Host application instance
func NewHostApp() (*HostApp, error) {
    // Implementation
}

// ApplyHosts applies all enabled host groups to the system
// It preserves existing system host entries and only modifies the Ghost section
func (app *HostApp) ApplyHosts() error {
    // Implementation
}
```

#### Struct Tags
- Use JSON tags for API/export: `json:"id"`
- Use omitempty for optional fields: `json:"description,omitempty"`

### Vue.js Frontend

#### Component Structure
```vue
<template>
  <!-- Template content with Element Plus components -->
</template>

<script>
import { ComponentA, ComponentB } from './components'

export default {
  name: 'ComponentName',
  components: {
    ComponentA,
    ComponentB
  },
  props: {
    propName: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      localState: 'value'
    }
  },
  methods: {
    methodName() {
      // Method implementation
    }
  }
}
</script>

<style scoped>
/* Component-specific styles */
</style>
```

#### Naming Conventions
- **PascalCase** for component names: `Sidebar.vue`, `MainPanel.vue`
- **camelCase** for props, data, methods: `selectedGroup`, `loadHostGroups()`
- **kebab-case** for CSS classes: `tree-item`, `sidebar-header`
- **UPPER_SNAKE_CASE** for CSS constants: `$PRIMARY_COLOR`

#### Event Handling
- Use `$emit` for parent communication: `@click="$emit('select-group', group)"`
- Use async/await for API calls
- Handle errors with try/catch and show user-friendly messages

### File Organization

#### Go Backend
```
/                    # Root
├── main.go          # Application entry point
├── app.go           # Wails app structure and API methods
├── models/          # Data models and structs
├── application/     # Business logic layer
├── storage/         # Data persistence
├── system/          # System integration (hosts file)
├── remote/          # Remote content fetching
└── permissions/     # Privilege management
```

#### Frontend
```
frontend/src/
├── main.js          # Vue app initialization
├── App.vue          # Root component
├── components/      # Reusable components
├── style.css        # Global styles
└── wailsjs/         # Generated Wails bindings
```

## Development Patterns

### Backend Patterns
- **Repository Pattern:** `ConfigStorage` handles data persistence
- **Service Layer:** `HostApp` contains business logic
- **Dependency Injection:** Pass dependencies through constructors
- **Context Usage:** Use `context.Context` for operations

### Frontend Patterns
- **Composition over Inheritance:** Use component composition
- **Prop Validation:** Always define prop types and requirements
- **Event-driven:** Use events for component communication
- **State Management:** Local component state for simple cases

## Important Notes

### Security
- Host file modification requires admin privileges
- Validate remote URLs before fetching content
- Sanitize user inputs when applying to system hosts

### Cross-Platform
- Use `runtime.GOOS` for platform-specific code
- Handle different file paths for Windows/Linux/macOS
- Test on all target platforms before release

### Wails Integration
- Backend methods must be exported (PascalCase) to be callable from frontend
- Use generated bindings in `wailsjs/go/` for frontend-backend communication
- Handle context lifecycle properly in `startup()` and `shutdown()`

### Performance
- Use file locking for concurrent access to hosts file
- Implement proper error handling for remote content fetching
- Consider caching for frequently accessed data