package webdav

import (
	"context"
	"fmt"
	"io"
	"path"
	"strings"

	"BingPaper/internal/storage"

	"github.com/studio-b12/gowebdav"
)

type WebDAVStorage struct {
	client          *gowebdav.Client
	publicURLPrefix string
}

func NewWebDAVStorage(url, username, password, publicURLPrefix string) (*WebDAVStorage, error) {
	client := gowebdav.NewClient(url, username, password)
	if err := client.Connect(); err != nil {
		// 有些 webdav 不支持 Connect，我们可以忽略错误或者做简单的探测
	}
	return &WebDAVStorage{
		client:          client,
		publicURLPrefix: publicURLPrefix,
	}, nil
}

func (w *WebDAVStorage) Put(ctx context.Context, key string, r io.Reader, contentType string) (storage.StoredObject, error) {
	// 确保目录存在
	dir := path.Dir(key)
	if dir != "." && dir != "/" {
		if err := w.client.MkdirAll(dir, 0755); err != nil {
			return storage.StoredObject{}, err
		}
	}

	err := w.client.WriteStream(key, r, 0644)
	if err != nil {
		return storage.StoredObject{}, err
	}

	publicURL := ""
	if w.publicURLPrefix != "" {
		publicURL = fmt.Sprintf("%s/%s", strings.TrimSuffix(w.publicURLPrefix, "/"), key)
	}

	return storage.StoredObject{
		Key:         key,
		ContentType: contentType,
		PublicURL:   publicURL,
	}, nil
}

func (w *WebDAVStorage) Get(ctx context.Context, key string) (io.ReadCloser, string, error) {
	reader, err := w.client.ReadStream(key)
	if err != nil {
		return nil, "", err
	}
	return reader, "", nil
}

func (w *WebDAVStorage) Delete(ctx context.Context, key string) error {
	return w.client.Remove(key)
}

func (w *WebDAVStorage) PublicURL(key string) (string, bool) {
	if w.publicURLPrefix != "" {
		return fmt.Sprintf("%s/%s", strings.TrimSuffix(w.publicURLPrefix, "/"), key), true
	}
	return "", false
}

func (w *WebDAVStorage) Exists(ctx context.Context, key string) (bool, error) {
	_, err := w.client.Stat(key)
	if err == nil {
		return true, nil
	}
	// gowebdav 的错误处理比较原始，通常 404 会返回错误
	// 这里假设报错就是不存在，或者可以根据错误消息判断
	if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
		return false, nil
	}
	return false, err
}
