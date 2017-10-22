package config

import (
	"sync"
)

// Config is a signleton which should be accessed/initialized only by the GetInstance function, it houses the program configuration
type Config struct {
	ClientConfig *ClientConfig
}

var instance *Config
var once sync.Once

// GetInstance returns a singleton instance of the program configuration
func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{ClientConfig: &ClientConfig{}}
	})
	return instance
}

// ClientConfig houses the configurations for the client portion of the program
type ClientConfig struct {
	MessagePort uint16
}

// ServerConfig houses the configurations for the client portion of the program
type ServerConfig struct {
	Address string
	Port    uint16
}
