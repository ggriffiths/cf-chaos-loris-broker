package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Altoros/cf-chaos-loris-broker/broker"
	"github.com/Altoros/cf-chaos-loris-broker/cmd"
	"github.com/Altoros/cf-chaos-loris-broker/config"
	"github.com/Altoros/cf-chaos-loris-broker/db"
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

	opts := cmd.CommandOpts{}
	_, err := goflags.ParseArgs(&opts, os.Args[1:])

	brokerLogger.Info("Using config file: " + opts.ConfigPath)

	config, err := config.LoadFromFile(opts.ConfigPath)
	if err != nil {
		brokerLogger.Error("Failed to load the config file", err, lager.Data{
			"broker-config-path": opts.ConfigPath,
		})
		return
	}

	db, err = db.New()
	if err != nil {
		brokerLogger.Error("Failed to connect to the mysql: %s", err)
	}
	defer db.Close()

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
		"port": opts.Port,
	})
	err = http.ListenAndServe(fmt.Sprintf(":%d", opts.Port), nil)
	if err != nil {
		brokerLogger.Error("Failed to start the server", err)
	}
}
