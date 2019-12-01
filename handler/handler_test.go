package handler_test

import (
	"github.com/eugeneverywhere/billing/config"
	"github.com/eugeneverywhere/billing/handler"
	"github.com/eugeneverywhere/billing/rabbit"
	"github.com/eugeneverywhere/billing/types"
	"github.com/lillilli/logger"
	"os"
	"testing"
)

func TestAccountAlreadyExists(t *testing.T) {
	sender := initSender()
	_ = sender.Send(&types.CreateAccount{
		Operation:         &types.Operation{Code: handler.OpCreateAccount},
		ExternalAccountID: "oJnNPGsiuzytMOJPatwt",
	})
}

func TestAccountCreate(t *testing.T) {
	sender := initSender()
	_ = sender.Send(&types.CreateAccount{
		Operation:         &types.Operation{Code: handler.OpCreateAccount},
		ExternalAccountID: "",
	})
}

func TestUnknownOperation(t *testing.T) {
	sender := initSender()
	_ = sender.Send(&types.CreateAccount{
		Operation:         &types.Operation{Code: -1},
		ExternalAccountID: "",
	})
}

func initSender() rabbit.Sender {
	cfg := &config.Config{
		Rabbit: config.RabbitConfig{
			Addr:          "amqp://user:user@localhost:5672/",
			OutputChannel: "output_channel",
			ErrorChannel:  "error_channel",
			InputChannel:  "input_channel",
		},
		Log: logger.Params{
			MinLevel: "DEBUG",
		},
	}

	logger.Init(cfg.Log)
	log := logger.NewLogger("test")

	rabbitConnection, err := rabbit.NewConnection(cfg.Rabbit.Addr)
	if err != nil {
		log.Errorf("Rabbit connecting failed: %v", err)
		os.Exit(1)
	}

	outputChannel, err := rabbitConnection.DeclareQueue(cfg.Rabbit.InputChannel)
	if err != nil {
		log.Errorf("Rabbit connecting to output queue failed: %v", err)
		os.Exit(1)
	}
	return rabbit.NewSender(rabbitConnection, outputChannel)

}
