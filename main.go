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
	ConfigPollWait = "apid_poll_wait"

	//ConfigNginxDir the directory that nginx is located in
	ConfigNginxDir = "nginx_dir"

	//ConfigNginxDir the directory that nginx is located in
	ConfigNginxPid = "nginx_pid_file"
)

func main() {

	v := viper.New()
	v.SetEnvPrefix("goz") // eg. env var "GOZ_APID_URI" will bind to config "apid_uri"
	v.AutomaticEnv()
	v.SetDefault(ConfigApidURI, "http://localhost:8181")
	v.SetDefault(ConfigPollWait, "5")

	//use openresty for now.  Must have LUAJIT installed
	v.SetDefault(ConfigNginxDir, " /usr/local/Cellar/openresty/1.9.15.1/")
	v.SetDefault(ConfigNginxPid, "/usr/local/var/run/openresty.pid")

	apidURI := v.GetString(ConfigApidURI)
	timeout := v.GetInt(ConfigPollWait)
	nginxDir := v.GetString(ConfigNginxDir)
	nginxPid := v.GetString(ConfigNginxPid)

	client, err := client.CreateApidClient(apidURI)

	if err != nil {
		log.Fatalf("Could not create cache.  Error is %s", err)
	}

	stageManager := new(nginx.StageManagerImpl)

	manager := nginx.NewManager(client, stageManager, nginxDir, nginxPid, timeout)

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
