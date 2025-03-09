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

var ConsumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Consume Asynq tasks",
	Args:  cobra.ExactArgs(2), // Yêu cầu đúng 1 argument (loại task)
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.MustLoad()

		consumerArg := args[0]
		queueNameArg := args[1]

		imageArg := conf.Asynq.IMAGE_CONSUMER_ARGUMENT
		audioArg := conf.Asynq.AUDIO_CONSUMER_ARGUMENT

		if consumerArg != imageArg && consumerArg != audioArg {
			fmt.Printf("Invalid consumer type. Use '%v' or '%v'", imageArg, audioArg)
			return
		}

		internal.Invoke(func(lc fx.Lifecycle, conf *config.Config, dictionaryService service.IDictionaryService) {
			consumerCmd(lc, conf, dictionaryService, consumerArg, queueNameArg)
		}).Run()
	},
}

func consumerCmd(
	lc fx.Lifecycle,
	conf *config.Config,
	dictionaryService service.IDictionaryService,
	consumerArg,
	queueName string,
) {
	// imageArg := conf.Asynq.IMAGE_CONSUMER_ARGUMENT
	// audioArg := conf.Asynq.AUDIO_CONSUMER_ARGUMENT

	// INIT CONSUMER
	consumer := asynqueue.CreateConsumer(queueName, 10, conf)

	handler := &handler.Handler{Config: conf, DictionaryService: dictionaryService}

	// Mapping
	mux := asynq.NewServeMux()
	if consumerArg == conf.Asynq.IMAGE_CONSUMER_ARGUMENT {
		mux.HandleFunc(publisher.VocabImageTaskType, handler.HandlerVocabImageTask)
	} else {
		mux.HandleFunc(publisher.VocabAudioTaskType, handler.HandlerVocabAudioTask)
	}

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
