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
	Use:   "vocab-audio-pub",
	Short: "",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		queueNameArg := args[0]

		internal.Invoke(func(lc fx.Lifecycle, conf *config.Config, dictionaryService service.IDictionaryService) {
			vocabAudioPublisherCmd(lc, conf, dictionaryService, queueNameArg)
		}).Run()
	},
}

func vocabAudioPublisherCmd(lc fx.Lifecycle, conf *config.Config, dictionaryService service.IDictionaryService, queueName string) {
	fmt.Println("✅ Cron job vocab-audio initializing...")
	c := cron.New(cron.WithSeconds())

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("✅ Cron job vocab-audio started")
			_, err := c.AddFunc("*/15 * * * * *", func() {
				err := scanVocabAudioDb(conf, dictionaryService, queueName)
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

func scanVocabAudioDb(conf *config.Config, dictionaryService service.IDictionaryService, queueName string) error {
	limitRecord := conf.EmojiFlashcard.DITCTIONARY_PUBLISH_LIMIT_RECORD

	dictionaries, err := dictionaryService.GetDictionaries("", false, true, 1, limitRecord)
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

	// batchSize := conf.Emoji_Flashcard.DITCTIONARY_PUBLISH_BATCH_SIZE
	for idx, t := range newTasks {
		recordId := dictionaries[idx].Id
		taskId := fmt.Sprintf("%v%v", VocabAudioTaskName, recordId)

		// batchIdx := idx/batchSize + 1
		// delayMinute := time.Duration(batchIdx)

		_, err := publisher.Enqueue(t,
			asynq.TaskID(taskId),
			asynq.Queue(queueName),
			// asynq.ProcessIn(delayMinute*time.Minute),
			asynq.ProcessIn(10*time.Second),
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
