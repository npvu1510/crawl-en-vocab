package asynqueue

import (
	"sync"

	"github.com/hibiken/asynq"
	"github.com/npvu1510/crawl-en-vocab/pkg/config"
)

var (
	inspector         *asynq.Inspector
	onceInitInspector sync.Once
)

func init() {
	conf := config.MustLoad()

	onceInitInspector.Do(func() {
		inspector = asynq.NewInspector(asynq.RedisClientOpt{
			Addr: conf.Redis.Address,
		})
	})

}

func GetInspector() *asynq.Inspector {
	return inspector
}

func CloseInspector() {
	if inspector != nil {
		inspector.Close()
	}
}
