package config

import (
	"fmt"
	"sync"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	App struct {
		LOCAL_IMAGE_PATH string `env:"LOCAL_IMAGE_PATH"`
		LOCAL_AUDIO_PATH string `env:"LOCAL_AUDIO_PATH"`
	}

	Emoji_Flashcard struct {
		CRAWLING_URL           string `env:"EMOJI_FLASHCARD_CRAWLING_URL"`
		DITCTIONARY_BATCH_SIZE int    `env:"EMOJI_FLASHCARD_DITCTIONARY_BATCH_SIZE" envDefault:"18"`
		WORKER_NUM             int    `env:"EMOJI_FLASHCARD_WORKER_NUM" envDefault:"10"`
		SOURCE                 string `env:"EMOJI_FLASHCARD_SRC" envDefault:"EMOJI_FLASHCARD_SRC"`
	}
	Postgres struct {
		Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
		Port     int    `env:"POSTGRES_PORT" envDefault:"5432"`
		User     string `env:"POSTGRES_USER" envDefault:"postgres"`
		Password string `env:"POSTGRES_PASSWORD" envDefault:""`
		DbName   string `env:"POSTGRES_DBNAME" envDefault:""`
	}

	Redis struct {
		Address string `env:"REDIS_ADDRESS" envDefault:"127.0.0.1:6379"`
	}
}

func init() {
	fmt.Println("INIT...")
	if err := godotenv.Load(".env"); err != nil {
		panic(fmt.Sprintf("Error while loading .env file: %v", err))
	}

}

var (
	once   sync.Once
	lock   sync.Mutex
	config *Config
)

func Load() Config {
	fmt.Println("LOAD...")

	var conf Config
	// log := logger.Get()
	once.Do(func() {
		if err := env.Parse(&conf); err != nil {
			// log.Fatalf("Error reading the environment variables: %v", err)
			// return
			panic(fmt.Sprintf("Error reading the environment variables: %v", err))
		}
	})
	return conf
}

func MustLoad() *Config {
	fmt.Println("MUSTLOAD...")
	lock.Lock()
	defer lock.Unlock()
	if config != nil {
		return config
	}
	_conf := Load()

	config = &_conf
	return config
}
