package viper

import (
	"os"
	"path/filepath"
	"github.com/fsnotify/fsnotify"
)

func (v *Viper) watchConfig() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}

	configDir := filepath.Dir(v.configFilePath)
	configFile := filepath.Base(v.configFilePath)

	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if filepath.Base(event.Name) == configFile && (event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create) {
					v.reloadConfig()
				}
			case <-watcher.Errors:
				return
			}
		}
	}()

	watcher.Add(configDir)
}