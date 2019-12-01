package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/eugeneverywhere/billing/cache"
	"github.com/eugeneverywhere/billing/config"
	"github.com/eugeneverywhere/billing/db"
	"github.com/eugeneverywhere/billing/dispatcher"
	"github.com/eugeneverywhere/billing/handler"
	"github.com/eugeneverywhere/billing/rabbit"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lillilli/logger"

	"github.com/lillilli/vconf"
)

var (
	configFile = flag.String("config", "", "set service config file")
)

const accountsCacheUpdateInterval = 1 * time.Second

func main() {
	flag.Parse()

	cfg := &config.Config{}

	if err := vconf.InitFromFile(*configFile, cfg); err != nil {
		fmt.Printf("unable to load config: %s\n", err)
		os.Exit(1)
	}

	logger.Init(cfg.Log)
	log := logger.NewLogger("synchronizer")

	database := db.New(db.GenerateMySQLDatabaseURL(cfg.DB), cfg.DB.MaxOpenConnections)

	if err := database.Connect(); err != nil {
		log.Errorf("DB connecting failed: %v", err)
		os.Exit(1)
	}

	defer database.Close()

	rabbitConnection, err := rabbit.NewConnection(cfg.Rabbit.Addr)
	if err != nil {
		log.Errorf("Rabbit connecting failed: %v", err)
		os.Exit(1)
	}

	inputChannel, err := rabbitConnection.DeclareQueue(cfg.Rabbit.InputChannel)
	if err != nil {
		log.Errorf("Rabbit connecting to input queue failed: %v", err)
		os.Exit(1)

	}

	inputSubscriber, err := rabbit.NewQueueSubscriber(rabbitConnection, inputChannel)
	if err != nil {
		log.Errorf("Rabbit subscribe to tickers queue failed: %v", err)
		os.Exit(1)
	}

	outputChannel, err := rabbitConnection.DeclareQueue(cfg.Rabbit.OutputChannel)
	if err != nil {
		log.Errorf("Rabbit connecting to output queue failed: %v", err)
		os.Exit(1)
	}

	sender := rabbit.NewSender(rabbitConnection, outputChannel)

	errorChannel, err := rabbitConnection.DeclareQueue(cfg.Rabbit.ErrorChannel)
	if err != nil {
		log.Errorf("Rabbit connecting to output queue failed: %v", err)
		os.Exit(1)
	}

	errSender := rabbit.NewSender(rabbitConnection, errorChannel)

	accountsCache := cache.NewAccountsCache(database)
	go startUpdateCache(accountsCache)

	log.Info("Start listening...")
	startProcessing(log, database, sender, errSender, inputSubscriber, accountsCache)
}

func startProcessing(log logger.Logger,
	db db.DB,
	sender rabbit.Sender,
	errSender rabbit.Sender,
	subscriber rabbit.QueueSubscriber,
	accountsCache cache.AccountsCache) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	inputChannel, err := subscriber.SubscribeOnQueue()
	if err != nil {
		log.Errorf("Failed to subscribe on input queue: %v", err)
		os.Exit(1)
	}

	handler := handler.NewHandler(db, accountsCache)
	dispatcher := dispatcher.NewDispatcher(handler, sender, errSender)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("Listening stopped")
				return
			case data := <-inputChannel:
				go dispatcher.Dispatch(data.Body)
			}
		}
	}()

	<-signals
	close(signals)
	cancel()
}

func startUpdateCache(accountsCache cache.AccountsCache) {
	accountsTicker := time.NewTicker(accountsCacheUpdateInterval)
	for {
		select {
		case <-accountsTicker.C:
			accountsCache.Update()
		}
	}
}
