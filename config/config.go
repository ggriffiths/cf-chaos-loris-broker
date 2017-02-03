package config

import (
	"errors"
	"github.com/cloudfoundry-incubator/candiedyaml"
	"os"
)

type Config struct {
	Plans         []Plan
	ServiceBroker ServiceBroker
}

type Plan struct {
	Name        string
	Schedule    string
	Description string
	Probability float64
}

type ServiceBroker struct {
	ServiceID   string
	Name        string
	Description string
	Auth        Auth
}

type Auth struct {
	Username string
	Password string
}

func LoadFromFile(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := candiedyaml.NewDecoder(file).Decode(&config); err != nil {
		return Config{}, err
	}
	// TODO: add validations here
	return config, nil
}

func (c Config) PlanByName(name string) (Plan, error) {
	for _, plan := range c.Plans {
		if plan.Name == name {
			return plan, nil
		}
	}
	return Plan{}, errors.New("Can't find plan")
}
