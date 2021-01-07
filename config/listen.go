// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2021/1/7

package config

import (
	"log"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

const writeOrCreateMask = fsnotify.Write | fsnotify.Create

type FileWatch struct {
	Path string
}

func (fw FileWatch) WatchConfig(f func(fsnotify.Event)) {
	initWG := sync.WaitGroup{}
	initWG.Add(1)
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		eventsWG := sync.WaitGroup{}
		eventsWG.Add(1)
		go func() {
			defer eventsWG.Done()
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok { // 'Events' channel is closed
						return
					}
					if filepath.Clean(event.Name) == filepath.Clean(fw.Path) {
						if event.Op&writeOrCreateMask != 0 {
							if f != nil {
								f(event)
							}
						} else if event.Op&fsnotify.Remove&fsnotify.Remove != 0 {
							return
						}
					}

				case err, ok := <-watcher.Errors:
					if ok { // 'Errors' channel is not closed
						log.Printf("watcher error: %v", err)
					}
					return
				}
			}
		}()
		if err := watcher.Add(fw.Path); err != nil {
			log.Printf("add path error: %v", err)
		}
		initWG.Done()   // done initializing the watch in this go routine, so the parent routine can move on...
		eventsWG.Wait() // now, wait for event loop to end in this go-routine...
	}()
	initWG.Wait() // make sure that the go routine above fully ended before returning
}
