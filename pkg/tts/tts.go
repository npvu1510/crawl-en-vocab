package tts

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

type ClientUploader struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	uploadPath string
}

var (
	Service *ClientUploader
	once    sync.Once
)

func init() {
	once.Do(func() {
		// os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "") // FILL IN WITH YOUR FILE PATH
		client, err := storage.NewClient(context.Background(), option.WithCredentialsFile(GCS_SERVICE_ACCOUNT_PATH))
		utils.CheckError(err)
		defer client.Close()

		Service = &ClientUploader{
			cl:         client,
			bucketName: BUCKET_NAME,
			projectID:  PROJECT_ID,
			uploadPath: "public/npvu1510",
		}

	})
}

func (c *ClientUploader) Close() {
	if Service != nil {
		Service.cl.Close()
	}
}

func GetClient() *storage.Client {
	return Service.cl
}

func SetBucketPath(path string) *ClientUploader {
	Service.uploadPath = path
	return Service
}

// UploadFile uploads an object
func (c *ClientUploader) UploadBlob(file []byte, object string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// file, err := os.Open(filePath)
	// if err != nil {
	// 	return fmt.Errorf("UploadFile Open failed: %w", err)
	// }
	// defer file.Close()

	buf := bytes.NewBuffer(file)

	// Upload an object with storage.Writer.
	wc := c.cl.Bucket(c.bucketName).Object(c.uploadPath + object).NewWriter(ctx)
	if _, err := io.Copy(wc, buf); err != nil {
		return fmt.Errorf("UploadBlob Copy failed: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("UploadBlob Close failed: %v", err)
	}

	return nil
}

func (c *ClientUploader) UploadFile(filePath string, object string) (string, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("UploadFile Open failed: %w", err)
	}
	defer file.Close()

	// Tự động nhận diện Content-Type dựa trên phần mở rộng của file
	ext := filepath.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream" // Giá trị mặc định nếu không nhận diện được
	}

	// Tạo writer và đặt Content-Type
	wc := c.cl.Bucket(c.bucketName).Object(c.uploadPath + object).NewWriter(ctx)
	wc.ContentType = mimeType // Đặt Content-Type chính xác

	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("UploadFile Copy failed: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("UploadFile Close failed: %v", err)
	}
	fmt.Printf("✅ Uploaded successfully: %s (Content-Type: %s)\n", object, mimeType)

	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.bucketName, c.uploadPath+object)
	return publicURL, nil
}
