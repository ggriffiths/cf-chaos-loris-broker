package config

type Config struct {
	ServiceBroker ServiceBroker
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
