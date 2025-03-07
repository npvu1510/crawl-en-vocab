package asynqueue

import (
	"sync"
	"time"

	"github.com/hibiken/asynq"
	"github.com/npvu1510/crawl-en-vocab/pkg/config"
)

var (
	server         *asynq.Server
	onceInitServer sync.Once
)

func init() {
	conf := config.MustLoad()
	onceInitServer.Do(func() {
		server = asynq.NewServer(
			asynq.RedisClientOpt{
				Addr: conf.Redis.Address,
			},
			asynq.Config{
				Concurrency: 5,
				RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
					return 10 * time.Second
				},
			},
		)
	})
}

func GetServer() *asynq.Server {
	return server
}
