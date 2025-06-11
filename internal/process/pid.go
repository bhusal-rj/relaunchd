package process

import "os"

type PIDInfo struct {
	PID       int    `json:"pid"`        // The actual process ID provided by the OS
	Name      string `json:"name"`       // Human-readable name for the process
	Command   string `json:"command"`    // The command that was executed to start the process
	StartTime string `json:"start_time"` // The time when the process was started
	Status    string `json:"status"`     // Current status of the process
}

// This is the manger that handles our process.
type PIDManager struct {
	pidFile string              // Path where we save the PID data
	pids    map[string]*PIDInfo // In memort cache of all the tracked process
}

func NewPIDManager(pidDir string) *PIDManager {
	if pidDir == "" {
		pidDir = "/tmp/relaunchd" // Default location for PID files
	}

	//Create the PID directory if it doesn't exist
	if err := os.MkdirAll(pidDir, 0755); err != nil {
		panic("Failed to create PID directory: " + err.Error())
	}
	return &PIDManager{
		pidFile: pidDir + "/pids.json",
		pids:    make(map[string]*PIDInfo),
	}
}

// Add a new process to tracking
func (pm *PIDManager) AddProcess(name string, pid int, command string) error {
	// Create PIDInfo, add to map, save to file
	return nil
}

// Get information about a specific process
func (pm *PIDManager) GetProcess(name string) (*PIDInfo, bool) {
	// Look up in map, check if still running, return info
	return nil, false
}

// Check if a process is still running on the system
func (pm *PIDManager) IsProcessRunning(pid int) bool {
	// Use OS calls to verify process exists
	return false

}

// Save current PID map to JSON file
func (pm *PIDManager) SavePIDs() error {
	// Convert map to JSON, write to file
	return nil
}
