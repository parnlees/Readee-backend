package config

import (
	cc "Readee-Backend/common"
	"Readee-Backend/type/common"
	"log"
	"os"
	"time"

	"github.com/patrickmn/go-cache"
	"gopkg.in/yaml.v3"
)

var AppCache *cache.Cache

func Init() {
	config := new(common.Config)

	yml, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(yml, config); err != nil {
		log.Fatal("Unable to parse configuration file", err)
	}

	AppCache = cache.New(5*time.Minute, 10*time.Minute)

	cc.Config = config
}
