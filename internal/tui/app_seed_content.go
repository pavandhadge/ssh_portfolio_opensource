package tui

import (
	"embed"
	"io/fs"
	"runtime"
	"strings"
	"sync"
)

//go:embed content/**/*.txt
var seedContentFS embed.FS
var (
	seedContentOnce  sync.Once
	seedContentCache map[string]string
)

func preloadSeedContent() {
	seedContentCache = make(map[string]string)

	paths := make([]string, 0, 24)
	_ = fs.WalkDir(seedContentFS, "content", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d == nil || d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".txt") {
			paths = append(paths, path)
		}
		return nil
	})

	if len(paths) == 0 {
		return
	}

	workers := runtime.NumCPU()
	if workers < 2 {
		workers = 2
	}
	if workers > len(paths) {
		workers = len(paths)
	}

	type filePayload struct {
		path string
		text string
	}

	pathCh := make(chan string, len(paths))
	outCh := make(chan filePayload, len(paths))
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range pathCh {
				b, err := seedContentFS.ReadFile(path)
				if err != nil {
					continue
				}
				s := strings.TrimSpace(string(b))
				if s == "" {
					continue
				}
				outCh <- filePayload{path: path, text: s}
			}
		}()
	}

	for _, path := range paths {
		pathCh <- path
	}
	close(pathCh)
	wg.Wait()
	close(outCh)

	for item := range outCh {
		seedContentCache[item.path] = item.text
	}
}

func contentText(path, fallback string) string {
	seedContentOnce.Do(preloadSeedContent)
	s := seedContentCache[path]
	if s == "" {
		return fallback
	}
	return s
}
