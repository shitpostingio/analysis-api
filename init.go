package main

import (
	"github.com/gorilla/mux"
	"github.com/shitpostingio/analysis-api/api/cache"
	"github.com/shitpostingio/analysis-api/api/database"
	"github.com/shitpostingio/analysis-api/api/handler"
	"github.com/shitpostingio/analysis-api/configuration"
	"github.com/shitpostingio/analysis-commons/downloader"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {

	setEnvVars()
	cfg, err := configuration.Load(cfgPath)
	if err != nil {
		log.Fatal("Unable to load configuration", err)
	}

	if !cfg.Testing {
		database.Connect(&cfg.Database)
	} else {
		log.SetLevel(log.DebugLevel)
	}

	err = cache.NewFingerprintRedisClient(cfg.Redis.Address, "", cfg.Redis.FingerprintDatabase)
	if err != nil {
		log.Fatal("Unable to connect to fingerprint cache", err)
	}

	err = cache.NewNSFWRedisClient(cfg.Redis.Address, "", cfg.Redis.NSFWDatabase)
	if err != nil {
		log.Fatal("Unable to connect to nsfw cache", err)
	}

	telegramDownloader := &downloader.TelegramDownloader{MaxDownloadSize: cfg.MaxDownloadSize}
	telegramHandler = handler.New(cfg.Testing, cfg.NSFWEndpoint, cfg.FingerprintEndpoint, cfg.GibberishEndpoint, telegramDownloader)
	directDownloader := &downloader.DirectDownloader{MaxDownloadSize: cfg.MaxDownloadSize}
	directHandler = handler.New(cfg.Testing, cfg.NSFWEndpoint, cfg.FingerprintEndpoint, cfg.GibberishEndpoint, directDownloader)
	r = mux.NewRouter()

}

func setEnvVars() {

	cp := os.Getenv(cfgPathKey)
	if cp != "" {
		cfgPath = cp
	}

	add := os.Getenv(bindAddressKey)
	if add != "" {
		bindAddress = add
	}

}
