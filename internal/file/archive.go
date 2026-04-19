package file

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)


var noCompressExt = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
	".mp3": true, ".mp4": true, ".mkv": true, ".avi": true, ".mov": true,
	".zip": true, ".gz": true, ".rar": true, ".7z": true,
	".pdf": true,
}

func shouldCompress(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return !noCompressExt[ext]
}


func CreateArchive(paths []string) (string, error) {
	tmpFile, err := os.CreateTemp("", "winddrop_*.zip")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	zipWriter := zip.NewWriter(tmpFile)
	defer zipWriter.Close()

	for _, path := range paths {
		err := addToZip(zipWriter, path)
		if err != nil {
			return "", err
		}
	}

	return tmpFile.Name(), nil
}

func addToZip(zipWriter *zip.Writer, basePath string) error {
	return filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		relPath, err := filepath.Rel(filepath.Dir(basePath), path)
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = relPath

	
		if shouldCompress(path) {
			header.Method = zip.Deflate
		} else {
			header.Method = zip.Store
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		return err
	})
}
