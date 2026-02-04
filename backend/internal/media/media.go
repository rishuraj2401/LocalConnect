package media

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SavedMedia struct {
	Path      string
	MediaType string
}

func SaveFile(ctx context.Context, dir string, file multipart.File, header *multipart.FileHeader) (SavedMedia, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return SavedMedia{}, err
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	path := filepath.Join(dir, filename)

	out, err := os.Create(path)
	if err != nil {
		return SavedMedia{}, err
	}
	defer out.Close()
	if _, err := io.Copy(out, file); err != nil {
		return SavedMedia{}, err
	}
	mediaType := "image"
	if strings.HasPrefix(header.Header.Get("Content-Type"), "video/") {
		mediaType = "video"
	}
	return SavedMedia{Path: path, MediaType: mediaType}, nil
}
