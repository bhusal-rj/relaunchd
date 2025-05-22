package watcher

import (
	"bhusal-rj/relaunchd/internal/config"
	"log"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	fsWatcher *fsnotify.Watcher   // File system watcher instance
	config    *config.WatchConfig // Configuration from the yaml file
	mu        sync.Mutex          // Mutex for thread-safe operations
	isRunning bool                // Flag to check if the watcher is running
	stopChan  chan struct{}       // Channel to signal stopping the watcher
}

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
		return nil // Alread the watcher is running
	}

	for _, path := range w.config.Paths {
		err := w.fsWatcher.Add(path)
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
			// Handle the event
			if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {

				log.Print("There has been an event")
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
