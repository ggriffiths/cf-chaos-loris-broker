package main

import (
	"fmt"
	"net/http"
	"os"

	// goflags "github.com/jessevdk/go-flags"
	"github.com/Altoros/cf-chaos-loris-broker/broker"
	"github.com/Altoros/cf-chaos-loris-broker/config"
	"github.com/pivotal-cf/brokerapi"
	"github.com/pivotal-golang/lager"
)

var (
	localPersisterPath string
	brokerStateRoot    string
	brokerConfigPath   string
)

func main() {
	brokerLogger := lager.NewLogger("service-broker")
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

	if brokerConfigPath == "" {
		brokerLogger.Error("No config file specified", nil)
		return
	}

	brokerLogger.Info("Using config file: " + brokerConfigPath)

	// conf, err := config.LoadFromFile(brokerConfigPath)
	// if err != nil {
	// 	brokerLogger.Error("Failed to load the config file", err, lager.Data{
	// 		"broker-config-path": brokerConfigPath,
	// 	})
	// 	return
	// }

	config := config.Config{}

	serviceBroker := broker.NewServiceBroker(
		// instancecreators.NewDefault(config, brokerLogger),
		// instancebinders.NewDefault(config, brokerLogger),
		// persisters.NewLocalPersister(localPersisterPath),
		config,
		brokerLogger,
	)

	credentials := brokerapi.BrokerCredentials{
		Username: config.ServiceBroker.Auth.Username,
		Password: config.ServiceBroker.Auth.Password,
	}

	brokerAPI := brokerapi.New(serviceBroker, brokerLogger, credentials)
	http.Handle("/", brokerAPI)
	brokerLogger.Info("Listening for requests", lager.Data{
		"port": config.ServiceBroker.Port,
	})
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.ServiceBroker.Port), nil)
	if err != nil {
		brokerLogger.Error("Failed to start the server", err)
	}
}
