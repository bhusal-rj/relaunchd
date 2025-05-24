package watcher

import (
	"bhusal-rj/relaunchd/internal/config"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	fsWatcher     *fsnotify.Watcher   // File system watcher instance
	config        *config.WatchConfig // Configuration from the yaml file
	mu            sync.Mutex          // Mutex for thread-safe operations
	isRunning     bool                // Flag to check if the watcher is running
	stopChan      chan struct{}       // Channel to signal stopping the watcher
	changeHandler ChangeHandler       // Handler for file changes
}

// Change handler is the func that gets called when the file changes
type ChangeHandler func()

// Create the new file system watcher
func New(cfg *config.WatchConfig) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &Watcher{
		fsWatcher: fsWatcher,
		config:    cfg,
		isRunning: false,
		stopChan:  make(chan struct{}),
		mu:        sync.Mutex{},
	}, nil
}

// Start the monitoring of the specified paths
func (w *Watcher) Start() error {
	w.mu.Lock()

	if w.isRunning {
		w.mu.Unlock()
		return nil // Already the watcher is running
	}

	for _, path := range w.config.Paths {
		//Walk through the directory and add all the subdirectories
		err := filepath.Walk(path, func(subPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				// Check if directory should be excluded
				for _, pattern := range w.config.Exclude {
					matched, err := filepath.Match(pattern, info.Name())
					if err != nil {
						return err
					}
					if matched {
						log.Printf("Skipping excluded directory: %s (matched pattern: %s)", subPath, pattern)
						return filepath.SkipDir
					}
				}

				// Add directory to watcher after checking all exclusion patterns
				if err := w.fsWatcher.Add(subPath); err != nil {
					return err
				}
				// log.Printf("Watching directory: %s", subPath)
			}
			return nil
		})

		if err != nil {
			w.mu.Unlock()
			return err
		}

	}
	w.isRunning = true
	w.mu.Unlock()
	go w.watch() //Start the watching goroutine
	return nil
}

func (w *Watcher) watch() {
	defer func() {
		w.mu.Lock()
		w.isRunning = false
		w.mu.Unlock()
	}()

	for {
		select {
		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return
			}

			// Check if this file should be excluded based on patterns
			filename := filepath.Base(event.Name)
			excluded := false

			for _, pattern := range w.config.Exclude {
				matched, err := filepath.Match(pattern, filename)
				if err == nil && matched {
					// log.Printf("Ignoring excluded file: %s (matched pattern: %s)", event.Name, pattern)
					excluded = true
					break
				}
			}

			// Skip this event if the file is excluded
			if excluded {
				continue
			}

			// Handle the event for non-excluded files
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {

				log.Printf("File changes detected: %s", event.Name)

				// Call the change handler if it is set
				w.mu.Lock()
				handler := w.changeHandler
				w.mu.Unlock()

				if handler != nil {
					handler()
				}
			}

		case err, ok := <-w.fsWatcher.Errors:
			if !ok {
				return

			}
			log.Printf("Error: %v", err)
		case <-w.stopChan:
			// Stop the watcher
			w.fsWatcher.Close()
			return
		}
	}
}

func (w *Watcher) Stop() error {
	w.mu.Lock()
	if !w.isRunning {
		w.mu.Unlock()
		return nil // Alread the watcher is stopped
	}
	close(w.stopChan)
	w.mu.Unlock()
	return nil
}

// SetsChangeHandler sets the function to call when file changes
func (w *Watcher) SetChangeHandler(handler ChangeHandler) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.changeHandler = handler
}
