package cmd

type CommandOpts struct {
	// -----> Global options
	Version            func() error `long:"version" short:"v" description:"Show Service Brokere version"`
	Deployment         string       `long:"name" short:"n" description:"Environment or Foundation name" env:"FOUNDATION_NAME"`
	CfApiUrl           string       `long:"cf-api" description:"Vault Address" env:"CF_API_URL"`
	CfUsername         string       `long:"cf-username" description:"Vault Address" env:"CF_USERNAME"`
	CfPassword         string       `long:"cf-password" description:"Vault Address" env:"CF_PASSWORD"`
	ChaosLorisHost     string       `long:"chaos-loris-host" short:"f" description:"Fetch variables"`
	Port               int          `long:"port" description:"Port of Service Broker" default:"8080"`
	ConfigPath         string       `long:"config" short:"c" description:"Path to config with plans"`
	ServiceID          string       `long:"service-broker-id" descrition:"service broker service id" default:"chaos-loris-broker" env:"SERVICE_BROKER_ID"`
	Name               string       `long:"service-broker-name" descrition:"service broker name" default:"chaos-loris-broker" env:"SERVICE_BROKER_NAME"`
	Description        string       `long:"service-broker-description" descrition:"service broker name" default:"Service for running destructive tests on a app" env:"SERVICE_BROKER_DESCRIPTION"`
	ChaosLorisUsername string       `long:"service-broker-username" description:"broker-username" default:"loris"`
	ChaosLorisPassword string       `long:"service-broker-password" description:"broker-password" default:"cha0s-l0r1s"`
	Help               bool         `long:"help" short:"h" description:"Show this help message"`
}
