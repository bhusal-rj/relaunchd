# Process Manager

## Overview
The Process Manager is a core component of relaunchd that handles application lifecycle management. It provides robust control over process execution, monitoring, and automatic restarts.

## Features
- Process Lifecycle Control: Start, stop, and restart processes programmatically
- Background Process Support: Run applications detached from the terminal
- Automatic Restart: Configurable restart policies when processes exit
- State Tracking: Monitor process status, uptime, PID, and exit conditions
- Graceful Termination: Properly shut down processes to prevent resource leaks
- Environment Configuration: Set working directory and environment variables
- I/O Handling: Manage standard input/output streams

## How It Works
The Process Manager acts as a wrapper around your application, providing:

- Process Creation: Builds the command from your configuration, including arguments and environment settings
- Execution Control: Manages the running state of processes
- Monitoring: For background processes, watches for exits and handles restarts
- Resource Management: Ensures proper cleanup when processes terminate

## Integration with File Watcher
The Process Manager integrates with the File Watcher to enable hot-reload development:

1. File Watcher detects changes in specified directories/files
2. When changes are detected, the Process Manager:
   - Gracefully terminates the running process
   - Applies a configurable restart delay if specified
   - Launches a fresh instance of the application
   - Tracks restart count against maximum limits

## Configuration Options
The Process Manager uses these configuration settings from your YAML file:

## Usage Examples
The Process Manager is used internally by relaunchd's commands:

- `relaunchd start`: Creates and starts a new managed process
- `relaunchd stop`: Gracefully terminates a running process
- `relaunchd restart`: Stops and restarts a process
- `relaunchd status`: Reports current process state and statistics

## Implementation Details
The Process Manager maintains a state machine for each process, tracking:

- Current execution state (running, stopped, restarting, failed)
- Process ID
- Start and stop times
- Restart count
- Last encountered error

This information enables accurate reporting and intelligent restart handling based on your configuration.