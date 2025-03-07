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
	RunE: func(cmd *cobra.Command, args []string) error {
		return internal.Invoke(clearTasksCmd).Start(context.Background())
	},
}

func clearTasksCmd(
	conf *config.Config,
) {
	inspector := asynqueue.GetInspector()

	//
	tasks, err := inspector.ListScheduledTasks("default", 1000)
	if err != nil {
		log.Fatalf("Error while getting tasks: %v", err)
	}

	//
	for _, task := range tasks {
		err := inspector.DeleteTask("default", task.ID)
		if err != nil {
			fmt.Printf("❌ Error while deleting task %s: %v\n", task.ID, err)
		} else {
			fmt.Printf("✅ Deleted task %s\n", task.ID)
		}
	}
	fmt.Println("✅ Clear tasks successfully")
}
