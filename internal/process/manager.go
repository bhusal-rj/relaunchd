package manager

import (
	"bhusal-rj/relaunchd/internal/config"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

type ProcessState string

const (
	StateNotStarted ProcessState = "not_started"
	StateRunning    ProcessState = "running"
	StateStopped    ProcessState = "stopped"
	StateFailed     ProcessState = "failed"
)

type ProcessManager struct {
	config *config.Config //actual configuration of the config for the process manager
	cmd    *exec.Cmd      //command to be excuted
	state  ProcessState   //current state of the process
	pid    int            //process id of the running process
	mutex  sync.Mutex     //mutex to protect the state and pid
}

// NewProcessManager creates a new ProcessManager instance
func NewProcessManager(cfg *config.Config) *ProcessManager {
	return &ProcessManager{
		config: cfg,
		state:  StateNotStarted,
		cmd:    nil,
		pid:    0,
		mutex:  sync.Mutex{},
	}
}

// Parse Command splits a command string into executable and arguments
func ParseCommand(commandStr string) (string, []string, error) {
	if commandStr == "" {
		return "", nil, fmt.Errorf("command string is empty")
	}

	parts := strings.Fields(commandStr)

	if len(parts) == 0 {
		return "", nil, fmt.Errorf("command string is empty")
	}

	executable := parts[0]

	args := parts[1:]
	return executable, args, nil
}

func (pm *ProcessManager) Start() error {

	//  Lock the mutex for the thread safety
	pm.mutex.Lock()

	// Unlock the mutex when the function returns
	defer pm.mutex.Unlock()

	// If the process is already running, return an error
	if pm.state == StateRunning {
		return fmt.Errorf("process is already running")
	}

	// Parse the command into executable and arguments
	executable, args, err := ParseCommand(pm.config.Command.Exec)

	if err != nil {
		return fmt.Errorf("failed to parse command: %v", err)
	}

	// Create the command
	pm.cmd = exec.Command(executable, args...)

	// Setup the standard output in console
	pm.cmd.Stdout = os.Stdout
	pm.cmd.Stderr = os.Stderr

	// Start the specified command and incase of error, set the state to failed
	if err := pm.cmd.Start(); err != nil {
		pm.state = StateFailed
		return fmt.Errorf("failed to start process: %v", err)
	}

	// Set the process ID and state
	pm.pid = pm.cmd.Process.Pid
	pm.state = StateRunning
	fmt.Printf("Process started with PID: %d\n", pm.pid)

	//Start the goroutine to wait for the process to finish
	go func() {
		_ = pm.cmd.Wait()

		pm.mutex.Lock()

		if pm.state == StateRunning {
			pm.state = StateStopped
			log.Printf("Process with PID %d has stopped\n", pm.pid)
			pm.pid = 0
		}
	}()

	return nil
}

func (pm *ProcessManager) Stop() error {

	// Lock the mutex for the thread safety
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	// If the process is not running there is no need to stop it so return an error
	if pm.state != StateRunning || pm.cmd == nil {
		return fmt.Errorf("process is not running")
	}

	//Try the graceful termination of the process
	// SIGTERM is the signal to terminate a process gracefully
	if err := pm.cmd.Process.Signal((syscall.SIGTERM)); err != nil {
		log.Printf("Failed to send SIGTERM to process: %v", err)

		if err := pm.cmd.Process.Kill(); err != nil {
			log.Printf("Failed to kill process: %v", err)
			return fmt.Errorf("failed to kill process: %v", err)
		}
	}

	pm.state = StateStopped
	pm.pid = 0
	fmt.Println("Process stopped")
	return nil
}

func (pm *ProcessManager) Restart() error {
	if pm.state != StateRunning {
		return fmt.Errorf("process is not running")
	}

	if err := pm.Stop(); err != nil {
		return fmt.Errorf("failed to stop process: %v", err)
	}

	// Wait for the short time to ensure the process is stopped
	time.Sleep(100 * time.Millisecond)

	return pm.Start()

}

func (pm *ProcessManager) Status() (ProcessState, int) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	return pm.state, pm.pid
}
