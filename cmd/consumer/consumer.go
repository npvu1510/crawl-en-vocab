package consumer

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/npvu1510/crawl-en-vocab/cmd/consumer/handler"
	"github.com/npvu1510/crawl-en-vocab/cmd/publisher"
	"github.com/npvu1510/crawl-en-vocab/internal"
	"github.com/npvu1510/crawl-en-vocab/internal/service"
	"github.com/npvu1510/crawl-en-vocab/pkg/asynqueue"
	"github.com/npvu1510/crawl-en-vocab/pkg/config"
	"github.com/npvu1510/crawl-en-vocab/pkg/utils"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var (
	ImageQueueName = "image_queue"
	AudioQueueName = "audio_queue"
)

var ConsumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Consume Asynq tasks",
	Args:  cobra.ExactArgs(1), // Yêu cầu đúng 1 argument (loại task)
	Run: func(cmd *cobra.Command, args []string) {
		consumerType := args[0] // Lấy argument từ CLI
		if consumerType != publisher.VocabImageTaskType && consumerType != publisher.VocabAudioTaskType {
			fmt.Println("Invalid consumer type. Use 'image' or 'audio'")
			return
		}

		internal.Invoke(func(lc fx.Lifecycle, conf *config.Config, dictionaryService service.IDictionaryService) {
			consumerCmd(lc, conf, dictionaryService, consumerType)
		}).Start(context.Background())
	},
}

func consumerCmd(
	lc fx.Lifecycle,
	conf *config.Config,
	dictionaryService service.IDictionaryService,
	consumerType string,
) {

	queueName := "default"

	if consumerType == publisher.VocabImageTaskType {
		queueName = ImageQueueName
	} else if consumerType == publisher.VocabAudioTaskType {
		queueName = AudioQueueName
	}

	consumer := asynqueue.CreateConsumer(queueName, 10, conf)

	handler := &handler.Handler{Config: conf, DictionaryService: dictionaryService}

	mux := asynq.NewServeMux()
	switch consumerType {
	case publisher.VocabImageTaskType:
		mux.HandleFunc(publisher.VocabImageTaskType, handler.HandlerVocabImageTask)
	case publisher.VocabAudioTaskType:
		mux.HandleFunc(publisher.VocabAudioTaskType, handler.HandlerVocabAudioTask)
	default:
		panic(fmt.Errorf("invalid consumer type!"))
	}

	// mux.HandleFunc()
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := consumer.Run(mux)
			utils.CheckError(err)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			consumer.Stop()
			consumer.Shutdown()
			return nil
		},
	})
}
