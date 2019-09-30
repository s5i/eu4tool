// eu4autosave is a tool which clones Ironman backup files into dated versions
// whenever they change.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/s5i/eu4tool/lib/decode"
	"github.com/s5i/eu4tool/lib/unzip"
)

var backupPath = flag.String("backup_path", "", fmt.Sprintf(`Save files directory. Defaults to %s`, defaultBackupPath))

func main() {
	flag.Parse()

	path := *backupPath

	if path == "" {
		path = defaultBackupPath
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create a filesystem watcher: %v", err)
	}
	defer w.Close()

	if err := w.Add(path); err != nil {
		log.Fatalf("Failed to watch %s for changes: %v", path, err)
	}

	log.Printf("Watching %s for *_Backup.eu4 changes...", path)

	for e := range w.Events {
		if e.Op != fsnotify.Create || !strings.HasSuffix(e.Name, "_Backup.eu4") {
			continue
		}

		// Ugly hack to avoid read errors just after file creation.
		time.Sleep(time.Second)

		dated, err := clone(e.Name)
		if err != nil {
			log.Printf("Failed to backup %s: %v", filepath.Base(e.Name), err)
			continue
		}

		log.Printf("%s backed up as %s.", filepath.Base(e.Name), filepath.Base(dated))

		baseSave := strings.TrimSuffix(e.Name, "_Backup.eu4") + ".eu4"
		err = touch(baseSave)
		if err != nil {
			log.Printf("Failed to touch %s: %v", filepath.Base(baseSave), err)
			continue
		}
		log.Printf("%s touched", filepath.Base(baseSave))
	}
}

func clone(path string) (string, error) {
	tmp := fmt.Sprintf("%s.tmp", path)
	if err := cp(path, tmp); err != nil {
		return "", fmt.Errorf("cp(%s, %s) failed: %v", path, tmp, err)
	}
	defer os.Remove(tmp)

	date, err := readDate(tmp)
	if err != nil {
		return "", fmt.Errorf("readDate(%s) failed: %v", tmp, err)
	}

	dated := strings.Replace(tmp, "_Backup.eu4.tmp", fmt.Sprintf("_%s.eu4", date), -1)
	if err := os.Rename(tmp, dated); err != nil {
		return "", fmt.Errorf("os.Rename(%s, %s) failed: %v", tmp, dated, err)
	}

	return dated, nil
}

func cp(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	if err != nil {
		return err
	}
	return d.Close()
}

func readDate(path string) (string, error) {
	b, err := unzip.Meta(path)
	if err != nil {
		return "", fmt.Errorf("unzip.Meta(%s) failed: %v", path, err)
	}

	date, err := decode.DateFromBinaryMeta(b)
	if err != nil {
		return "", fmt.Errorf("decode.DateFromBinaryMeta failed: %v", err)
	}

	return date, nil
}

func touch(path string) error {
	now := time.Now()
	err := os.Chtimes(path, now, now)
	if err != nil {
		return fmt.Errorf("os.Chtimes(%s, %v, %v): %v", path, now, now, err)
	}
	return nil
}
