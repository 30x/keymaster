package main

import (
	"github.com/spf13/viper"
)

const (
	ConfigApidUri  = "apid_uri"
)

func main() {

	v := viper.New()
	v.SetEnvPrefix("goz") // eg. env var "GOZ_APID_URI" will bind to config "apid_uri"
	v.AutomaticEnv()
	v.SetDefault(ConfigApidUri, "http://localhost:8181/bundles")

	//apidUri := v.Get(ConfigApidUri)

}
