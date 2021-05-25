package setup

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"path"
)

func Config() {
	conf := flag.String("config", "config.toml", "specifies config file location")
	root := flag.String("root", ".", "specifies root path for templates and public files")

	flag.Parse()

	if conf != nil {
		configFile := path.Clean(*conf)

		switch path.Ext(configFile) {
		case "hcl":
			viper.SetConfigType("hcl")
		case "json":
			viper.SetConfigType("json")
		case "toml":
			viper.SetConfigType("toml")
		case "yaml":
			viper.SetConfigType("yaml")
		}

		viper.AddConfigPath(path.Dir(configFile))

		baseName := path.Base(configFile)

		if baseName == "." {
			baseName = "config.toml"
		}

		extName := path.Ext(baseName)
		viper.SetConfigName(baseName[:len(baseName)-len(extName)])
	} else {
		viper.SetConfigType("toml")
		viper.SetConfigName("config")
		viper.AddConfigPath("")
	}

	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	viper.SetDefault("HTTP_BASE", "/")
	viper.SetDefault("HTTP_HOST", "localhost")
	viper.SetDefault("HTTP_SCHEMA", "http")

	viper.SetDefault("DEBUG_USER", 0)

	if root != nil {
		rootpath := path.Clean(*root)
		viper.SetDefault("ROOT_PATH", rootpath)
	}
}
