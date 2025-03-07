package consumer

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"

// 	"github.com/hibiken/asynq"
// 	"github.com/npvu1510/crawl-en-vocab/internal"
// 	"github.com/npvu1510/crawl-en-vocab/internal/model"
// 	"github.com/npvu1510/crawl-en-vocab/internal/service"
// 	"github.com/npvu1510/crawl-en-vocab/pkg/config"
// 	"github.com/spf13/cobra"
// )

// var VocabAudioConsumerCmd = &cobra.Command{
// 	Use:   "vocab-audio-consumer",
// 	Short: "Consume vocabulary audio tasks",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		internal.Invoke(consumerCmd).Start(context.Background())
// 	},
// }

// type VocabAudioConsumer struct {
// 	Service service.IDictionaryService
// 	Config  *config.Config
// }

// func (v *VocabAudioConsumer) handlerVocabAudioTask(ctx context.Context, t *asynq.Task) error {

// 	// UNMARSHAL
// 	var d *model.Dictionary
// 	if err := json.Unmarshal(t.Payload(), &d); err != nil {
// 		return fmt.Errorf("HandleAudioTask Unmarshal failed: %v", err)
// 	}

// 	// // HANDLE TASK (STORAGE IMAGE)
// 	// gcsService := gcs.Service

// 	// filename := utils.Definition2FileName(d.Definition) + ".svg"
// 	// fullFilePath := filepath.Join(v.Config.App.LOCAL_IMAGE_PATH, filename)

// 	// audioPath, err := utils.DownloadAudio(d.Audio, "./"+fullFilePath)
// 	// utils.CheckError(err)

// 	// publicUrl, err := gcsService.UploadFile(audioPath, "/vu1/audios/"+filename)
// 	// if err != nil {
// 	// 	return fmt.Errorf("HandleAudioTask storage audio failed: %v", err)
// 	// }

// 	// // Update db with publicUrl
// 	// v.Service.UpdateAudio(d, publicUrl)

// 	// return nil
// }
