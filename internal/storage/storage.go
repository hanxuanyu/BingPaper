package storage

import (
	"context"
	"io"
)

type StoredObject struct {
	Key         string
	ContentType string
	Size        int64
	PublicURL   string
}

type Storage interface {
	Put(ctx context.Context, key string, r io.Reader, contentType string) (StoredObject, error)
	Get(ctx context.Context, key string) (io.ReadCloser, string, error)
	Delete(ctx context.Context, key string) error
	PublicURL(key string) (string, bool)
}

var GlobalStorage Storage

func InitStorage() error {
	// 实际初始化在 main.go 中根据配置调用对应的初始化函数
	return nil
}
