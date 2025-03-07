package gcs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"github.com/npvu1510/crawl-en-vocab/pkg/utils"
	"google.golang.org/api/option"
)

const (
	PROJECT_ID               = "saoladigital"            // FILL IN WITH YOURS
	BUCKET_NAME              = "static.saoladigital.com" // FILL IN WITH YOURS
	TTS_SERVICE_ACCOUNT_PATH = "service-account-tts.json"
	GCS_SERVICE_ACCOUNT_PATH = "service-account-gcs.json"
)

type GCSService struct {
	client     *storage.Client
	projectID  string
	bucketName string
	uploadPath string
}

var (
	Service *GCSService
	once    sync.Once
)

func init() {
	once.Do(func() {
		// os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "") // FILL IN WITH YOUR FILE PATH
		client, err := storage.NewClient(context.Background(), option.WithCredentialsFile(GCS_SERVICE_ACCOUNT_PATH))
		utils.CheckError(err)
		defer client.Close()

		Service = &GCSService{
			client:     client,
			bucketName: BUCKET_NAME,
			projectID:  PROJECT_ID,
			uploadPath: "public/npvu1510",
		}

	})
}

func GetClient() *storage.Client {
	return Service.client
}

func (c *GCSService) Close() {
	if Service != nil {
		Service.client.Close()
	}
}

func (g *GCSService) SetBucketPath(path string) {
	g.uploadPath = path
}

// UploadFile uploads an object
func (c *GCSService) UploadBlob(file []byte, filename string) (string, error) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// file, err := os.Open(filePath)
	// if err != nil {
	// 	return "",fmt.Errorf("UploadFile Open failed: %w", err)
	// }
	// defer file.Close()

	buf := bytes.NewBuffer(file)

	ext := filepath.Ext(filename)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Upload an object with storage.Writer.
	fullFilePath := c.uploadPath + "/" + filename
	wc := c.client.Bucket(c.bucketName).Object(fullFilePath).NewWriter(ctx)
	wc.ContentType = mimeType
	if _, err := io.Copy(wc, buf); err != nil {
		return "", fmt.Errorf("UploadBlob Copy failed: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("UploadBlob Close failed: %v", err)
	}

	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.bucketName, fullFilePath)
	return publicURL, nil
}

func (c *GCSService) UploadFile(filePath string, filename string) (string, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("UploadFile Open failed: %w", err)
	}
	defer file.Close()

	ext := filepath.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream" // Giá trị mặc định nếu không nhận diện được
	}

	// Tạo writer và đặt Content-Type
	fullFilePath := c.uploadPath + "/" + filename
	wc := c.client.Bucket(c.bucketName).Object(fullFilePath).NewWriter(ctx)
	wc.ContentType = mimeType // Đặt Content-Type chính xác

	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("UploadFile Copy failed: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("UploadFile Close failed: %v", err)
	}
	fmt.Printf("✅ Uploaded successfully: %s (Content-Type: %s)\n", fullFilePath, mimeType)

	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.bucketName, fullFilePath)
	return publicURL, nil
}
