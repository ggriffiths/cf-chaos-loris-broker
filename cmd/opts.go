package cmd

import (
	boshcmd "github.com/cloudfoundry/bosh-cli/cmd"
)

type CommandOpts struct {
	// -----> Global options
	Version        func() error `long:"version" short:"v" description:"Show Service Brokere version"`
	Deployment     string       `long:"name" short:"n" description:"Environment or Foundation name" env:"FOUNDATION_NAME"`
	CfApiUrl       string       `long:"cf-api" description:"Vault Address" env:"CF_API_URL"`
	CfUsername     string       `long:"cf-username" description:"Vault Address" env:"CF_USERNAME"`
	CfPassword     string       `long:"cf-password" description:"Vault Address" env:"CF_PASSWORD"`
	ChaosLorisHost string       `long:"fetch" short:"h" description:"Fetch variables"`
	Port           int          `long:"port" description:"Port of Service Broker" default:"8080"`
	ConfigPath     string       `long:"config" short:"c" description:"Fetch variables"`
	Help           bool         `long:"help" short:"h" description:"Show this help message"`
}
