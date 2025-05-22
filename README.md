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



## 🗺️ Development Roadmap

### Phase 1: Core Foundation
- [x] Set up Go project structure with modules
- [x] Implement basic YAML configuration parser
- [x] Create file watching system using fsnotify
- [ ] Build simple process management (start/stop)
- [ ] Implement basic CLI command structure
- [ ] Set up logging framework

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
