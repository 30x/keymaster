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
	ConfigPollWait       = "apid_poll_wait"
	ConfigCacheDirectory = "cache_directory"
)

func main() {

	v := viper.New()
	v.SetEnvPrefix("goz") // eg. env var "GOZ_APID_URI" will bind to config "apid_uri"
	v.AutomaticEnv()
	v.SetDefault(ConfigApidURI, "http://localhost:8181")
	v.SetDefault(ConfigPollWait, "5")
	v.SetDefault(ConfigCacheDirectory, "/tmp/apidBundleCache")

	apidURI := v.GetString(ConfigApidURI)
	cacheDir := v.GetString(ConfigCacheDirectory)
	timeout := v.GetInt(ConfigPollWait)

	cache, err := client.CreateBundleCache(apidURI, cacheDir, timeout)

	if err != nil {
		log.Fatalf("Could not create cache.  Error is %s", err)
	}

	//loop forever writing configs
	for {

		log.Printf("Attempting to load bundles from cache")

		deployment, err := cache.GetBundles()

		if err != nil {
			log.Printf("Error occured getting bundle from cache.  Error is :%s", err)
			time.Sleep(time.Second * time.Duration(timeout))
			continue
		}

		writeConfig(deployment)

		time.Sleep(time.Second * time.Duration(timeout))
	}
}

func writeConfig(deployment *client.Deployment) {

}
