# Ghost - AI-Powered Host Manager

An advanced cross-platform GUI application built with Go and Wails framework for managing and switching host files. Features AI-driven development practices and supports multi-host group management with remote host capabilities.

## New Feature: Tree Structure UI

Latest version introduces a brand new left-right panel interface design:

- **Left Tree Structure** - Clearly displays all host groups
- **Right Edit Panel** - Provides detailed content editing capabilities
- **Search Functionality** - Quickly find specific host groups
- **Status Management** - Real-time display of group enable/disable status
- **Data Loss Protection** - Detects unsaved changes to prevent accidental loss

## Core Features

1. **Multi-Host Group Management** - Create, edit, delete and manage multiple host file groups with intuitive UI
2. **Remote Host Support** - Configure URLs to periodically fetch remote host files and apply them locally
3. **Smart Toggle Controls** - Selectively enable/disable host groups with visual indicators
4. **One-Click Application** - Merge and apply all enabled host groups to system hosts file in one click
5. **Cross-Platform Compatibility** - Automatically detects hosts file locations across Windows, macOS, and Linux systems with appropriate permission handling

## Technical Architecture

- **Frontend** - Vue 3 + Element Plus UI library with modern component architecture
- **Backend** - Go language implementing business logic and system interactions
- **Desktop Framework** - Wails framework bridging Go backend with web technologies
- **Data Persistence** - JSON-based file storage solution with backup capabilities
- **Remote Fetching** - Secure URL-based host content retrieval with validation
- **AI-Enhanced Development** - AI-assisted code generation, refactoring, and optimization

## AI-Driven Development Highlights

- **Component Architecture** - Clean separation of concerns with dedicated components for local/remote hosts, system preview, and action controls
- **Intelligent State Management** - Smart dirty state detection with floating save buttons that appear only when content is modified
- **Responsive UI Design** - Adaptive layout with tree view navigation and dynamic content panels
- **Smart Validation** - Real-time form validation and content verification
- **Automated Refactoring** - AI-assisted migration from monolithic component to modular architecture

## Cross-Platform Compatibility

- **Windows**: `C:\Windows\System32\drivers\etc\hosts`
- **Unix/Linux/macOS**: `/etc/hosts`
- **Permissions**: Requires administrator privileges for system hosts file modification

## Development Approach

Built using AI-assisted development methodologies with:
- Component-based architecture for improved maintainability
- Intelligent state management patterns
- Automated code quality checks
- AI-enhanced refactoring and optimization
- Smart UI/UX design patterns

## Development Setup

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Build

To build a redistributable, production mode package, use `wails build`.