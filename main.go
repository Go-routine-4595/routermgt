package main

import (
	"fmt"
	"github.com/Go-routine-4995/routermgt/adapter/controllers"
	"github.com/Go-routine-4995/routermgt/adapter/repository/postgres"
	"github.com/Go-routine-4995/routermgt/logging"
	"github.com/Go-routine-4995/routermgt/service"
	"gopkg.in/yaml.v2"
	"os"
	"sync"
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
		PubKey     string `yaml:"pubKey"`
		ClientCert string `yaml:"client-cert"`
		ClientKey  string `yaml:"client-key"`
		ServerCert string `yaml:"server-cert"`
		Address    string `yaml:"address"`
		User       string `yaml:"user"`
		Password   string `yaml:"password"`
		Database   string `yaml:"database"`
	} `yaml:"database"`
	Server struct {
		Url string `yaml:"nats"`
	} `yaml:"server"`
}

func main() {
	var (
		conf string
		wg   *sync.WaitGroup
	)
	fmt.Println("Starting OSS Routers/service v", version)
	args := os.Args
	if len(args) < 2 {
		conf = config
	} else {
		conf = args[1]
	}

	wg = new(sync.WaitGroup)
	cfg := openFile(conf)

	// new repo
	//r := simdb.NewSimDB()
	wg.Add(1)
	r := postgres.NewPostgres(cfg.Database.Address,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Database,
		cfg.Database.ClientCert,
		cfg.Database.ClientKey,
		cfg.Database.ServerCert,
		wg)

	// new service
	svc := service.NewService(r)

	// new logger
	svc = logging.NewLoggingService(svc)

	// new Api
	api := controllers.NewApiService(svc, cfg.Service.Nats, cfg.Service.Subject, wg)
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
