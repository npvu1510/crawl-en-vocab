package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/hibiken/asynq"
	"github.com/npvu1510/crawl-en-vocab/internal/model"
	"github.com/npvu1510/crawl-en-vocab/pkg/gcs"
	"github.com/npvu1510/crawl-en-vocab/pkg/tts"
	"github.com/npvu1510/crawl-en-vocab/pkg/utils"
)

func (h *Handler) HandlerVocabAudioTask(ctx context.Context, t *asynq.Task) error {
	var wg sync.WaitGroup

	// UNMARSHAL
	var d *model.Dictionary
	if err := json.Unmarshal(t.Payload(), &d); err != nil {
		return fmt.Errorf("HandleAudioTask Unmarshal failed: %v", err)
	}

	// HANDLE TASK (STORE AUDIO)
	wg.Add(2)

	// GB
	go func() {
		defer wg.Done()

		filename := utils.Definition2FileName(d.Definition) + "_gb.mp3"
		publicUrl, err := ttsAndUpload(d.Definition, "en-GB", filename, h.Config.GCS.AUDIO_PATH)
		utils.CheckError(err)

		h.DictionaryService.UpdateAudioGb(d, publicUrl)
	}()

	// US
	go func() {
		defer wg.Done()

		filename := utils.Definition2FileName(d.Definition) + "_us.mp3"
		publicUrl, err := ttsAndUpload(d.Definition, "en-US", filename, h.Config.GCS.AUDIO_PATH)
		utils.CheckError(err)

		h.DictionaryService.UpdateAudioUs(d, publicUrl)
	}()

	wg.Wait()
	return nil
}

func ttsAndUpload(text, languageCode, filename, gcs_path string) (string, error) {
	ttsService := tts.Service
	gcsService := gcs.Service

	gcsService.SetBucketPath(gcs_path)
	ttsService.SetLanguageCode(languageCode)

	speechData, err := ttsService.TextToSpeech(text)
	utils.CheckError(err)
	publicUrl, err := gcsService.UploadBlob(speechData, filename)
	fmt.Println(publicUrl)

	if err != nil {
		return "", fmt.Errorf("HandleAudioTask storage audio failed: %v", err)

	}
	return publicUrl, nil
}

// languageCode := "en-US"
// if lang == "gb" {
// 	languageCode = "en-GB"
// }
