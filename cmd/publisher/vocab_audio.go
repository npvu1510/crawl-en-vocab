package publisher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jinzhu/copier"
	"github.com/npvu1510/crawl-en-vocab/internal"
	"github.com/npvu1510/crawl-en-vocab/internal/model"
	"github.com/npvu1510/crawl-en-vocab/internal/service"
	"github.com/npvu1510/crawl-en-vocab/pkg/asynqueue"
	"github.com/npvu1510/crawl-en-vocab/pkg/config"
	"github.com/npvu1510/crawl-en-vocab/pkg/utils"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

const (
	VocabAudioTaskName = "vocab_audio"
	VocabAudioTaskType = "vocab:audio"
)

var VocabAudioPublisherCmd = &cobra.Command{
	Use:   "cronjob",
	Short: "Chạy cron job định kỳ",
	Run: func(cmd *cobra.Command, args []string) {
		internal.Invoke(vocabAudioPublisherCmd).Run()
	},
}

func vocabAudioPublisherCmd(lc fx.Lifecycle, conf *config.Config, dictionaryService service.IDictionaryService) {
	fmt.Println("✅ Cron job vocab-audio initializing...")
	c := cron.New(cron.WithSeconds())

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("✅ Cron job vocab-audio started")
			_, err := c.AddFunc("*/10 * * * * *", func() {
				err := scanVocabAudioDb(conf, dictionaryService)
				utils.CheckError(err)
			})
			utils.CheckError(err)
			c.Start()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			c.Stop()
			fmt.Println("❌ Cron job vocab-audio stopped")
			return nil
		},
	})
}

func scanVocabAudioDb(conf *config.Config, dictionaryService service.IDictionaryService) error {
	dictionaries, err := dictionaryService.GetDictionaries()
	utils.CheckError(err)

	// CREATE NEW TASKS
	newTasks := make([]*asynq.Task, 0)
	for _, d := range dictionaries {
		cloneDictionary := &model.Dictionary{}
		_ = copier.CopyWithOption(cloneDictionary, d, copier.Option{IgnoreEmpty: true, DeepCopy: true})

		bytesData, err := json.Marshal(cloneDictionary)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %v", err)
		}
		vocabAudioTask := asynq.NewTask(VocabAudioTaskType, bytesData)
		newTasks = append(newTasks, vocabAudioTask)
	}

	// ENQUEUE TASKS
	publisher := asynqueue.GetClient()

	batchSize := conf.Emoji_Flashcard.DITCTIONARY_PUBLISH_BATCH_SIZE
	for idx, t := range newTasks {
		recordId := dictionaries[idx].Id
		taskId := fmt.Sprintf("%v%v", VocabAudioTaskName, recordId)

		// fmt.Println(VocabAudioTaskName)
		// fmt.Println(recordId)

		batchIdx := idx/batchSize + 1
		delayMinute := time.Duration(batchIdx)

		_, err := publisher.Enqueue(t,
			asynq.TaskID(taskId),
			asynq.ProcessIn(delayMinute*time.Minute),
			asynq.MaxRetry(2),
			asynq.Timeout(90*time.Second))

		switch {
		case errors.Is(err, asynq.ErrTaskIDConflict):
			fmt.Printf("vocabAudioPublisherCmd ErrTaskIDConflict (%v): %v\n", taskId, err)
		case err != nil:
			fmt.Printf("vocabAudioPublisherCmd failed: %v\n", err)
		default:
			fmt.Printf("Enqueued task: %s\n", taskId)
		}

	}

	return nil
}
