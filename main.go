package main

import (
	"time"

	"log"

	"github.com/30x/keymaster/client"
	"github.com/30x/keymaster/nginx"
	"github.com/spf13/viper"
)

const (
	//ConfigApidURI defualt config value for the apid location
	ConfigApidURI = "apid_uri"
	//ConfigPollWait the number of seconds to wait after successfully polling apid before polling again
	ConfigPollWait       = "apid_poll_wait"
	ConfigCacheDirectory = "cache_directory"

	//ConfigNginxDir the directory that nginx is located in
	ConfigNginxDir = "nginx_dir"
)

func main() {

	v := viper.New()
	v.SetEnvPrefix("goz") // eg. env var "GOZ_APID_URI" will bind to config "apid_uri"
	v.AutomaticEnv()
	v.SetDefault(ConfigApidURI, "http://localhost:8181")
	v.SetDefault(ConfigPollWait, "5")
	v.SetDefault(ConfigCacheDirectory, "/tmp/apidBundleCache")

	//use openresty for now.  Must have LUAJIT installed
	v.SetDefault(ConfigNginxDir, " /usr/local/Cellar/openresty/1.9.15.1/")

	apidURI := v.GetString(ConfigApidURI)
	cacheDir := v.GetString(ConfigCacheDirectory)
	timeout := v.GetInt(ConfigPollWait)
	nginxDir := v.GetString(ConfigNginxDir)

	cache, err := client.CreateBundleCache(apidURI, cacheDir, timeout)

	if err != nil {
		log.Fatalf("Could not create cache.  Error is %s", err)
	}

	manager := nginx.NewManager(cache, nginxDir)

	//loop forever writing configs
	for {

		log.Printf("Runnig manager")

		err := manager.ApplyDeployment()

		if err != nil {
			log.Printf("An error occured when attempting to apply the latest deployment.  Error is %s", err)
		}

		time.Sleep(time.Second * time.Duration(timeout))
	}
}
