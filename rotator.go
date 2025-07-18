package golog

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Rotator handles log file rotation.
type Rotator struct {
	filePath   string
	maxSize    int64 // in bytes
	maxBackups int
	compress   bool
}

// NewRotator creates a new rotator.
func NewRotator(filePath string, maxSizeMB, maxBackups int, compress bool) *Rotator {
	return &Rotator{
		filePath:   filePath,
		maxSize:    int64(maxSizeMB) * 1024 * 1024,
		maxBackups: maxBackups,
		compress:   compress,
	}
}

// RotateIfNeeded rotates the log file if it exceeds the size limit.
func (r *Rotator) RotateIfNeeded(file *os.File) error {
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat log file: %v", err)
	}

	if info.Size() < r.maxSize {
		return nil
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("failed to close log file: %v", err)
	}

	newPath := fmt.Sprintf("%s.%s", r.filePath, time.Now().Format("20060102_150405"))
	if err := os.Rename(r.filePath, newPath); err != nil {
		return fmt.Errorf("failed to rename log file: %v", err)
	}

	if r.compress {
		if err := compressFile(newPath); err != nil {
			return fmt.Errorf("failed to compress log file: %v", err)
		}
		os.Remove(newPath)
		newPath += ".gz"
	}

	r.cleanupBackups()

	file, err = os.OpenFile(r.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to reopen log file: %v", err)
	}

	return nil
}

// compressFile compresses a file using gzip.
func compressFile(filePath string) error {
	in, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(filePath + ".gz")
	if err != nil {
		return err
	}
	defer out.Close()

	gz := gzip.NewWriter(out)
	defer gz.Close()

	_, err = io.Copy(gz, in)
	return err
}

// cleanupBackups removes old log files if the number exceeds maxBackups.
func (r *Rotator) cleanupBackups() {
	files, err := filepath.Glob(r.filePath + ".*")
	if err != nil {
		return
	}

	if len(files) <= r.maxBackups {
		return
	}

	// Sort files by modification time (newest first)
	type fileInfo struct {
		name  string
		mtime time.Time
	}
	var fileInfos []fileInfo
	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			continue
		}
		fileInfos = append(fileInfos, fileInfo{name: f, mtime: info.ModTime()})
	}

	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].mtime.After(fileInfos[j].mtime)
	})

	// Remove oldest files
	for i := r.maxBackups; i < len(fileInfos); i++ {
		os.Remove(fileInfos[i].name)
	}
}
