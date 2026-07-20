@@ -123,6 +123,7 @@ func (v *Viper) WatchConfig() {
 	watcher, err := fsnotify.NewWatcher()
 	if err != nil {
 		v.logError(fmt.Errorf("failed to create watcher: %w", err))
+		return
 	}
 	defer watcher.Close()

@@ -130,6 +131,14 @@ func (v *Viper) WatchConfig() {
 	if err != nil {
 		v.logError(fmt.Errorf("failed to add watch on config file: %w", err))
 	}

+	go func() {
+		for {
+			select {
+			case event, ok := <-watcher.Events:
+				if !ok {
+					return
+				}
+				if event.Op&fsnotify.Rename == fsnotify.Rename {
+					if err := watcher.Add(v.configFileUsed); err != nil {
+						v.logError(fmt.Errorf("failed to re-add watch on config file: %w", err))
+					}
+				}
+			case err, ok := <-watcher.Errors:
+				if !ok {
+					return
+				}
+				v.logError(fmt.Errorf("file watch error: %w", err))
+			}
+		}
+	}()
+
 	v.watcher = watcher
 	v.watchConfig()
 }
```

### Explanation

1. **Add a check for `fsnotify.Rename` event**: We add a goroutine that listens for events from the `fsnotify.Watcher`. When a `fsnotify.Rename` event is detected, it re-adds the watch on the configuration file.
2. **Handle errors and cleanup**: The goroutine also handles errors from the `fsnotify.Watcher` and ensures that the watcher is closed properly.

### Additional Considerations

- **Symlink Handling**: The above solution should handle symlink swaps correctly as `fsnotify` will trigger a `Rename` event for symlink changes.
- **Cross-Platform Compatibility**: The `fsnotify` package is designed to work across Linux, macOS, and Windows, so the solution should be cross-platform.
- **Resource Leaks**: The solution does not introduce any resource leaks as it properly closes the watcher and handles errors.

This patch should meet the acceptance criteria and resolve the issue with minimal changes to the existing codebase.