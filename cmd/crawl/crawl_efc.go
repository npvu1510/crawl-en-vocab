package crawl

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/npvu1510/crawl-en-vocab/internal"
	"github.com/npvu1510/crawl-en-vocab/internal/model"
	"github.com/npvu1510/crawl-en-vocab/internal/service"
	"github.com/npvu1510/crawl-en-vocab/pkg/config"
	"github.com/npvu1510/crawl-en-vocab/pkg/utils"
	"gorm.io/gorm"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

const VocabTask = "vocab"
const TypeVocabGBTask = "vocab:gb:chatgpt"
const TypeVocabUSTask = "vocab:us:chatgpt"
const TypeVocabMeanTask = "vocab:mean:chatgpt"

var CrawlEfcCmd = &cobra.Command{
	Use:   "crawl-efc",
	Short: "Crawl English Flashcards",
	RunE: func(cmd *cobra.Command, args []string) error {
		return internal.Invoke(crawlEfcCmd).Start(context.Background())
	},
}

func crawlEfcCmd(
	lc fx.Lifecycle,
	db *gorm.DB,
	conf *config.Config,
	categoryService service.ICategoryService,
	dictionaryService service.IDictionaryService,
) error {

	// Implementation
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Implementation
			startCrawling(conf, db, categoryService, dictionaryService)
			return nil
		},
	})

	return nil
}

func getPageDocument(url string) *goquery.Document {
	res, err := http.Get(url)
	utils.CheckError(err)

	defer res.Body.Close()

	document, err := goquery.NewDocumentFromReader(res.Body)
	utils.CheckError(err)

	return document
}

// ##############################CATEGORY##############################
func element2Category(opElement *goquery.Selection) *model.Category {
	categoryName := strings.TrimSpace(opElement.Text())
	categoryIdStr, isExist := opElement.Attr("value")

	if !isExist {
		return nil
	}

	if categoryIdStr == "" {
		return nil
	}

	categoryId, err := strconv.Atoi(categoryIdStr)
	utils.CheckError(err)

	return model.NewCategory(categoryId, categoryName)
}

func getCategories(document *goquery.Document, workerNum int, categoryService service.ICategoryService) {
	// Channels
	optionsCh := make(chan *goquery.Selection)

	// Worker pool
	var wg sync.WaitGroup
	for workerId := 1; workerId <= int(workerNum); workerId++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()

			for opElement := range optionsCh {
				// Crawling
				fmt.Printf("🔃 Worker %v: CATEGORY...\n", workerId)

				category := element2Category(opElement)
				if category == nil {
					continue
				}

				if category.Id == -1 {
					fmt.Printf("❌ Worker %v: CATEGORY '%s' doesn't contain an id.\n", workerId, category.Name)
					continue
				}

				// Save to database
				// category.Save()
				categoryService.CreateCategory(category)

			}
		}(workerId)
	}

	mainDiv := document.Find("#word-list-and-picture-dictionary")
	mainDiv.Find("#collection_id option").Each(func(index int, optionElement *goquery.Selection) {
		optionsCh <- optionElement

	})

	close(optionsCh)
	wg.Wait()

}

func getCategoryMap(db *gorm.DB) map[int]model.Category {
	var categories []model.Category
	db.Find(&categories) // Load tất cả categories

	categoryMap := make(map[int]model.Category)
	for _, category := range categories {
		categoryMap[category.Id] = category
	}

	return categoryMap
}

func element2Dictionary(itemElement *goquery.Selection, categoryMap map[int]model.Category) *model.Dictionary {
	// PARSE IDS STRING
	idsStr, isExist := itemElement.Attr("data-collection-ids")
	if !isExist {
		return nil
	}

	// CREATE DICTIONARIES_CATEGORIES records
	categories := make([]model.Category, 0)
	if idsStr != "" {
		for _, categoryIdStr := range strings.Split(idsStr, ",") {
			categoryId, err := strconv.Atoi(categoryIdStr)
			utils.CheckError(err)

			// Finding
			category, ok := categoryMap[categoryId]

			if !ok {
				continue
			}

			categories = append(categories, category)
		}
	}
	// DEFINITION
	definition := itemElement.Find(".definition").Text()

	dictonary := model.NewDictionary(definition, categories)
	return dictonary
}
func getDictionaries(document *goquery.Document, workerNum int, batchSize int, db *gorm.DB, dictionaryService service.IDictionaryService) {
	// Get categories Ids
	categoryMap := getCategoryMap(db)

	// Channels
	dictionaryCh := make(chan *goquery.Selection)

	// Worker pool
	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Batch
	var dictionaryBatch = make([]*model.Dictionary, 0)

	for workerId := 1; workerId <= workerNum; workerId++ {
		wg.Add(1)
		go func(workerId int) {
			// var assetsWg sync.WaitGroup

			defer wg.Done()

			for element := range dictionaryCh {
				// ################### CRAWLING ###################
				fmt.Printf("🔃 Worker %v: DICTIONARY...\n", workerId)

				dictionary := element2Dictionary(element, categoryMap)
				if dictionary == nil {
					continue
				}

				mutex.Lock()
				dictionaryBatch = append(dictionaryBatch, dictionary)

				if len(dictionaryBatch) >= batchSize {
					tempBatch := make([]*model.Dictionary, len(dictionaryBatch))
					copy(tempBatch, dictionaryBatch)

					dictionaryBatch = make([]*model.Dictionary, 0)
					mutex.Unlock()

					err := dictionaryService.CreateDictionaries(tempBatch, len(tempBatch))
					utils.CheckError(err)

				} else {
					mutex.Unlock()
				}

			}

		}(workerId)

	}

	// Send list of dictionaries to channel
	mainDiv := document.Find("#word-list-and-picture-dictionary")
	mainDiv.Find(`#picture-dictionary-page-container li`).Each(func(i int, dictionaryEle *goquery.Selection) {
		dictionaryCh <- dictionaryEle
	})

	close(dictionaryCh)
	wg.Wait()

	if len(dictionaryBatch) >= 0 {
		err := dictionaryService.CreateDictionaries(dictionaryBatch, len(dictionaryBatch))
		utils.CheckError(err)
	}
}

func startCrawling(conf *config.Config, db *gorm.DB, categoryService service.ICategoryService, dictionaryService service.IDictionaryService) {
	crawlingUrl := conf.EMOJI_FLASHCARD.CRAWLING_URL

	workerNum, dictionary_batch_size := conf.EMOJI_FLASHCARD.WORKER_NUM, conf.EMOJI_FLASHCARD.DITCTIONARY_BATCH_SIZE

	// GET page document
	document := getPageDocument(crawlingUrl)

	// GET Categories
	getCategories(document, workerNum, categoryService)

	// GET Dictionaries
	getDictionaries(document, workerNum, dictionary_batch_size, db, dictionaryService)
}
