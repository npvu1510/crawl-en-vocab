package cmd

import (
	"context"
	"fmt"

	"github.com/npvu1510/crawl-en-vocab/internal"
	"github.com/npvu1510/crawl-en-vocab/pkg/config"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

const VocabTask = "vocab"
const TypeVocabGBTask = "vocab:gb:chatgpt"
const TypeVocabUSTask = "vocab:us:chatgpt"
const TypeVocabMeanTask = "vocab:mean:chatgpt"

var CrawlEfcCmd = &cobra.Command{
	Use:   "crawl-efc",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		return internal.Invoke(crawlEfcCmd).Start(context.Background())
	},
}

func crawlEfcCmd(
	lc fx.Lifecycle,
	conf *config.Config,

) error {

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("OnStart crawlEfcCmd")

			fmt.Printf("%+v", conf)
			return nil
		},
	})

	return nil
}
