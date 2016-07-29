package main

import (
	"time"

	"log"

	"github.com/30x/keymaster/client"
	"github.com/spf13/viper"
)

const (
	//ConfigApidURI defualt config value for the apid location
	ConfigApidURI = "apid_uri"
	//ConfigPollWait the number of seconds to wait after successfully polling apid before polling again
	ConfigPollWait = "apid_poll_wait"
)

func main() {

	v := viper.New()
	v.SetEnvPrefix("goz") // eg. env var "GOZ_APID_URI" will bind to config "apid_uri"
	v.AutomaticEnv()
	v.SetDefault(ConfigApidURI, "http://localhost:8181")
	v.SetDefault(ConfigPollWait, "5")

	apidURI := v.GetString(ConfigApidURI)
	timeout := v.GetInt(ConfigPollWait)

	cache := &keymaster.BundleCache{
		ApidURI: apidURI,
	}

	//loop forever writing configs
	for {

		log.Printf("Attempting to load bundles from cache")

		bundles, err := cache.GetBundles()

		if err != nil {
			log.Printf("Error occured getting bundle from cache.  Error is :%s", err)
			time.Sleep(time.Second * time.Duration(timeout))
			continue
		}

		if len(bundles) > 0 {
			writeConfig(bundles)
		}

		time.Sleep(time.Second * time.Duration(timeout))
	}
}

func writeConfig(bundles []keymaster.Bundle) {

}
