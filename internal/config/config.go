package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName    string          `json:"serviceName"`
	ServiceAddress string          `json:"servicePort"`
	ServiceID      string          `json:"serviceID"`
	RPCAddress     string          `json:"rpcAddress"`
	TrustedService map[string]bool `json:"trustedService"`
	Environment    Environment     `json:"environment"`
	FilePath       string
	RunSince       time.Time

	MaxCore string

	DeployConfigPath string
}

const logTagConfig = "[Init Config]"

var config *Config

func Init() {
	godotenv.Load("conf/.env")

	conf := Config{
		ServiceName:    os.Getenv("SERVICE_NAME"),
		ServiceAddress: os.Getenv("SERVICE_ADDR"),
		ServiceID:      os.Getenv("SERVICE_ID"),
		RPCAddress:     os.Getenv("GPRC_ADDR"),
		FilePath:       os.Getenv("FILE_PATH"),
		MaxCore:        os.Getenv("MAX_CORE"),
	}

	if conf.ServiceName == "" {
		log.Fatalf("%s service name should not be empty", logTagConfig)
	}

	if conf.ServiceAddress == "" {
		log.Fatalf("%s service port should not be empty", logTagConfig)
	}

	envString := os.Getenv("ENVIRONMENT")
	if envString != "dev" && envString != "prod" && envString != "local" {
		log.Fatalf("%s environment must be either local, dev or prod, found: %s", logTagConfig, envString)
	}

	deployPath := os.Getenv("DEPLOY_CONFIG_PATH")
	if deployPath == "" {
		log.Fatalf("%s deploy path should not be empty", logTagConfig)
	}
	conf.Environment = Environment(envString)
	conf.DeployConfigPath = deployPath
	conf.RunSince = time.Now()

	conf.TrustedService = map[string]bool{conf.ServiceID: true}
	if trusted := os.Getenv("TRUSTED_SERVICES"); trusted == "" {
		conf.TrustedService["STELLAR_HENTAI"] = true
	} else {
		for _, svc := range strings.Split(trusted, ",") {
			if _, ok := conf.TrustedService[svc]; !ok {
				conf.TrustedService[svc] = true
			}
		}
	}

	config = &conf
}

func Get() (conf *Config) {
	conf = config
	return
}
