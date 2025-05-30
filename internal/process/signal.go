package process

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

// SetupSignalHandling configures platform-specific signal handling
func SetupSignalHandling(callback func()) {
	sigs := make(chan os.Signal, 1)

	if runtime.GOOS == "windows" {
		// Windows supports Interrupt (Ctrl+C) but not SIGTERM
		signal.Notify(sigs, os.Interrupt)
	} else {
		// Unix-like systems support both SIGINT and SIGTERM
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	}

	go func() {
		<-sigs
		callback()
	}()
}
