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
4. **Auto-Sync Enabled Groups** - When remote host content is updated (via "获取host内容" button or "Refresh Remote Groups"), automatically applies changes to system hosts if the group is enabled
5. **One-Click Application** - Merge and apply all enabled host groups to system hosts file in one click
6. **Cross-Platform Compatibility** - Automatically detects hosts file locations across Windows, macOS, and Linux systems with appropriate permission handling

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

### Permission Management

Ghost Host Manager handles permissions differently across platforms:

- **Windows**: Uses application manifest to request permanent administrator privileges. Once approved via UAC, the application runs with elevated permissions throughout the session.
- **Linux**: Utilizes graphical sudo tools (pkexec, gksudo, kdesudo) to request permissions when needed. Authentication may be required periodically based on system sudo timeout settings.
- **macOS**: Uses AppleScript to request administrator privileges when needed. Authentication may be required periodically based on system sudo timeout settings.

For detailed information about permission handling on each platform, please see our [Permissions Guide](./docs/permissions.md).

## Remote Host Behavior

- **Enable/Disable**: Toggling a remote host group only enables/disables it without fetching fresh content. The system hosts file is updated to reflect the enabled/disabled state, but the remote content itself is not refreshed.
- **Read-only Mode**: When a host group (local or remote) is enabled, the editor switches to read-only preview mode to prevent accidental modifications to active content.
- **Content Updates**: Click "获取host内容" or "Refresh Remote Groups" to manually update remote host content.
- **Auto-Apply**: When a remote host group is enabled and its content is updated, changes are automatically applied to the system hosts file if the group is enabled

## Development Approach

Built using AI-assisted development methodologies with:
- Component-based architecture for improved maintainability
- Intelligent state management patterns
- Automated code quality checks
- AI-enhanced refactoring and optimization
- Smart UI/UX design patterns

## User Interface

The application provides a clean, intuitive interface with:
- **Action Bar** - Contains "Refresh Remote Groups" and "Backup Now" buttons
- **Automatic Apply Behavior** - Hosts are automatically applied to the system when groups are toggled, saved, or refreshed (if the group is enabled)
- **System Host Preview** - Shows current system hosts file content with refresh capability

## Documentation

For detailed information about using Ghost Host Manager, please see our documentation:

- [Quick Start Guide](./docs/quick-start.md) - Installation and basic usage
- [Permissions Guide](./docs/permissions.md) - Detailed explanation of permission handling on each platform
- [Full Documentation Index](./docs/index.md) - Complete list of available documentation

## Development Setup

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Build

To build a redistributable, production mode package, use `wails build`.
