package loader

import (
	"os"
	"time"
)

// Watch polls the given file paths at the provided interval and emits the
// path on the returned channel whenever the file's modification time changes.
// The second return value is a stop function; call it to terminate the
// watcher goroutine and close the channel.
//
// It uses simple polling instead of fsnotify to avoid adding heavy
// dependencies for the game-jam scope and to keep cross-platform behaviour
// predictable.
func Watch(paths []string, interval time.Duration) (<-chan string, func()) {
    ch := make(chan string)
    stop := make(chan struct{})
    // Record initial mod times (zero value if file does not exist yet).
    modTimes := make(map[string]time.Time)
    for _, p := range paths {
        if fi, err := os.Stat(p); err == nil {
            modTimes[p] = fi.ModTime()
        } else {
            modTimes[p] = time.Time{}
        }
    }

    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()
        defer close(ch)
        for {
            select {
            case <-ticker.C:
                for _, p := range paths {
                    fi, err := os.Stat(p)
                    if err != nil {
                        // Ignore missing file errors during polling.
                        continue
                    }
                    if fi.ModTime().After(modTimes[p]) {
                        modTimes[p] = fi.ModTime()
                        ch <- p
                    }
                }
            case <-stop:
                return
            }
        }
    }()

    // Provide a stop function.
    return ch, func() { close(stop) }
} 