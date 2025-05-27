## 🚀 relaunchd

**relaunchd** is a lightweight, developer-friendly process manager and file watcher written in Go.  
It monitors files and directories for changes and automatically restarts your application using a simple YAML configuration.  
Inspired by tools like PM2 and nodemon, it's ideal for hot-reload development workflows in any language or framework.

---

## ✨ Features

- 🔧 **YAML-based configuration**  
  Define what to watch and which command to run — no scripting needed.

- 👀 **File & directory watching**  
  Supports glob patterns to recursively monitor source files.

- 🔄 **Automatic restarts**  
  On file change, your specified command is stopped and relaunched seamlessly.

- 🧠 **Background process support (PM2-style)**  
  Run and manage long-lived processes in the background with PID tracking.

- 📊 **CLI Interface**  
  Commands like `relaunchd start`, `stop`, `status` for easy control.

- 🖥️ **Cross-platform compatibility**  
  Works on Linux, macOS, and Windows.

- 🧪 **Minimal dependencies**  
  Written in Go, portable and fast with zero runtime bloat.

## 🛠️ Installation & Usage

### Prerequisites

- Go 1.16 or higher
- Git (for cloning the repository)

### Installing relaunchd

#### From Source

```bash
# Clone the repository
git clone https://github.com/bhusal-rj/relaunchd.git
cd relaunchd

# Build the binary
go build -o relaunchd cmd/relaunchd/main.go

# Optional: Move to PATH
sudo mv relaunchd /usr/local/bin/ # Linux/macOS
# or add to your PATH on Windows
```

#### Using Go Install

```bash
go install github.com/bhusal-rj/relaunchd/cmd/relaunchd@latest
```

### Basic Usage

1. Create a configuration file named `relaunch.yml` in your project:

```yaml
watch:
  - "**/*.go"  # Watch all Go files recursively
  - "!vendor/" # Exclude vendor directory
command: "go run main.go"
debounce: 500  # Wait 500ms after changes before restarting
```

2. Start relaunchd:

```bash
relaunchd start
```

### Environment-Specific Considerations

#### Linux

```bash
# Run in background and redirect output
relaunchd start &> relaunchd.log &

# Check process status
relaunchd status

# Stop the process
relaunchd stop
```

#### macOS

Works identical to Linux. For automatic startup, consider using `launchd`:

```bash
# Create a launch agent (example)
cat > ~/Library/LaunchAgents/com.user.relaunchd.plist << EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.user.relaunchd</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/relaunchd</string>
        <string>start</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>WorkingDirectory</key>
    <string>/path/to/your/project</string>
</dict>
</plist>
EOF

# Load the agent
launchctl load ~/Library/LaunchAgents/com.user.relaunchd.plist
```

#### Windows

```powershell
# Run in background
Start-Process relaunchd -ArgumentList "start" -NoNewWindow

# Check process status
relaunchd status

# Stop the process
relaunchd stop
```

For Windows services, consider using [NSSM](https://nssm.cc/) to register relaunchd as a system service.

### Development Mode

When developing relaunchd itself:

```bash
# Run with hot reload (requires air: https://github.com/cosmtrek/air)
air

# Run tests
go test ./...

# Run with debug output
go run cmd/relaunchd/main.go --debug start
```

## 🗺️ Development Roadmap

### Phase 1: Core Foundation
- [x] Set up Go project structure with modules
- [x] Implement basic YAML configuration parser
- [x] Create file watching system using fsnotify
- [x] Build simple process management (start/stop)
- [x] Implement basic CLI command structure


### Phase 2: Process Management
- [ ] Develop background process handling with PID tracking
- [ ] Implement graceful shutdown mechanisms
- [ ] Add signal handling (SIGTERM, SIGINT, etc.)
- [ ] Create process status reporting functionality
- [ ] Build the "status" command implementation
- [ ] Add process restart capabilities

### Phase 3: Advanced File Watching
- [ ] Implement glob pattern support
- [ ] Add directory recursion capabilities
- [ ] Create file change debouncing mechanism
- [ ] Develop file type filtering
- [ ] Implement watch exclusion patterns
- [ ] Add support for multiple watch configurations

### Phase 4: Full CLI Experience
- [ ] Complete all CLI commands (start, stop, status, list)
- [ ] Add command flags and options
- [ ] Implement configuration validation
- [ ] Create helpful error messages
- [ ] Add colorized console output
- [ ] Implement verbosity levels for output

### Phase 5: Cross-Platform Compatibility
- [ ] Test and fix Windows-specific issues
- [ ] Ensure macOS compatibility
- [ ] Handle path differences between operating systems
- [ ] Verify process management across platforms
- [ ] Create platform-specific installation guides

### Phase 6: Polish and Release
- [ ] Create comprehensive documentation
- [ ] Build example configurations
- [ ] Implement version command
- [ ] Package for distribution
- [ ] Set up CI/CD pipeline
- [ ] Create user guides and tutorials
