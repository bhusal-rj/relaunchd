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
