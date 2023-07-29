package main

import (
	"fmt"
	"github.com/Go-routine-4995/routermgt/adapter/controllers"
	"github.com/Go-routine-4995/routermgt/adapter/repository/simdb"
	"github.com/Go-routine-4995/routermgt/logging"
	"github.com/Go-routine-4995/routermgt/service"
	"gopkg.in/yaml.v2"
	"os"
)

const (
	config  = "conf.yml"
	version = 0.01
)

type Config struct {
	Service struct {
		Nats    string `yaml:"nats"`
		Subject string `yaml:"subject"`
	} `yaml:"service"`
	Database struct {
		PubKey string `yaml:"pubKey"`
	} `yaml:"database"`
	Server struct {
		Url string `yaml:"nats"`
	} `yaml:"server"`
}

func main() {
	var (
		conf string
	)
	fmt.Println("Starting OSS Routers/service v", version)
	args := os.Args
	if len(args) < 2 {
		conf = config
	} else {
		conf = args[1]
	}

	cfg := openFile(conf)

	// new repo
	r := simdb.NewSimDB()

	// new service
	svc := service.NewService(r)

	// new logger
	svc = logging.NewLoggingService(svc)

	// new Api
	api := controllers.NewApiService(svc, cfg.Service.Nats, cfg.Service.Subject)
	api.Start()

}

func openFile(s string) Config {
	f, err := os.Open(s)
	if err != nil {
		processError(err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		processError(err)
	}

	return cfg
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}
