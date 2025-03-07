package consumer

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/npvu1510/crawl-en-vocab/cmd/publisher"
	"github.com/npvu1510/crawl-en-vocab/internal/service"
	"github.com/npvu1510/crawl-en-vocab/pkg/asynqueue"
	"github.com/npvu1510/crawl-en-vocab/pkg/config"
	"github.com/npvu1510/crawl-en-vocab/pkg/utils"
	"go.uber.org/fx"
)

func consumerCmd(
	lc fx.Lifecycle,
	conf *config.Config,
	dictionaryService service.IDictionaryService,
) {

	consumer := asynqueue.GetServer()

	vocabImageConsumer := &VocabImageConsumer{Service: dictionaryService, Config: conf}

	mux := asynq.NewServeMux()
	mux.HandleFunc(publisher.VocabImageTaskType, vocabImageConsumer.handlerVocabImageTask)
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
