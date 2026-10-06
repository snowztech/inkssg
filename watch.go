package inkssg

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// watch rebuilds the site whenever a source file under dir changes.
// It debounces bursts of events (editors often emit several per save).
// outAbs is the absolute path of the build output, which is excluded.
func watch(ctx context.Context, dir, outAbs string, onChange func() error) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("watcher: %w", err)
	}
	defer w.Close()

	if err := addWatchTree(w, dir, outAbs); err != nil {
		return err
	}

	const debounce = 100 * time.Millisecond
	var timer *time.Timer
	trigger := make(chan struct{}, 1)

	for {
		select {
		case <-ctx.Done():
			return nil
		case ev, ok := <-w.Events:
			if !ok {
				return nil
			}
			if shouldIgnore(ev.Name, outAbs) {
				continue
			}
			// Chmod means only attributes changed, not content. On macOS a
			// rebuild is followed by another Chmod on the page it just read,
			// so reacting to it rebuilds in a loop.
			if ev.Op == fsnotify.Chmod {
				continue
			}
			if ev.Op&fsnotify.Create != 0 {
				if info, err := os.Stat(ev.Name); err == nil && info.IsDir() {
					_ = addWatchTree(w, ev.Name, outAbs)
				}
			}
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(debounce, func() {
				select {
				case trigger <- struct{}{}:
				default:
				}
			})
		case err, ok := <-w.Errors:
			if !ok {
				return nil
			}
			fmt.Fprintf(os.Stderr, "watch error: %v\n", err)
		case <-trigger:
			if err := onChange(); err != nil {
				fmt.Fprintf(os.Stderr, "rebuild failed: %v\n", err)
			}
		}
	}
}

func addWatchTree(w *fsnotify.Watcher, root, outAbs string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			return nil
		}
		if shouldIgnore(path, outAbs) {
			return filepath.SkipDir
		}
		return w.Add(path)
	})
}

func shouldIgnore(path, outAbs string) bool {
	abs, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	if abs == outAbs || strings.HasPrefix(abs, outAbs+string(filepath.Separator)) {
		return true
	}
	base := filepath.Base(abs)
	if strings.HasPrefix(base, ".") && base != "." {
		return true
	}
	return false
}
