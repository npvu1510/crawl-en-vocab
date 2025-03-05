package internal

import (
	"fmt"

	"github.com/npvu1510/crawl-en-vocab/pkg/config"
	"go.uber.org/fx"
)

func Invoke(invokers ...interface{}) *fx.App {
	conf := config.MustLoad()

	fmt.Printf("%+v", conf)

	app := fx.New(
		fx.Provide(newDatabaseConnection),

		fx.Supply(conf),
		fx.Invoke(invokers...),
	)

	return app
}
