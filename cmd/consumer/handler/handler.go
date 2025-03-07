package handler

import (
	"github.com/npvu1510/crawl-en-vocab/internal/service"
	"github.com/npvu1510/crawl-en-vocab/pkg/config"
)

type Handler struct {
	Config            *config.Config
	DictionaryService service.IDictionaryService
}
