package mdlsvc

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/tools/go/packages"

	"goa.design/clue/log"
	"goa.design/model/codegen"
)

// watch implements functionality to listen to changes in the model files
// when notifications are received from the filesystem, the model is rebuild
// and the editor page is refreshed via live reload server `lrserver`
func watch(ctx context.Context, dir, pkg string, reload func()) error {
	pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedFiles, Dir: dir}, pkg+"//...")
	if err != nil {
		return err
	}
	if len(pkgs) == 0 {
		fmt.Println("Nothing to watch")
		return nil
	}
	fmt.Println("Watching:", filepath.Dir(pkgs[0].GoFiles[0]))

	// Watch model design and regenerate on change
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	for _, p := range pkgs { // we need to watch the subpackages too
		if err = watcher.Add(filepath.Dir(p.GoFiles[0])); err != nil {
			return err
		}
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				if strings.HasPrefix(filepath.Base(ev.Name), codegen.TmpDirPrefix) {
					// ignore temporary (generated) files
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

				log.Infof(ctx, ev.String())
				reload()

			case err := <-watcher.Errors:
				log.Errorf(ctx, err, "Error watching files")
			}
		}
	}()

	return nil
}
