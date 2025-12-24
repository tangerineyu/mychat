package upload

import (
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
)

func SaveFile(file *multipart.FileHeader, dstFolder string) (string, error) {
	if _, err := os.Stat(dstFolder); os.IsNotExist(err) {
		_ = os.MkdirAll(dstFolder, os.ModePerm)
	}
	ext := path.Ext(file.Filename)
	newFileName := uuid.New().String() + ext
	dstPath := path.Join(dstFolder, newFileName)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			return
		}
	}(src)

	out, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			return
		}
	}(out)
	_, err = io.Copy(out, src)
	return "/" + dstPath, err
}
func CheckImageExt(fileName string) bool {
	ext := strings.ToLower(path.Ext(fileName))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp":
		return true
	default:
		return false
	}
}
