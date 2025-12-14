package jsonstore

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Store struct {
	Path string
}

func New(path string) Store {
	return Store{Path: path}
}

func (s Store) Load() (fileSchema, error) {
	b, err := os.ReadFile(s.Path)
	if err != nil {
		if os.IsNotExist(err) {
			// empty store
			return fileSchema{Version: schemaVersion, SavedAt: time.Now().UTC(), Todos: []todoRow{}}, nil
		}
		return fileSchema{}, err
	}

	var fs fileSchema
	if err := json.Unmarshal(b, &fs); err != nil {
		return fileSchema{}, ErrCorruptData
	}
	if fs.Version != schemaVersion {
		return fileSchema{}, ErrCorruptData
	}
	if fs.Todos == nil {
		fs.Todos = []todoRow{}
	}
	return fs, nil
}

func (s Store) Save(fs fileSchema) error {
	fs.Version = schemaVersion
	fs.SavedAt = time.Now().UTC()

	dir := filepath.Dir(s.Path)
	tmp := s.Path + ".tmp"

	b, err := json.MarshalIndent(fs, "", "")
	if err != nil {
		return err
	}

	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}

	if _, err := f.Write(b); err != nil {
		_ = f.Close()
		return err
	}

	if err := f.Sync(); err != nil {
		_ = f.Close()
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmp, s.Path); err != nil {
		return err
	}

	// best-effort: sync directory
	if d, err := os.Open(dir); err == nil {
		_ = d.Sync()
		_ = d.Close()
	}

	return nil
}
