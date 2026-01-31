# CODEBUDDY.md

This file provides guidance to CodeBuddy Code when working with code in this repository.

## Project Overview

Ghost is a cross-platform desktop application for managing system hosts files. It's built with Go (backend) and Vue 3 (frontend) using the Wails framework, which allows building native desktop applications using web technologies.

### Core Features

- **Multi-Host Group Management**: Create, edit, delete, and manage multiple host file groups
- **Remote Host Support**: Configure URLs to periodically fetch remote host files
- **Smart Toggle Controls**: Selectively enable/disable host groups with visual indicators
- **Automatic Application**: When groups are toggled, saved, or refreshed, enabled groups are automatically applied to the system hosts file
- **Cross-Platform Compatibility**: Supports Windows, macOS, and Linux with platform-specific permission handling

## Development Commands

### Running the Application

```bash
# Development mode (with hot reload)
wails dev

# The frontend dev server runs on http://localhost:8173
# A dev server for accessing Go methods runs on http://localhost:34115
```

### Building

```bash
# Using the build script (recommended)
./scripts/build.sh

# This script:
# - Copies wails.exe.dev.manifest to wails.exe.manifest in build/windows/
# - Runs wails build
# - Cleans up the temporary manifest file regardless of build success/failure

# Build production executable directly
wails build

# The output executable is named "ghost.exe" (located in build/bin/)
```
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Run frontend dev server (for frontend-only development)
npm run dev

# Build frontend for production
npm run build

# Preview production build
npm run preview
```

## Architecture

### High-Level Structure

The application follows a clean separation of concerns with distinct layers:

```
├── main.go              # Wails app entry point
├── app.go               # App struct that bridges Wails and business logic
├── application/         # Business logic layer
│   └── host_app.go      # Core application logic for host management
├── models/              # Data models
│   └── host_models.go   # HostGroup, AppConfig, HostManager structures
├── storage/             # Data persistence layer
│   └── config_storage.go # JSON file storage for config and data
├── system/              # System hosts file operations
│   └── host_manager.go  # Reading/writing system hosts file
├── permissions/         # Cross-platform privilege management
│   └── elevation.go     # Admin/root privilege handling
├── remote/              # Remote host fetching
│   └── remote_fetcher.go # HTTP-based remote host content retrieval
└── frontend/            # Vue 3 frontend
    ├── src/
    │   ├── App.vue      # Main app component
    │   ├── components/  # Modular Vue components
    │   └── main.js      # Vue app initialization
```

### Key Architectural Patterns

**1. Wails Binding Layer (`app.go`)**

The `App` struct in `app.go` is bound to Wails and exposes methods to the frontend. This is the only interface between Go and JavaScript:

- All methods are public and exported
- Methods handle conversion between Go types and JSON for frontend consumption
- Methods delegate to `application.HostApp` for business logic

**2. Business Logic Layer (`application/host_app.go`)**

The `HostApp` struct contains core business logic:
- CRUD operations for host groups
- Remote host fetching and validation
- Applying enabled groups to system hosts
- Auto-refresh functionality for remote hosts
- Backup operations

Key methods:
- `GetHostGroups()`, `AddHostGroup()`, `UpdateHostGroup()`, `DeleteHostGroup()` - CRUD
- `ToggleHostGroup()` - Enable/disable groups (triggers auto-apply)
- `ApplyHosts()` - Applies all enabled groups to system hosts file
- `RefreshRemoteGroups()` - Fetches fresh content from all remote URLs
- `RefreshRemoteGroup(id)` - Fetches content for a specific remote group

**3. System Integration (`system/host_manager.go`)**

The `HostManager` handles all system hosts file operations:
- Detects platform-specific hosts file location
- Reads/writes system hosts with permission checking
- Applies Ghost-managed sections using markers (`# >>> Ghost Host Entries`)
- Creates and manages backups of system hosts file

Important: Ghost sections are delimited by start/end markers and preserve existing system hosts content outside these markers.

**4. Data Storage (`storage/config_storage.go`)**

Uses JSON files for persistence:
- `~/.ghost/config.json` - Application configuration
- `~/.ghost/data.json` - Host groups and state
- `~/.ghost/backups/` - Backup files with timestamps

The `ConfigStorage` struct handles all file I/O with mutex-based concurrency control.

**5. Frontend Architecture (`frontend/src/`)**

Component-based architecture with Vue 3 and Element Plus UI:

- **App.vue**: Main container managing overall state and coordinating components
- **Sidebar.vue**: Tree view displaying all host groups with search and enable/disable toggles
- **MainPanel.vue**: Right panel showing selected group details or system host preview
  - Contains LocalHostEditor.vue for local group editing
  - Contains RemoteHostEditor.vue for remote group editing
  - Contains SystemHostPreview.vue for system host file viewing
- **ActionBar.vue**: Top action bar with global actions (refresh remotes, backup)
- **AddGroupModal.vue**: Modal for creating new host groups

State Management:
- App.vue holds the main state and passes props down
- Components emit events to parent for state updates
- Dirty state tracking prevents accidental data loss when switching groups
- When a group is enabled/disabled/saved/refreshed, ApplyHosts is automatically called if the group is enabled

**6. Cross-Platform Permission Handling (`permissions/elevation.go`)**

Platform-specific privilege management:
- **Windows**: Uses application manifest (`resources.rc`) for permanent admin privileges via UAC
- **Linux**: Attempts graphical sudo tools (pkexec, gksudo, kdesudo) for elevation
- **macOS**: Uses AppleScript with administrator privileges

