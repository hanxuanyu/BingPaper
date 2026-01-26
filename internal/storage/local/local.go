package local

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"BingDailyImage/internal/storage"
)

type LocalStorage struct {
	root string
}

func NewLocalStorage(root string) (*LocalStorage, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, err
	}
	return &LocalStorage{root: root}, nil
}

func (l *LocalStorage) Put(ctx context.Context, key string, r io.Reader, contentType string) (storage.StoredObject, error) {
	path := filepath.Join(l.root, key)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return storage.StoredObject{}, err
	}

	f, err := os.Create(path)
	if err != nil {
		return storage.StoredObject{}, err
	}
	defer f.Close()

	n, err := io.Copy(f, r)
	if err != nil {
		return storage.StoredObject{}, err
	}

	return storage.StoredObject{
		Key:         key,
		ContentType: contentType,
		Size:        n,
	}, nil
}

func (l *LocalStorage) Get(ctx context.Context, key string) (io.ReadCloser, string, error) {
	path := filepath.Join(l.root, key)
	f, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	// 这里很难从文件扩展名以外的地方获得 contentType，除非存储时记录
	// 简单处理
	return f, "", nil
}

func (l *LocalStorage) Delete(ctx context.Context, key string) error {
	path := filepath.Join(l.root, key)
	return os.Remove(path)
}

func (l *LocalStorage) PublicURL(key string) (string, bool) {
	return "", false
}
