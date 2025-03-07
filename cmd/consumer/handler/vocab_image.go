package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/hibiken/asynq"
	"github.com/npvu1510/crawl-en-vocab/internal/model"
	"github.com/npvu1510/crawl-en-vocab/pkg/gcs"
	"github.com/npvu1510/crawl-en-vocab/pkg/utils"
)

func (h *Handler) HandlerVocabImageTask(ctx context.Context, t *asynq.Task) error {

	// UNMARSHAL
	var d *model.Dictionary
	if err := json.Unmarshal(t.Payload(), &d); err != nil {
		return fmt.Errorf("HandleAudioTask Unmarshal failed: %v", err)
	}

	// HANDLE TASK (STORAGE IMAGE)
	gcsService := gcs.Service

	filename := utils.Definition2FileName(d.Definition) + ".svg"
	fullFilePath := filepath.Join("./", h.Config.App.LOCAL_IMAGE_PATH, filename)

	imagePath, err := utils.DownloadImage(d.Image, fullFilePath)
	utils.CheckError(err)

	gcsService.SetBucketPath(h.Config.GCS.IMAGE_PATH)
	publicUrl, err := gcsService.UploadFile(imagePath, filename)
	if err != nil {
		return fmt.Errorf("HandleAudioTask storage image failed: %v", err)
	}

	// Update db with publicUrl
	h.DictionaryService.UpdateImage(d, publicUrl)

	return nil
}
