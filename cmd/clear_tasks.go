package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/npvu1510/crawl-en-vocab/internal"
	"github.com/npvu1510/crawl-en-vocab/pkg/asynqueue"
	"github.com/npvu1510/crawl-en-vocab/pkg/config"
	"github.com/spf13/cobra"
)

var ClearTasksCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all tasks",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return internal.Invoke(func(config *config.Config) {

			if args[0] != config.Asynq.IMAGE_QUEUE_NAME && args[0] != config.Asynq.AUDIO_QUEUE_NAME {
				fmt.Println("Invalid queue name")
				return
			}

			clearTasksCmd(args[0])
		}).Start(context.Background())
	},
}

func clearTasksCmd(queueName string) {
	var isEmpty bool = true

	inspector := asynqueue.GetInspector()

	for {
		//
		tasks, err := inspector.ListScheduledTasks(queueName)
		if err != nil {
			log.Fatalf("Error while getting tasks: %v", err)
		}

		if len(tasks) == 0 {
			break
		}

		//
		for _, task := range tasks {
			isEmpty = false

			err := inspector.DeleteTask(queueName, task.ID)
			if err != nil {
				fmt.Printf("❌ Error while deleting task %s: %v\n", task.ID, err)
			} else {
				fmt.Printf("✅ Deleted task %s\n", task.ID)
			}
		}
	}

	if !isEmpty {
		fmt.Println("✅ Clear tasks successfully")
	} else {
		fmt.Println("✅ Task queue is empty")
	}
}
