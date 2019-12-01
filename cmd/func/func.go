package main

import (
	"flag"
	"fmt"
	"github.com/eugeneverywhere/billing/config"
	"github.com/eugeneverywhere/billing/rabbit"
	"github.com/eugeneverywhere/billing/types"
	"github.com/lillilli/logger"
	"github.com/lillilli/vconf"
	"os"
	"time"
)

var (
	configFile = flag.String("config", "", "set service config file")
)

const acc1 = "FTA1"
const acc2 = "FTA2"

var sender rabbit.Sender

func main() {
	flag.Parse()

	cfg := &config.Config{}

	if err := vconf.InitFromFile(*configFile, cfg); err != nil {
		fmt.Printf("unable to load config: %s\n", err)
		os.Exit(1)
	}

	logger.Init(cfg.Log)
	log := logger.NewLogger("billing service")

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

	sender = rabbit.NewSender(rabbitConnection, inputChannel)

	log.Info("Sending messages...")
	createAccount(acc1)
	createAccount(acc2)
	time.Sleep(2 * time.Second)

	addAmount(acc1, 100.30)
	addAmount(acc2, 200.70)
	addAmount(acc2, -5000)
	for i := 0; i < 100; i++ {
		transferAmount(acc2, acc1, 1000)
		transferAmount(acc1, acc2, 1000)
	}
	log.Info("Done")

}

var counter int

func createAccount(extId string) {
	_ = sender.Send(types.CreateAccount{
		Operation: &types.Operation{
			ConsumerID:  0,
			OperationID: counter,
			Code:        types.OpCreateAccount,
		},
		ExternalAccountID: extId,
	})
	counter++
}

func addAmount(extId string, amount float64) {
	_ = sender.Send(types.AddAmount{
		Operation: &types.Operation{
			ConsumerID:  0,
			OperationID: counter,
			Code:        types.OpAddAmount,
		},
		Amount:            amount,
		ExternalAccountID: extId,
	})
	counter++
}

func transferAmount(src string, tgt string, amount float64) {
	_ = sender.Send(types.TransferAmount{
		Operation: &types.Operation{
			ConsumerID:  0,
			OperationID: counter,
			Code:        types.OpTransfer,
		},
		Amount: amount,
		Source: src,
		Target: tgt,
	})
	counter++
}
