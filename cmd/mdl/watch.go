package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jaschaephraim/lrserver"
	"golang.org/x/tools/go/packages"
)

// watch implements functionality to listen to changes in the model files
// when notifications are received from the filesystem, the model is rebuild
// and the editor page is refreshed via live reload server `lrserver`
func watch(pkg string, reload func()) error {
	// Watch model design and regenerate on change
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedFiles}, pkg+"//...")
	if err != nil {
		return err
	}
	fmt.Println("Watching:", filepath.Dir(pkgs[0].GoFiles[0]))
	for _, p := range pkgs { // we need to watch the subpackages too
		if err = watcher.Add(filepath.Dir(p.GoFiles[0])); err != nil {
			return err
		}
	}

	// Create live reload server and hookup to watcher
	lr := lrserver.New(lrserver.DefaultName, lrserver.DefaultPort)
	lr.SetStatusLog(nil)
	lr.SetErrorLog(nil)
	go func() {
		if err := lr.ListenAndServe(); err != nil {
			fail(err.Error())
		}
	}()
	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				if strings.HasPrefix(filepath.Base(ev.Name), tmpDirPrefix) {
					continue
				}

				// debounce, because some editors do several file operations when you save
				// we wait for the stream of events to become silent for `interval`
				interval := 100 * time.Millisecond
				timer := time.NewTimer(interval)
			outer:
				for {
					select {
					case ev = <-watcher.Events:
						timer.Reset(interval)
					case <-timer.C:
						break outer
					}
				}

				fmt.Println(ev.String())
				reload()
				lr.Reload(ev.Name)

			case err := <-watcher.Errors:
				fmt.Fprintln(os.Stderr, "Error watching files:", err)
			}
		}
	}()

	return nil
}
