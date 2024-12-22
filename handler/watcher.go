package handler

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/fsnotify/fsnotify"
	"github.com/networkteam/slogutils"

	"github.com/esdete2/mjml-dev/config"
)

type Watcher struct {
	processor    *Processor
	fsWatcher    *fsnotify.Watcher
	config       *config.Config
	debounceTime time.Duration
	mu           sync.Mutex
	timer        *time.Timer
	notifier     ReloadNotifier
	done         chan struct{}
}

type ReloadNotifier interface {
	NotifyReload()
}

func NewWatcher(proc *Processor, cfg *config.Config, notifier ReloadNotifier) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, errors.Wrap(err, "creating fsnotify watcher")
	}

	w := &Watcher{
		processor:    proc,
		fsWatcher:    fsWatcher,
		config:       cfg,
		debounceTime: 100 * time.Millisecond,
		notifier:     notifier,
		done:         make(chan struct{}),
	}

	if err := w.addDirsToWatch(); err != nil {
		fsWatcher.Close()
		return nil, err
	}

	return w, nil
}

func (w *Watcher) addDirsToWatch() error {
	// Watch documents directory
	if err := filepath.Walk(w.config.Paths.Documents, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return w.fsWatcher.Add(path)
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "watching documents directory")
	}

	// Watch partials directory
	if err := filepath.Walk(w.config.Paths.Partials, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return w.fsWatcher.Add(path)
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "watching partials directory")
	}

	return nil
}

func (w *Watcher) Watch() error {
	slog.With("documents", w.config.Paths.Documents).With("partials", w.config.Paths.Partials).Info("Watching for changes")

	go func() {
		defer w.fsWatcher.Close()
		for {
			select {
			case event, ok := <-w.fsWatcher.Events:
				if !ok {
					return
				}

				// Skip temporary files and non-MJML files
				if strings.HasPrefix(filepath.Base(event.Name), ".") ||
					!strings.HasSuffix(event.Name, ".mjml") {
					continue
				}

				// Handle the file event
				if err := w.handleFileEvent(event); err != nil {
					slog.Error("Error handling file event", slogutils.Err(err))
				}

			case err, ok := <-w.fsWatcher.Errors:
				if !ok {
					return
				}
				slog.Error("File watcher error", slogutils.Err(err))

			case <-w.done:
				return
			}
		}
	}()

	return nil
}

func (w *Watcher) Stop() {
	close(w.done)
}

func (w *Watcher) handleFileEvent(event fsnotify.Event) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Cancel previous timer if it exists
	if w.timer != nil {
		w.timer.Stop()
	}

	// Ignore chmod events
	if event.Op == fsnotify.Chmod {
		return nil
	}

	// Set new timer for debouncing
	w.timer = time.AfterFunc(w.debounceTime, func() {
		isPartial := strings.HasPrefix(event.Name, w.config.Paths.Partials)

		// If it's a partial or create/remove/rename operation, rebuild all templates
		if isPartial || event.Op&(fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
			slog.Info("Rebuilding all templates...")
			if err := w.processor.Process(); err != nil {
				slog.Error("Error rebuilding templates", slogutils.Err(err))
			} else {
				w.notifier.NotifyReload()
			}
			return
		}

		// For document write changes, rebuild only the changed template
		if event.Op&fsnotify.Write != 0 {
			relPath, err := filepath.Rel(w.config.Paths.Documents, event.Name)
			if err != nil {
				slog.Error("Error getting relative path", slogutils.Err(err))
				return
			}

			templateName := strings.TrimSuffix(filepath.ToSlash(relPath), ".mjml")
			slog.With("template", templateName).Info("Rebuilding template")

			if err := w.processor.ProcessSingle(templateName); err != nil {
				slog.Error("Error rebuilding single template", slogutils.Err(err))
			} else {
				w.notifier.NotifyReload()
			}
		}
	})

	return nil
}
