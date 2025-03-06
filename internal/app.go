package internal

import (
	"fmt"

	"github.com/npvu1510/crawl-en-vocab/internal/repository"
	"github.com/npvu1510/crawl-en-vocab/internal/service"
	"github.com/npvu1510/crawl-en-vocab/pkg/config"
	"go.uber.org/fx"
)

func Invoke(invokers ...interface{}) *fx.App {
	conf := config.MustLoad()

	fmt.Printf("%+v", conf)

	app := fx.New(
		fx.Provide(
			// Database
			newDatabaseConnection,

			// Repositories
			repository.NewCategoryRepository,
			repository.NewDictionaryRepository,

			// Services
			service.NewCategoryService,
			service.NewDictionaryService),

		fx.Supply(conf),
		fx.Invoke(invokers...),
		// fx.Invoke(func(db *gorm.DB) {}),
	)

	return app
}
