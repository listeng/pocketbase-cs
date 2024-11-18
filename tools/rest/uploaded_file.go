package rest

import (
	"net/http"

	"github.com/pocketbase/pocketbase/tools/filesystem"
)

// DefaultMaxMemory defines the default max memory bytes that
// will be used when parsing a form request body.
const DefaultMaxMemory = 32 << 20 // 32mb

// FindUploadedFiles extracts all form files of "key" from a http request
// and returns a slice with filesystem.File instances (if any).
func FindUploadedFiles(r *http.Request, key string) ([]*filesystem.File, error) {
	if r.MultipartForm == nil {
		err := r.ParseMultipartForm(DefaultMaxMemory)
		if err != nil {
			return nil, err
		}
	}

	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		return nil, http.ErrMissingFile
	}

	// 检查两个key是否都没有文件
	keyFiles := r.MultipartForm.File[key]
	keyArrayFiles := r.MultipartForm.File[key+"[]"]
	if len(keyFiles) == 0 && len(keyArrayFiles) == 0 {
		return nil, http.ErrMissingFile
	}

	// 预分配足够的容量
	result := make([]*filesystem.File, 0, len(keyFiles)+len(keyArrayFiles))

	// 处理 key 的文件
	for _, fh := range keyFiles {
		file, err := filesystem.NewFileFromMultipart(fh)
		if err != nil {
			return nil, err
		}
		result = append(result, file)
	}

	// 处理 key[] 的文件
	for _, fh := range keyArrayFiles {
		file, err := filesystem.NewFileFromMultipart(fh)
		if err != nil {
			return nil, err
		}
		result = append(result, file)
	}

	return result, nil
}
