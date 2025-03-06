package publisher

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/npvu1510/crawl-en-vocab/internal"
	"github.com/npvu1510/crawl-en-vocab/internal/service"
	"github.com/npvu1510/crawl-en-vocab/pkg/config"
	"github.com/npvu1510/crawl-en-vocab/pkg/uploader"
	"github.com/npvu1510/crawl-en-vocab/pkg/utils"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

// var VocabImagePublisherCmd = &cobra.Command{
// 	Use:   "cronjob",
// 	Short: "Chạy cron job định kỳ",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		c := cron.New(cron.WithSeconds())              // Cho phép biểu thức cron hỗ trợ giây
// 		_, err := c.AddFunc("*/10 * * * * *", func() { // Chạy mỗi 10 giây
// 			fmt.Println("Hello from cron job at", time.Now().Format("15:04:05"))
// 		})
// 		if err != nil {
// 			fmt.Println("Lỗi khi tạo cronjob:", err)
// 			return
// 		}

// 		c.Start()
// 		fmt.Println("Cron job started. Press Ctrl+C to exit.")

// 		// Chặn main goroutine để không kết thúc chương trình
// 		select {}
// 	},
// }

var VocabImagePublisherCmd = &cobra.Command{
	Use:   "cronjob",
	Short: "Chạy cron job định kỳ",
	Run: func(cmd *cobra.Command, args []string) {
		internal.Invoke(vocabImagePublisherCmd).Run()
	},
}

func vocabImagePublisherCmd(lc fx.Lifecycle, conf *config.Config, dictionaryService service.IDictionaryService) {
	fmt.Println("✅ Cron job vocab-image initializing...")
	c := cron.New(cron.WithSeconds())

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("✅ Cron job vocab-image started")
			_, err := c.AddFunc("*/3 * * * * *", func() {
				scanVocabImageDb(conf, dictionaryService)
			})
			utils.CheckError(err)
			c.Start()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			c.Stop()
			fmt.Println("❌ Cron job vocab-image stopped")
			return nil
		},
	})
}

func scanVocabImageDb(conf *config.Config, dictionaryService service.IDictionaryService) {
	fmt.Println("scanVocabImageDb")
	dictionaries, err := dictionaryService.GetDictionaries()
	utils.CheckError(err)

	gcsService := uploader.Service

	var wg sync.WaitGroup

	for _, d := range dictionaries {
		filename := utils.Definition2FileName(d.Definition) + ".svg"
		fullFilePath := filepath.Join(conf.App.LOCAL_IMAGE_PATH, filename)

		wg.Add(1)
		go func() {
			var err error

			defer wg.Done()
			imagePath, err := utils.DownloadImage(d.Image, "./"+fullFilePath)
			utils.CheckError(err)

			err = gcsService.UploadFile(imagePath, "/vu1/images/"+filename)
			fmt.Println(err)

		}()
	}

	wg.Wait()
}
