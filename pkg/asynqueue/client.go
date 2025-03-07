package asynqueue

import (
	"sync"

	"github.com/hibiken/asynq"
	"github.com/npvu1510/crawl-en-vocab/pkg/config"
)

var (
	client *asynq.Client
	once   sync.Once
	conf   *config.Config
)

func init() {
	conf = config.MustLoad()
	once.Do(func() {
		client = asynq.NewClient(
			asynq.RedisClientOpt{
				Addr: conf.Redis.Address,
				// Password: conf.Redis.Password,
				// DB:       conf.Asynq.DB,
			},
		)
	})
}

func Close() {
	if client != nil {
		client.Close()
	}
}

func GetClient() *asynq.Client {
	return client
}
