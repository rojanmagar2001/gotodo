package jsonstore

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type lock struct {
	dir string
}

func acquireLock(dbPath string) (*lock, error) {
	dir := dbPath + ".lockdir"

	if err := os.Mkdir(dir, 0o700); err != nil {
		if os.IsExist(err) {
			return nil, ErrLocked
		}
		return nil, err
	}

	// write debug info (best-effort)
	_ = os.WriteFile(filepath.Join(dir, "info.txt"),
		[]byte(fmt.Sprintf("pid=%d\nat=%s\n", os.Getpid(), time.Now().UTC().Format(time.RFC3339Nano))),
		0o600,
	)

	return &lock{dir: dir}, nil
}

func (l *lock) release() error {
	return os.RemoveAll(l.dir)
}
