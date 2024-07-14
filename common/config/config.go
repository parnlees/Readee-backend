package config

import (
	cc "Readee-Backend/common"
	"Readee-Backend/type/common"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func Init() {
	config := new(common.Config)

	yml, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(yml, config); err != nil {
		log.Fatal("Unable to parse configuration file", err)
	}

	cc.Config = config
}
