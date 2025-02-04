package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Config struct {
	Proxy struct {
		Ip   string `yaml:"ip"`   //Proxy server IP
		Port string `yaml:"port"` //Port on which the application is running
		Code string `yaml:"code"` //Code for access as a request handler (If it is empty, you will connect to the proxy as a regular user and will not be able to receive requests)
	} `yaml:"proxy"`

	LocalStorage struct {
		Directory string `yaml:"directory"`    //The directory where the versions of the executable files will be located
		Name      string `yaml:"filename_exe"` //Name of the executable file
	} `yaml:"local_storage"`
}

// configPath stores the path to the configuration file
var configPath = "config.yaml"

func SetConfigPath(path string) {
	configPath = path
}

// SetConfigPath allows you to change the path to the configuration file
func LoadConfig() *Config {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	return &cfg
}