The `IsAdmin()` function checks current privilege level, and `RequestElevation()` prompts for elevated privileges when needed.

**7. Remote Host Fetching (`remote/remote_fetcher.go`)**

Handles HTTP-based remote host content retrieval:
- Validates fetched content as proper hosts file format
- 30-second timeout for HTTP requests
- Content validation checks for IP/domain pairs
- User-Agent header identifies the client

### Data Flow

**Applying Host Groups to System:**

1. Frontend calls `ApplyHosts()` via Wails binding
2. `HostApp.ApplyHosts()` checks write permissions
3. Loads all host groups from storage
4. Filters for enabled groups
5. Passes to `HostManager.ApplyHostGroups()`
6. `HostManager`:
   - Reads current system hosts
   - Removes existing Ghost-managed sections
   - Adds new Ghost section with enabled group content
   - Writes back to system hosts file

**Creating a New Host Group:**

1. Frontend shows AddGroupModal
2. User fills form (name, content, optional URL for remote)
3. Frontend calls `AddHostGroup()` via Wails binding
4. `HostApp.AddHostGroup()`:
   - Generates UUID for group
   - Validates required fields
   - Sets timestamps
   - Saves to data.json
5. Frontend reloads groups list

**Remote Host Update Flow:**

1. User clicks "获取host内容" button or "Refresh Remote Groups"
2. Frontend calls `RefreshRemoteGroup(id)` or `RefreshRemoteGroups()`
3. `RemoteFetcher` fetches content from URL
4. Content is validated as proper hosts format
5. Group content is updated in storage
6. If the group is enabled, `ApplyHosts()` is automatically called
7. System hosts file is updated with new content

### Important Behavioral Notes

**Auto-Apply Behavior:**
- When a group is toggled enabled/disabled: Auto-applies all enabled groups
- When a group is saved: Auto-applies if that group is enabled
- When remote content is refreshed: Auto-applies if the refreshed group is enabled
- When "Refresh Remote Groups" is clicked: Auto-applies if any groups are enabled

**Read-Only Mode:**
- When a host group (local or remote) is enabled, the editor switches to read-only preview mode
- This prevents accidental modifications to active content

**Remote Host Behavior:**
- Enabling/disabling a remote host group does NOT fetch fresh content
- Content updates only happen when user explicitly clicks "获取host内容" or "Refresh Remote Groups"
- Once content is updated and if group is enabled, changes auto-apply to system

**Permission Handling:**
- Windows: Application runs with permanent admin privileges once approved via UAC
- Linux/macOS: Authentication may be required periodically based on sudo timeout settings
- Permission elevation happens when writing to system hosts file

## Configuration Files

### Wails Configuration (`wails.json`)

```json
{
  "name": "ghost",
  "outputfilename": "ghost",
  "frontend:install": "npm install",
  "frontend:build": "npm run build",
  "frontend:dev:watcher": "npm run dev"
}
```

### Go Module (`go.mod`)

- Go version: 1.23
- Main dependencies: `github.com/wailsapp/wails/v2`, `github.com/google/uuid`

### Frontend Package (`frontend/package.json`)

- Vue 3 + Element Plus for UI
- Vite for build tooling and dev server
- Dev server runs on port 8173

## Testing

This codebase currently does not have automated tests. When adding tests:
- Go tests should be placed in `*_test.go` files alongside the code they test
- Use `go test ./...` to run all tests
- No test runner configuration files exist yet

## Data Storage Locations

- **Windows**: `%USERPROFILE%\.ghost\`
- **Unix/Linux/macOS**: `~/.ghost/`

Files:
- `config.json` - App configuration
- `data.json` - Host groups
- `backups/` - Backup files (both data and system hosts)

## Build Script

A build script is provided at `scripts/build.sh` for automated building with manifest handling:

```bash
./scripts/build.sh
```

**What the script does:**
1. Copies `build/windows/wails.exe.dev.manifest` to `build/windows/wails.exe.manifest`
2. Executes `wails build`
3. Removes the temporary `wails.exe.manifest` file (regardless of build success/failure)
4. Reports build status with clear step-by-step output

**Output:** The built executable is located at `build/bin/ghost.exe`

## System Hosts File Locations

- **Windows**: `C:\Windows\System32\drivers\etc\hosts`
- **Unix/Linux/macOS**: `/etc/hosts`

Ghost-managed entries are wrapped in:
```
# >>> Ghost Host Entries
[group content]
# <<< Ghost Host Entries
```

## Common Development Tasks

**Building the application:**
- Use the provided build script: `./scripts/build.sh`
- This handles manifest file copying and cleanup automatically
- Output: `build/bin/ghost.exe`

**Adding a new host group:**
1. Use AddGroupModal in UI or call `AddHostGroup()` directly
2. Groups get auto-generated UUIDs
3. Data saved to `~/.ghost/data.json`

**Modifying system hosts file operations:**
- Edit `system/host_manager.go`
- Key methods: `ReadSystemHosts()`, `WriteSystemHosts()`, `ApplyHostGroups()`

**Adding new Wails bindings:**
1. Add method to `App` struct in `app.go`
2. Method will be automatically exposed to frontend
3. Regenerate bindings with `wails generate module` if needed

**Frontend component development:**
- Components in `frontend/src/components/`
- Use Element Plus for UI components
- Wails bindings available via `../wailsjs/go/main/App`

**Cross-platform considerations:**
- Use `runtime.GOOS` for platform-specific code
- Permission handling varies by platform (see `permissions/elevation.go`)
- File paths use `filepath.Join()` for cross-platform compatibility
