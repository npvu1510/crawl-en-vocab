package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
	return
}

func DownloadImage(url string, fullFilePath string) (string, error) {
	//
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download image failed: %d", resp.StatusCode)
	}

	//
	dir := filepath.Dir(fullFilePath)

	//
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return "", err
	}

	//
	file, err := os.Create(fullFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	//
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	fileName := filepath.Base(fullFilePath)
	fmt.Printf("✅ Save image %s successfully\n", fileName)
	return fullFilePath, nil
}

func Definition2FileName(definition string) string {
	return strings.ReplaceAll(strings.ReplaceAll(definition, " ", "_"), "-", "_")
}
